package story

import (
	"fmt"
	"time"
	"strings"
	"context"
	"database/sql"

	"gitlab.com/pedrokoblitz/opus-go/actions"
	"gitlab.com/pedrokoblitz/opus-go/quality"
)

/**
 * 
 * STORY STORE
 * 
 * 



 **/

const (
	CreateStoryMySQLQuery = "INSERT INTO stories (name, slug, description) VALUES (?, ?, ?)"
	UpdateStoryMySQLQuery = "UPDATE stories SET name = ?, slug = ?, description = ? WHERE id = ?"
	DeleteStoryMySQLQuery = "DELETE FROM stories WHERE id = ?"

	CheckSlugMySQLQuery = "SELECT id FROM stories WHERE slug = ?"

	BatchAddRolesMySQLQuery = "INSERT INTO roles (name, slug, description, story_id) VALUES "
	BatchAddDefsMySQLQuery = "INSERT INTO definitions (name, slug, description, story_id) VALUES "
	BatchAddSettingsMySQLQuery = "INSERT INTO settings (name, slug, description, story_id) VALUES "
	BatchAddRelationshipsMySQLQuery = "INSERT INTO graph_relationships (name, slug, description, source_definition_id, target_definition_id) VALUES "

	GetStoryBySlugMySQLQuery = `
		SELECT id, name, slug, created_at
		FROM stories 
		WHERE slug = ?
		AND deleted_at IS NULL;
	`

	GetStoryByIDMySQLQuery = `
		SELECT id, name, slug, created_at
		FROM stories 
		WHERE id = ? 
		AND deleted_at IS NULL;
	`

	GetStoryRolesMySQLQuery = "SELECT id, name, description FROM roles WHERE story_id = ?"
	GetStorySettingsMySQLQuery = "SELECT id, name, value FROM settings WHERE story_id = ?"
	GetStoryRelationshipsMySQLQuery = "SELECT id, name, source_definition_id. target_definition_id FROM graph_relationships WHERE story_id = ?"
	GetStoryDefinitionsMySQLQuery = "SELECT id, name, description, created_at FROM definitions WHERE story_id = ?"

	// TODO
	GetRolePermissionsMySQLQuery = ""
	GetSchemasMySQLQuery = ""
)

type StoryStore struct {
	adapter *sql.DB
}

func NewStoryStore(adapter *sql.DB) *StoryStore {
	return &StoryStore{adapter: adapter}
}

func (s *StoryStore) CheckSlug(ctx context.Context, p actions.Payload) error {
    story := p.(*Story)
    err := s.adapter.QueryRowContext(ctx, CheckSlugMySQLQuery, story.Slug).Scan(&story.ID)

    if err == sql.ErrNoRows {
        return nil
    }
    if err != nil {
        return quality.ErrInternal.WithDetail(err.Error())
    }

    return quality.ErrInternal.WithDetail("slug already exists")
}

// Get retrieves a single story by ID
func (s *StoryStore) Get(ctx context.Context, p actions.Payload) error {
    story := p.(*Story)
    err := s.adapter.QueryRowContext(ctx, GetStoryByIDMySQLQuery, story.ID).Scan(
        &story.ID,
        &story.Name,
        &story.Slug,
        &story.CreatedAt,
        // &story.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return quality.ErrNotFound
    }
    if err != nil {
        return quality.ErrInternal.WithDetail(err.Error())
    }

    return nil
}

// Get retrieves a single story
func (s *StoryStore) GetBySlug(ctx context.Context, p actions.Payload) error {
    story := p.(*Story)
    err := s.adapter.QueryRowContext(ctx, GetStoryBySlugMySQLQuery, story.Slug).Scan(
        &story.ID,
        &story.Name,
        &story.Slug,
        &story.CreatedAt,
        // &story.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return quality.ErrNotFound
    }
    if err != nil {
        return quality.ErrInternal.WithDetail(err.Error())
    }

    return nil
}


func (s *StoryStore) Create(ctx context.Context, p actions.Payload) error {
  story := p.(*Story)
	res, err := s.adapter.ExecContext(ctx,
		CreateStoryMySQLQuery,
		story.Name, story.Slug, story.Description,
	)

	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	story.ID = uint(id)

	return nil
}

func (s *StoryStore) Update(ctx context.Context, p actions.Payload) error {
    story := p.(*Story)
	_, err := s.adapter.ExecContext(ctx,
		UpdateStoryMySQLQuery,
		story.Name, story.Slug, story.Description,
	)

	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}

func (s *StoryStore) Delete(ctx context.Context, p actions.Payload) error {
  story := p.(*actions.DefaultID)
	_, err := s.adapter.ExecContext(ctx,
		DeleteStoryMySQLQuery,
		story.ID,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}

func (s *StoryStore) AddRoles(ctx context.Context, p actions.Payload) error {
	story := p.(*Story)
	if len(story.RoleModels) == 0 {
		return nil
	}

	query := BatchAddRolesMySQLQuery
	values := make([]string, 0, len(story.RoleModels))
	args := make([]interface{}, 0, len(story.RoleModels)*3)

	for _, role := range story.RoleModels {
		if role.ID == 0 {
			values = append(values, "(?, ?, ?, ?)")
			args = append(args, role.Name, role.Slug, role.Description, story.ID)
		}
	}

	if len(values) == 0 || len(args) == 0 {
		return nil
	}

	query += strings.Join(values, ", ")
	_, err := s.adapter.ExecContext(ctx, query, args...)
	return err
}

func (s *StoryStore) AddRelationships(ctx context.Context, p actions.Payload) error {
	story := p.(*Story)
	if len(story.RoleModels) == 0 {
		return nil
	}

	query := BatchAddRelationshipsMySQLQuery
	values := make([]string, 0, len(story.RelationshipModels))
	args := make([]interface{}, 0, len(story.RelationshipModels)*3)

	for _, rel := range story.RelationshipModels {
		if rel.ID == 0 {
			values = append(values, "(?, ?, ?)")
			args = append(args, rel.Name, rel.Slug, rel.SourceDefinitionID, rel.TargetDefinitionID, story.ID)
		}
	}

	if len(values) == 0 || len(args) == 0 {
		return nil
	}

	query += strings.Join(values, ", ")
	_, err := s.adapter.ExecContext(ctx, query, args...)
	return err
}

func (s *StoryStore) AddDefinitions(ctx context.Context, p actions.Payload) error {
	story := p.(*Story)
	if len(story.DefinitionModels) == 0 {
		return nil
	}

	query := BatchAddDefsMySQLQuery
	values := make([]string, 0, len(story.DefinitionModels))
	args := make([]interface{}, 0, len(story.DefinitionModels)*3)

	for _, def := range story.DefinitionModels {
		values = append(values, "(?, ?, ?)")
		args = append(args, def.Name, def.Slug, story.ID)
	}

	if len(values) == 0 || len(args) == 0 {
		return nil
	}

	query += strings.Join(values, ", ")
	_, err := s.adapter.ExecContext(ctx, query, args...)
	return err
}

func (s *StoryStore) AddSettings(ctx context.Context, p actions.Payload) error {
	story := p.(*Story)
	if len(story.SettingModels) == 0 {
		return nil
	}

	query := BatchAddSettingsMySQLQuery
	values := make([]string, 0, len(story.SettingModels))
	args := make([]interface{}, 0, len(story.SettingModels)*3)

	for _, setting := range story.SettingModels {
		values = append(values, "(?, ?, ?)")
		args = append(args, setting.Name, setting.Slug, story.ID)
	}

	if len(values) == 0 || len(args) == 0 {
		return nil
	}

	query += strings.Join(values, ", ")
	_, err := s.adapter.ExecContext(ctx, query, args...)
	return err
}

func (s *StoryStore) Index(ctx context.Context, p actions.Payload) error {
	stories := p.(*StoryResults)
  query := "SELECT id, name, slug, created_at FROM stories"

  // Execute query
  // rows, err := s.adapter.QueryContext(ctx, query, args...)
  rows, err := s.adapter.QueryContext(ctx, query)
  if err != nil {
    return quality.ErrInternal.WithDetail(fmt.Sprintf("query error: %v", err))
  }
  defer rows.Close()

  for rows.Next() {
    var story Story
    err := rows.Scan(
      &story.ID,
      &story.Name,
      &story.Slug,
      &story.CreatedAt,
    )
    if err != nil {
      return quality.ErrInternal.WithDetail(fmt.Sprintf("scan error: %v", err))
    }
    stories.Results = append(stories.Results, story)
  }

  if err = rows.Err(); err != nil {
    return quality.ErrInternal.WithDetail(fmt.Sprintf("rows error: %v", err))
  }
	return nil
}

func (s *StoryStore) IndexWithFilters(ctx context.Context, p actions.Payload) error {
  query := "SELECT id, name, slug, active, created_at, updated_at FROM stories"
  args := []interface{}{}
  whereClauses := []string{}
  stories := p.(*StoryResults)

  if name, ok := stories.Filters["name"].(string); ok {
    whereClauses = append(whereClauses, "name LIKE ?")
    args = append(args, "%"+name+"%")
  }

  if active, ok := stories.Filters["active"].(bool); ok {
    whereClauses = append(whereClauses, "active = ?")
    args = append(args, active)
  }

  if createdAfter, ok := stories.Filters["createdAfter"].(time.Time); ok {
    whereClauses = append(whereClauses, "created_at >= ?")
    args = append(args, createdAfter)
  }

  if createdBefore, ok := stories.Filters["createdBefore"].(time.Time); ok {
    whereClauses = append(whereClauses, "created_at <= ?")
    args = append(args, createdBefore)
  }

  if updatedAfter, ok := stories.Filters["updatedAfter"].(time.Time); ok {
    whereClauses = append(whereClauses, "updated_at >= ?")
    args = append(args,updatedAfter)
  }

  if updatedBefore, ok := stories.Filters["updatedBefore"].(time.Time); ok {
    whereClauses = append(whereClauses, "updated_at <= ?")
    args = append(args,updatedBefore)
  }

  // Build WHERE clause if we have filters
  if len(whereClauses) > 0 {
    query += " WHERE " + strings.Join(whereClauses, " AND ")
  }

  // Default ordering
  query += " ORDER BY created_at DESC"

  // Apply pagination if specified
  if page, pageOk := stories.Filters["page"].(int); pageOk {
    pageSize := 10 // default
    if size, sizeOk := stories.Filters["pageSize"].(int); sizeOk {
      pageSize = size
    }
    offset := (page - 1) * pageSize
    query += " LIMIT ? OFFSET ?"
    args = append(args, pageSize, offset)
  }

  // Execute query
  rows, err := s.adapter.QueryContext(ctx, query, args...)
  if err != nil {
    return quality.ErrInternal.WithDetail(fmt.Sprintf("query error: %v", err))
  }
  defer rows.Close()

  for rows.Next() {
    var story Story
    err := rows.Scan(
      &story.ID,
      &story.Name,
      &story.Slug,
      &story.Active,
      &story.CreatedAt,
      &story.UpdatedAt,
    )
    if err != nil {
      return quality.ErrInternal.WithDetail(fmt.Sprintf("scan error: %v", err))
    }
    stories.Results = append(stories.Results, story)
  }

  if err = rows.Err(); err != nil {
    return quality.ErrInternal.WithDetail(fmt.Sprintf("rows error: %v", err))
  }

  return nil
}

func (s *StoryStore) GetPermissionsForRole(ctx context.Context, p actions.Payload) error {
	role := p.(*Role)
	rows, err := s.adapter.QueryContext(
		ctx,
		GetRolePermissionsMySQLQuery,
		role.ID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	var newPermissionList []Permission
	for rows.Next() {
		permission := Permission{}
		err := rows.Scan(&permission.ID, &permission.Name, &permission.Action)
		if err != nil {
			return err
		}
		newPermissionList = append(newPermissionList, permission)
	}
	role.PermissionModels = newPermissionList

	return nil
}

func (s *StoryStore) GetRoles(ctx context.Context, p actions.Payload) error {
	story := p.(*Story)
	rows, err := s.adapter.QueryContext(
		ctx,
		GetStoryRolesMySQLQuery,
		story.ID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	// var newRoleList []Role
	for rows.Next() {

		role := Role{StoryID: story.ID}
		err := rows.Scan(&role.ID, &role.Name, &role.Description)
		if err != nil {
			return err
		}
		// err = s.GetPermissionsForRole(ctx, &role)
		// if err != nil {
		// 	return err
		// }
		// newRoleList = append(newRoleList, role)
		story.RoleModels = append(story.RoleModels, role)
	}
	// story.RoleModels = newRoleList
	return nil
}

func (s *StoryStore) GetDefinitions(ctx context.Context, p actions.Payload) error {
	story := p.(*Story)
	rows, err := s.adapter.QueryContext(
		ctx,
		GetStoryDefinitionsMySQLQuery,
		story.ID,
	)
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		def := Definition{StoryID: story.ID}
		err := rows.Scan(&def.ID, &def.Name, &def.Description, &def.CreatedAt)
		if err != nil {
			return err
		}
		// err = s.GetSchemasForDefinition(ctx, &def)
		// if err != nil {
		// 	return err
		// }
		story.DefinitionModels = append(story.DefinitionModels, def)
	}
	return nil
}

func (s *StoryStore) GetSchemasForDefinition(ctx context.Context, p actions.Payload) error {
	definition := p.(*Definition)
	rows, err := s.adapter.QueryContext(
		ctx,
		GetSchemasMySQLQuery,
		definition.ID,
	)
	defer rows.Close()

	if err != nil {
		return err
	}
	var newSchemaList []DefinitionSchema
	for rows.Next() {
		schema := DefinitionSchema{}
		err := rows.Scan(&schema.ID, &schema.Name, &schema.Description)
		if err != nil {
			return err
		}
		newSchemaList = append(newSchemaList, schema)
	}
	definition.SchemaModels = newSchemaList
	return nil
}

func (s *StoryStore) GetSettings(ctx context.Context, p actions.Payload) error {
	story := p.(*Story)
	rows, err := s.adapter.QueryContext(
		ctx,
		GetStorySettingsMySQLQuery,
		story.ID,
	)
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		setting := Setting{StoryID: story.ID}
		err := rows.Scan(&setting.ID, &setting.Name, &setting.Value)
		if err != nil {
			return err
		}
		story.SettingModels = append(story.SettingModels, setting)
	}
	return nil
}

func (s *StoryStore) GetRelationships(ctx context.Context, p actions.Payload) error {
	story := p.(*Story)
	rows, err := s.adapter.QueryContext(
		ctx,
		GetStoryRelationshipsMySQLQuery,
		story.ID,
	)
	defer rows.Close()
	if err != nil {
		return err
	}
	var newRelationshipList []GraphRelationship
	for rows.Next() {
		relationship := GraphRelationship{StoryID: story.ID}
		err := rows.Scan(&relationship.ID, &relationship.Name, &relationship.SourceDefinitionID, &relationship.TargetDefinitionID)
		if err != nil {
			return err
		}
		newRelationshipList = append(newRelationshipList, relationship)
	}
	story.RelationshipModels = newRelationshipList
	return nil
}
