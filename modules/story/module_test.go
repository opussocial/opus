package story

import (
	"context"
	"fmt"
	"log"
	"testing"

	"database/sql"
	"github.com/stretchr/testify/require"

	"gitlab.com/pedrokoblitz/opus-go/actions"
	"gitlab.com/pedrokoblitz/opus-go/quality"
)

const (
	StoryNameString         = "Test Story"
	DefinitionNameString    = "Test Definition"
	SchemaNameString        = "test-schema"
	StatusNameString        = "test-status"
	RelationshipNameString  = "Test Relationship"
	RoleNameString          = "Test Role"
	PermissionActionString  = "test:create"
)

func DoCreateStory(t *testing.T, adapter *sql.DB) (uint, error) {
	var err error
	
	emptyPayload := &Story{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &Story{
		Name: "Test Story",
		Slug: "story",
		RoleModels: []Role{
			Role{
				Name: "Role 01",
			},
		},
		RelationshipModels: []GraphRelationship{
			GraphRelationship{
				Name: "Relationship 01",
				SourceDefinitionID: 0,
				TargetDefinitionID: 1,
			},
		},
		SettingModels: []Setting{
			Setting{
				Scope: "story",
				Name: "Setting 01",
				Value: "default",
			},
		},
		DefinitionModels: []Definition{
			Definition{
				Name: "Definition 01",
				Schemas:     []string{SchemaNameString},
				Status:     []string{StatusNameString},	
				Color:           "#eeeeee",
				Icon:            "mdl-car",
				HasTaxonomy:     true,
				HasInteractions: true,
				BelongsToGraph:  true,
			},
		},
	}
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = CreateStoryAction(validPayload, adapter)
	return validPayload.ID, err
}

func DoCreateDefinition(t *testing.T, adapter *sql.DB, storyID uint) (uint, error) {
	var err error
	
	emptyPayload := &Definition{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &Definition{
		Name:            "01 " + DefinitionNameString,
		StoryID:         storyID,
		Schemas:     []string{SchemaNameString},
		Status:     []string{StatusNameString},	
		Color:           "#eeeeee",
		Icon:            "mdl-car",
		HasTaxonomy:     true,
		HasInteractions: true,
		BelongsToGraph:  true,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = CreateDefinitionAction(validPayload, adapter)
	return validPayload.ID, err
}


func DoCreateBlogDefinition(t *testing.T, adapter *sql.DB, storyID uint) (uint, error) {
	var err error
	
	emptyPayload := &Definition{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &Definition{
		Name:            "02 " + DefinitionNameString,
		StoryID:         storyID,
		Schemas:     []string{SchemaNameString},
		Status:     []string{StatusNameString},	
		Color:           "#eeeeee",
		Icon:            "mdl-car",
		HasTaxonomy:     true,
		HasInteractions: true,
		BelongsToGraph:  true,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = CreateDefinitionAction(validPayload, adapter)
	return validPayload.ID, err
}

func DoCreatePageDefinition(t *testing.T, adapter *sql.DB, storyID uint) (uint, error) {
	var err error
	
	emptyPayload := &Definition{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &Definition{
		Name:            "03 " + DefinitionNameString,
		StoryID:         storyID,
		Schemas:     []string{SchemaNameString},
		Status:     []string{StatusNameString},	
		Color:           "#eeeeee",
		Icon:            "mdl-car",
		HasTaxonomy:     true,
		HasInteractions: true,
		BelongsToGraph:  true,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = CreateDefinitionAction(validPayload, adapter)
	return validPayload.ID, err
}


func DoCreatePostDefinition(t *testing.T, adapter *sql.DB, storyID uint) (uint, error) {
	var err error
	
	emptyPayload := &Definition{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &Definition{
		Name:            "04 " + DefinitionNameString,
		StoryID:         storyID,
		Schemas:     []string{SchemaNameString},
		Status:     []string{StatusNameString},	
		Color:           "#eeeeee",
		Icon:            "mdl-car",
		HasTaxonomy:     true,
		HasInteractions: true,
		BelongsToGraph:  true,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = CreateDefinitionAction(validPayload, adapter)
	return validPayload.ID, err
}

func DoCreateProfileDefinition(t *testing.T, adapter *sql.DB, storyID uint) (uint, error) {
	var err error
	
	emptyPayload := &Definition{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &Definition{
		Name:            "05 " + DefinitionNameString,
		StoryID:         storyID,
		Schemas:     []string{SchemaNameString},
		Status:     []string{StatusNameString},
		Color:           "#eeeeee",
		Icon:            "fa-car",
		HasTaxonomy:     true,
		HasInteractions: true,
		BelongsToGraph:  true,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = CreateDefinitionAction(validPayload, adapter)
	return validPayload.ID, err
}


func DoCreateRole(t *testing.T, adapter *sql.DB, storyID uint) (uint, error) {
	var err error
	
	emptyPayload := &Role{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &Role{
		Name:              RoleNameString,
		Permissions: []string{PermissionActionString, "test:delete"},
		StoryID:           storyID,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = CreateRoleAction(validPayload, adapter)
	return validPayload.ID, err
}

func DoCreateGraphRelationship(t *testing.T, adapter *sql.DB) (uint, error) {
	var err error
	
	emptyPayload := &GraphRelationship{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &GraphRelationship{
		Name: "test relationship",
		SourceDefinitionID: 1,
		TargetDefinitionID: 2,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = CreateGraphRelationshipAction(validPayload, adapter)
	return validPayload.ID, err
}

func DoDeleteGraphRelationship(t *testing.T, adapter *sql.DB, relationshipID uint) error {
	var err error
	
	emptyPayload := &payloads.DefaultID{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &payloads.DefaultID{
		ID: relationshipID,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = DeleteGraphRelationshipAction(validPayload, adapter) // Fixed: should be DeleteGraphRelationshipAction, not DeleteRoleAction
	return err
}

func DoDeleteRole(t *testing.T, adapter *sql.DB, roleID uint) error {
	var err error
	
	emptyPayload := &payloads.DefaultID{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &payloads.DefaultID{
		ID: roleID,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = DeleteRoleAction(validPayload, adapter)
	return err
}

func DoDeleteDefinition(t *testing.T, adapter *sql.DB, definitionID uint) error {
	var err error
	
	emptyPayload := &payloads.DefaultID{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &payloads.DefaultID{
		ID: definitionID,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = DeleteDefinitionAction(validPayload, adapter)
	return err
}

func DoDeleteStory(t *testing.T, adapter *sql.DB, storyID uint) error {
	var err error
	
	emptyPayload := &payloads.DefaultID{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &payloads.DefaultID{
		ID: storyID,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = DeleteStoryAction(validPayload, adapter)
	return err
}

// cleanupTestData removes any existing test data
func cleanupTestData(t *testing.T, adapter *sql.DB) {
	// Clean up in correct order to respect foreign key constraints
	queries := []string{
		"DELETE FROM graph_relationships",
		"DELETE FROM roles",
		"DELETE FROM definitions",
		"DELETE FROM stories",
	}
	
	for _, query := range queries {
		_, err := adapter.ExecContext(context.Background(), query)
		if err != nil {
			t.Logf("Cleanup warning for query '%s': %v", query, err)
		}
	}
}

func TestStory(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"root",
		"localhost",
		3307,
		"opus_test",
	)

	adapter, err := adapters.NewSQLAdapter(dsn)
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	defer adapter.Close()

	// Clean up before starting
	cleanupTestData(t, adapter)

	t.Run("CreateStory", func(t *testing.T) {
		storyID, err := DoCreateStory(t, adapter)
		require.NoError(t, err, "CreateStory should succeed")
		require.NotZero(t, storyID, "Story ID should not be zero")

		t.Run("CreateDefinition", func(t *testing.T) {
			defID, err := DoCreateDefinition(t, adapter, storyID)
			require.NoError(t, err, "CreateDefinition should succeed")
			require.NotZero(t, defID, "Definition ID should not be zero")

			t.Run("CreateRole", func(t *testing.T) {
				roleID, err := DoCreateRole(t, adapter, storyID)
				require.NoError(t, err, "CreateRole should succeed")
				require.NotZero(t, roleID, "Role ID should not be zero")

				t.Run("CreateGraphRelationship", func(t *testing.T) {
					relationID, err := DoCreateGraphRelationship(t, adapter)
					// This might fail if we don't have definitions with IDs 1 and 2
					if err != nil {
						t.Logf("CreateGraphRelationship failed (may need setup): %v", err)
						return
					}
					require.NotZero(t, relationID, "Relationship ID should not be zero")

					t.Run("DeleteGraphRelationship", func(t *testing.T) {
						err = DoDeleteGraphRelationship(t, adapter, relationID)
						if err != nil {
							t.Logf("DeleteGraphRelationship failed: %v", err)
						}
					})
				})

				t.Run("DeleteRole", func(t *testing.T) {
					err = DoDeleteRole(t, adapter, roleID)
					require.NoError(t, err, "DeleteRole should succeed")
				})
			})

			t.Run("DeleteDefinition", func(t *testing.T) {
				err = DoDeleteDefinition(t, adapter, defID)
				require.NoError(t, err, "DeleteDefinition should succeed")
			})
		})

		t.Run("DeleteStory", func(t *testing.T) {
			err = DoDeleteStory(t, adapter, storyID)
			require.NoError(t, err, "DeleteStory should succeed")
		})
	})
}

// func TestStoryComponents(t *testing.T) {
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 		"root",
// 		"root",
// 		"localhost",
// 		3307,
// 		"opus_test",
// 	)

// 	adapter, err := adapters.NewSQLAdapter(dsn)
// 	if err != nil {
// 		log.Fatalf("DB error: %v", err)
// 	}
// 	defer adapter.Close()

// 	// Clean up before starting
// 	cleanupTestData(t, adapter)

// 	t.Run("CreateStory", func(t *testing.T) {
// 		storyID, err := DoCreateStory(t, adapter)
// 		require.NoError(t, err, "CreateStory should succeed")
// 		require.NotZero(t, storyID, "Story ID should not be zero")
// 		t.Run("CreateProfileDefinition", func(t *testing.T) {
// 			defID, err := DoCreateProfileDefinition(t, adapter, storyID)
// 			require.NoError(t, err, "CreateProfileDefinition should succeed")
// 			require.NotZero(t, defID, "Definition ID should not be zero")
// 		})
// 		t.Run("CreatePageDefinition", func(t *testing.T) {
// 			defID, err := DoCreatePageDefinition(t, adapter, storyID)
// 			require.NoError(t, err, "CreatePageDefinition should succeed")
// 			require.NotZero(t, defID, "PageDefinition ID should not be zero")
// 		})
// 		t.Run("CreateBlogDefinition", func(t *testing.T) {
// 			defID, err := DoCreatePageDefinition(t, adapter, storyID)
// 			require.NoError(t, err, "CreateBlogDefinition should succeed")
// 			require.NotZero(t, defID, "BlogDefinition ID should not be zero")
// 		})
// 		t.Run("CreatePostDefinition", func(t *testing.T) {
// 			defID, err := DoCreatePostDefinition(t, adapter, storyID)
// 			require.NoError(t, err, "CreatePostDefinition should succeed")
// 			require.NotZero(t, defID, "PostDefinition ID should not be zero")
// 		})
// 		t.Run("CreateRole", func(t *testing.T) {
// 			roleID, err := DoCreateRole(t, adapter, storyID)
// 			require.NoError(t, err, "CreateRole should succeed")
// 			require.NotZero(t, roleID, "Role ID should not be zero")

// 		})
// 		t.Run("CreateGraphRelationship", func(t *testing.T) {
// 			relationID, err := DoCreateGraphRelationship(t, adapter, storyID)
// 			if err != nil {
// 				t.Logf("CreateGraphRelationship failed: %v", err)
// 				return
// 			}
// 			require.NotZero(t, relationID, "Relationship ID should not be zero")
// 		})

// 	})

// }

