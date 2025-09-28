package story

import (
	"strings"
	"context"
    "database/sql"
    "gitlab.com/pedrokoblitz/opus-go/actions"
    "gitlab.com/pedrokoblitz/opus-go/quality"
)

/**
 * 
 * ROLE STORE
 * 
 * 
 * 
 **/

const (
	CreateRoleMySQLQuery = `
		INSERT INTO roles (name, slug, description, story_id) VALUES (?, ?, ?, ?)
	`
	UpdateRoleMySQLQuery = "UPDATE roles SET name = ?, slug = ?, description = ? WHERE id = ?"
	DeleteRoleMySQLQuery = "DELETE FROM roles WHERE id = ?"
	DeleteRoleByStoryIDMySQLQuery = "DELETE FROM roles WHERE story_id = ?"

	GetPermissionsByActionsMySQLQuery = "SELECT id FROM permissions "
	BatchAddPermissionToRoleMySQLQuery = "INSERT INTO role_permissions (role_id, permission_id) VALUES "

	InsertPermissionMySQLQuery = "INSERT INTO permissions (action, name) VALUES (?, ?)"
	DeletePermissionMySQLQuery = "DELETE FROM permissions WHERE id = ?"
	InsertRolePermissionMySQLQuery = "INSERT INTO roles_permissions (role_id, permission_id) VALUES (?, ?)"
	DeleteRolePermissionMySQLQuery = "DELETE FROM roles_permissions WHERE role_id = ? AND permission_id = ?"

    
    GetRolePermissionsQuery = `
        SELECT p.id, p.action, p.name 
        FROM permissions p
        INNER JOIN roles_permissions rp ON p.id = rp.permission_id
        WHERE rp.role_id = ?
    `
    DeleteRolePermissionsQuery = "DELETE FROM roles_permissions WHERE role_id = ?"
    GetPermissionByActionQuery = "SELECT id FROM permissions WHERE action = ?"
    GetRolesByStoryQuery = "SELECT id, name, slug, description FROM roles WHERE story_id = ?"
    GetRoleBySlugQuery = "SELECT id, name, slug, description, story_id FROM roles WHERE slug = ? AND story_id = ?"
)

type RoleStore struct {
	adapter *sql.DB
}

func NewRoleStore(adapter *sql.DB) *RoleStore {
	return &RoleStore{adapter: adapter.(*sql.DB)}
}

func (s *RoleStore) Create(ctx context.Context, p payloads.Payload) error {
	role := p.(*Role)
    var storyID uint
    err := s.adapter.QueryRowContext(ctx, "SELECT id FROM stories WHERE slug = ?", role.Story).Scan(&storyID)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
    role.StoryID = storyID

	res, err := s.adapter.ExecContext(ctx,
		CreateRoleMySQLQuery,
		role.Name, role.Slug, role.Description, role.StoryID,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	role.ID = uint(id)

	return nil
}

func (s *RoleStore) Update(ctx context.Context, p payloads.Payload) error {
	role := p.(*Role)
	_, err := s.adapter.ExecContext(ctx,
		UpdateRoleMySQLQuery,
		role.Name, role.Slug, role.Description, role.ID,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}

func (s *RoleStore) Delete(ctx context.Context, p payloads.Payload) error {
	role := p.(*payloads.DefaultID)
	_, err := s.adapter.ExecContext(ctx,
		DeleteRoleMySQLQuery,
		role.ID,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}

func (s *RoleStore) GetPermissionsByActions(ctx context.Context, p payloads.Payload) error {
    role := p.(*Role)
    if len(role.Permissions) == 0 {
        return nil
    }
    query := GetPermissionsByActionsMySQLQuery
	values := ""
	args := make([]interface{}, 0, len(role.Permissions))
	for i, action := range role.Permissions {
		if i == 0 {
			values += "WHERE action = ?"
		} else {
			values += " OR action = ?"
		}
		args = append(args, action)
	}

	rows, err := s.adapter.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var permission Permission
		err := rows.Scan(&permission.ID)
		if err != nil {
			return err
		}
		role.PermissionModels = append(role.PermissionModels, permission)
	}
	return nil
}

// AddPermissions relates specific permissions to role
func (s *RoleStore) AddPermissionRelationships(ctx context.Context, p payloads.Payload) error {
    role := p.(*Role)
    if len(role.PermissionModels) == 0 {
        return nil
    }

    values := make([]string, 0, len(role.PermissionModels))
    args := make([]interface{}, 0, len(role.PermissionModels)*2)

    for _, perm := range role.PermissionModels {
        values = append(values, "(?, ?)")
        args = append(args, role.ID, perm.ID)
    }

    // Create local query string
    query := "INSERT INTO role_permissions (role_id, permission_id) VALUES " + 
             strings.Join(values, ", ")
    
    _, err := s.adapter.ExecContext(ctx, query, args...)
    return err
}


// SyncPermissionRelationships syncs permission relationships for a role
func (s *RoleStore) SyncPermissionRelationships(ctx context.Context, p payloads.Payload) error {
    role := p.(*Role)
    
    // First, delete existing permission relationships
    _, err := s.adapter.ExecContext(ctx, DeleteRolePermissionsQuery, role.ID)
    if err != nil {
        return quality.ErrInternal.WithDetail(err.Error())
    }
    
    // Then add the new permission relationships
    return s.AddPermissionRelationships(ctx, p)
}

// GetPermissions retrieves all permissions for a role
func (s *RoleStore) GetPermissions(ctx context.Context, roleID uint) ([]Permission, error) {
    rows, err := s.adapter.QueryContext(ctx, GetRolePermissionsQuery, roleID)
    if err != nil {
        return nil, quality.ErrInternal.WithDetail(err.Error())
    }
    defer rows.Close()
    
    var permissions []Permission
    for rows.Next() {
        var perm Permission
        err := rows.Scan(&perm.ID, &perm.Action, &perm.Name)
        if err != nil {
            return nil, quality.ErrInternal.WithDetail(err.Error())
        }
        permissions = append(permissions, perm)
    }
    
    return permissions, nil
}

// GetPermissionIDByAction retrieves a permission ID by its action string
func (s *RoleStore) GetPermissionIDByAction(ctx context.Context, action string) (uint, error) {
    var permissionID uint
    err := s.adapter.QueryRowContext(ctx, GetPermissionByActionQuery, action).Scan(&permissionID)
    if err != nil {
        return 0, quality.ErrInternal.WithDetail(err.Error())
    }
    return permissionID, nil
}

// GetRolesByStory retrieves all roles for a specific story
func (s *RoleStore) GetRolesByStory(ctx context.Context, storyID uint) ([]Role, error) {
    rows, err := s.adapter.QueryContext(ctx, GetRolesByStoryQuery, storyID)
    if err != nil {
        return nil, quality.ErrInternal.WithDetail(err.Error())
    }
    defer rows.Close()
    
    var roles []Role
    for rows.Next() {
        var role Role
        err := rows.Scan(&role.ID, &role.Name, &role.Slug, &role.Description)
        if err != nil {
            return nil, quality.ErrInternal.WithDetail(err.Error())
        }
        role.StoryID = storyID
        roles = append(roles, role)
    }
    
    return roles, nil
}

// GetRoleBySlug retrieves a role by its slug and story ID
func (s *RoleStore) GetRoleBySlug(ctx context.Context, slug string, storyID uint) (*Role, error) {
    var role Role
    err := s.adapter.QueryRowContext(ctx, GetRoleBySlugQuery, slug, storyID).Scan(
        &role.ID, &role.Name, &role.Slug, &role.Description, &role.StoryID,
    )
    if err != nil {
        return nil, quality.ErrInternal.WithDetail(err.Error())
    }
    
    return &role, nil
}

// AddSinglePermission adds a single permission to a role
func (s *RoleStore) AddSinglePermission(ctx context.Context, roleID, permissionID uint) error {
    _, err := s.adapter.ExecContext(ctx, InsertRolePermissionMySQLQuery, roleID, permissionID)
    if err != nil {
        return quality.ErrInternal.WithDetail(err.Error())
    }
    return nil
}

// RemoveSinglePermission removes a single permission from a role
func (s *RoleStore) RemoveSinglePermission(ctx context.Context, roleID, permissionID uint) error {
    _, err := s.adapter.ExecContext(ctx, DeleteRolePermissionMySQLQuery, roleID, permissionID)
    if err != nil {
        return quality.ErrInternal.WithDetail(err.Error())
    }
    return nil
}

// CreatePermission creates a new permission
func (s *RoleStore) CreatePermission(ctx context.Context, action, name string) (uint, error) {
    res, err := s.adapter.ExecContext(ctx, InsertPermissionMySQLQuery, action, name)
    if err != nil {
        return 0, quality.ErrInternal.WithDetail(err.Error())
    }
    
    id, err := res.LastInsertId()
    if err != nil {
        return 0, quality.ErrInternal.WithDetail(err.Error())
    }
    
    return uint(id), nil
}

// DeletePermission deletes a permission
func (s *RoleStore) DeletePermission(ctx context.Context, permissionID uint) error {
    _, err := s.adapter.ExecContext(ctx, DeletePermissionMySQLQuery, permissionID)
    if err != nil {
        return quality.ErrInternal.WithDetail(err.Error())
    }
    return nil
}

// BulkAddPermissions adds multiple permissions to a role in a single transaction
func (s *RoleStore) BulkAddPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
    if len(permissionIDs) == 0 {
        return nil
    }
    
    query := BatchAddPermissionToRoleMySQLQuery
    values := make([]string, 0, len(permissionIDs))
    args := make([]interface{}, 0, len(permissionIDs)*2)
    
    for _, permID := range permissionIDs {
        values = append(values, "(?, ?)")
        args = append(args, roleID, permID)
    }
    
    query += strings.Join(values, ", ")
    _, err := s.adapter.ExecContext(ctx, query, args...)
    if err != nil {
        return quality.ErrInternal.WithDetail(err.Error())
    }
    
    return nil
}
