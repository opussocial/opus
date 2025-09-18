package story

import (
	"fmt"
	"strings"
	"context"

	"gitlab.com/pedrokoblitz/opus-go/internal/adapters"
	"gitlab.com/pedrokoblitz/opus-go/internal/payloads"
	"gitlab.com/pedrokoblitz/opus-go/internal/quality"
)




/**
 * 
 * DEFINITION STORE
 * 
 * 
 **/

const (

	CreateDefMySQLQuery = "INSERT INTO definitions (name, slug, description, story_id) VALUES (?, ?, ?, ?)"
	UpdateDefMySQLQuery = "UPDATE definitions SET name = ?, description = ? WHERE id = ?"
	DeleteDefMySQLQuery = "DELETE FROM definitions WHERE id = ?"

	GetSchemasBySlugsMySQLQuery = "SELECT id FROM content_schemas "

	BatchAddSchemaMySQLQuery = "INSERT INTO definitions_schemas (definition_id, schema_id) VALUES "
	BatchAddHierarchyMySQLQuery = "INSERT INTO definition_hierarchy (parent_id, child_id) VALUES "
	BatchAddStatusMySQLQuery = "INSERT INTO definition_status (name, slug, definition_id) VALUES "
	BatchAddStateMachineMySQLQuery = "INSERT INTO definition_state_machine (status_id, next_id) VALUES "
	BatchCreateDefPermissionsMySQLQuery = `
		INSERT INTO permissions (name, action) VALUES (?, ?), (?, ?), (?, ?), (?, ?)
	`
    GetHierarchyMySQLQuery = "SELECT parent_id, child_id FROM definition_hierarchy WHERE parent_id = ? OR child_id = ?"
    DeleteHierarchyMySQLQuery = "DELETE FROM definition_hierarchy WHERE parent_id = ? AND child_id = ?"
    GetPossibleRootsQuery = `
        SELECT d.id, d.name, d.slug 
        FROM definitions d
        LEFT JOIN definition_hierarchy dh ON d.id = dh.child_id
        WHERE dh.parent_id IS NULL AND d.story_id = ?
    `
    GetPossibleParentsQuery = `
        SELECT d.id, d.name, d.slug 
        FROM definitions d
        WHERE d.id IN (
            SELECT parent_id FROM definition_hierarchy WHERE child_id = ?
        ) AND d.story_id = ?
    `
    GetPossibleChildrenQuery = `
        SELECT d.id, d.name, d.slug 
        FROM definitions d
        WHERE d.id IN (
            SELECT child_id FROM definition_hierarchy WHERE parent_id = ?
        ) AND d.story_id = ?
    `
)

type DefinitionStore struct {
	adapter adapters.DatabaseAdapter
}

func NewDefinitionStore(adapter adapters.DatabaseAdapter) *DefinitionStore {
	return &DefinitionStore{adapter: adapter.(adapters.DatabaseAdapter)}
}

func (s *DefinitionStore) Create(ctx context.Context, p payloads.Payload) error {
	def := p.(*Definition)
	res, err := s.adapter.ExecContext(ctx,
		CreateDefMySQLQuery,
		def.Name, def.Slug, def.Description, def.StoryID,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	def.ID = uint(id)
	return nil
}

func (s *DefinitionStore) Update(ctx context.Context, p payloads.Payload) error {
	def := p.(*Definition)
	_, err := s.adapter.ExecContext(ctx,
		UpdateDefMySQLQuery,
		def.Name, def.Description, def.ID,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}

func (us *DefinitionStore) Delete(ctx context.Context, p payloads.Payload) error {
	def := p.(*payloads.DefaultID)
	_, err := us.adapter.ExecContext(ctx,
		DeleteDefMySQLQuery,
		def.ID,
	)
	if err != nil {
		return quality.ErrInternal.WithDetail(err.Error())
	}
	return nil
}

func (s *DefinitionStore) GetSchemasBySlugs(ctx context.Context, p payloads.Payload) error {
	def := p.(*Definition)
	query := GetSchemasBySlugsMySQLQuery
	values := ""
	args := make([]interface{}, 0, len(def.Schemas))
	for i, schema := range def.Schemas {
		if i == 0 {
			values += "WHERE schema = ?"
		} else {
			values += " OR schema = ?"
		}
		args = append(args, schema)
	}

	rows, err := s.adapter.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var newSchemaList []DefinitionSchema
	for rows.Next() {
		var schema DefinitionSchema
		err := rows.Scan(&schema.ID)
		if err != nil {
			return err

		}
		newSchemaList = append(newSchemaList, schema)
	}
	def.SchemaModels = newSchemaList
	return nil
}

func (s *DefinitionStore) AddSchemaRelationships(ctx context.Context, p payloads.Payload) error {
	def := p.(*Definition)
	if len(def.SchemaModels) == 0 {
		return nil
	}
	query := BatchAddSchemaMySQLQuery
	values := make([]string, 0, len(def.SchemaModels))
	args := make([]interface{}, 0, len(def.SchemaModels)*2)

	for _, schema := range def.SchemaModels {
		values = append(values, "(?, ?)")
		args = append(args, def.ID, schema.ID)
	}

	query += strings.Join(values, ", ")
	_, err := s.adapter.ExecContext(ctx, query, args...)
	return err
}

func (s *DefinitionStore) AddStatus(ctx context.Context, p payloads.Payload) error {
	def := p.(*Definition)
	if len(def.StatusModels) == 0 {
		return nil
	}

	query := BatchAddStatusMySQLQuery
	values := make([]string, 0, len(def.StatusModels))
	args := make([]interface{}, 0, len(def.StatusModels)*3)

	for _, status := range def.StatusModels {
		values = append(values, "(?, ?, ?)")
		args = append(args, status.Name, status.Slug, def.ID)
	}

	query += strings.Join(values, ", ")
	_, err := s.adapter.ExecContext(ctx, query, args...)
	return err
}

func (s *DefinitionStore) CreatePermissions(ctx context.Context, p payloads.Payload) error {
	def := p.(*Definition)
	query := BatchCreateDefPermissionsMySQLQuery
	slug := def.Slug
	story := def.Story
	createPermission := story + ":" + slug + ":create"	
	updatePermission := story + ":" + slug + ":update"	
	deletePermission := story + ":" + slug + ":delete"	
	viewPermission := 	story + ":" + slug + ":view"
	_, err := s.adapter.ExecContext(
		ctx,
		query,
		"Create " + def.Name,
		createPermission,
		"Update " + def.Name,
		updatePermission,
		"Delete " + def.Name,
		deletePermission,
		"View " + def.Name,
		viewPermission,
	)
	return err
}


func (s *DefinitionStore) SyncSchemaRelationships(ctx context.Context, p payloads.Payload) error {
    def := p.(*Definition)
    
    // First, delete existing relationships
    deleteQuery := "DELETE FROM definitions_schemas WHERE definition_id = ?"
    _, err := s.adapter.ExecContext(ctx, deleteQuery, def.ID)
    if err != nil {
       return quality.ErrInternal.WithDetail(err.Error())
    }
    // Then add the new relationships (reusing existing function)
    return s.AddSchemaRelationships(ctx, p)
}

func (s *DefinitionStore) SyncStatus(ctx context.Context, p payloads.Payload) error {
    def := p.(*Definition)
    
    // First, delete existing statuses
    deleteQuery := "DELETE FROM definition_status WHERE definition_id = ?"
    _, err := s.adapter.ExecContext(ctx, deleteQuery, def.ID)
    if err != nil {
        return quality.ErrInternal.WithDetail(err.Error())
    }
    return s.AddStatus(ctx, p)
}

// You might also want to add this function to handle state machine sync
func (s *DefinitionStore) SyncStateMachine(ctx context.Context, p payloads.Payload) error {
    def := p.(*Definition)
    
    // First, get all status IDs for this definition to delete state machine entries
    statusIDsQuery := "SELECT id FROM definition_status WHERE definition_id = ?"
    rows, err := s.adapter.QueryContext(ctx, statusIDsQuery, def.ID)
    if err != nil {
        return quality.ErrInternal.WithDetail(err.Error())
    }
    defer rows.Close()
    
    var statusIDs []uint
    for rows.Next() {
        var id uint
        err := rows.Scan(&id)
        if err != nil {
            return quality.ErrInternal.WithDetail(err.Error())
        }
        statusIDs = append(statusIDs, id)
    }
    
    // Delete existing state machine relationships for these statuses
    if len(statusIDs) > 0 {
        placeholders := strings.Repeat("?, ", len(statusIDs)-1) + "?"
        deleteQuery := fmt.Sprintf("DELETE FROM definition_state_machine WHERE status_id IN (%s) OR next_id IN (%s)", placeholders, placeholders)
        
        args := make([]interface{}, 0, len(statusIDs)*2)
        for _, id := range statusIDs {
            args = append(args, id)
        }
        for _, id := range statusIDs {
            args = append(args, id)
        }
        
        _, err = s.adapter.ExecContext(ctx, deleteQuery, args...)
        if err != nil {
            return quality.ErrInternal.WithDetail(err.Error())
        }
    }
    
    return nil
}
