package elements

import (
	"fmt"
	"log"
	// "time"
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/pedrokoblitz/opus-go/internal/adapters"
	"gitlab.com/pedrokoblitz/opus-go/internal/payloads"
	"gitlab.com/pedrokoblitz/opus-go/internal/quality"
)

func DoCreateProfile(t *testing.T, adapter *adapters.SQLAdapter) (uint, error) {
	var err error
	
	emptyPayload := &Element{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &Element{
		Name:    "Test Profile",
		DefinitionID: 1,
		Person: &PersonSchema{
			Prefix: "",
			GivenName: "",
			AdditionalName: "",
			FamilyName: "",
			AlternateName: "",
			Suffix: "",

			Gender: "",
			// BirthDate: "",
			Nationality: "",
			BirthPlace: "",
		},
		CreatedByProfileID: 1,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = CreateElementAction(validPayload, adapter)
	return validPayload.ID, err
}

func DoCreateElement(t *testing.T, adapter *adapters.SQLAdapter, profileID uint) (uint, error) {
	var err error
	
	emptyPayload := &Element{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &Element{
		Name:      "Test Element",
		CreatedByProfileID: profileID,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = CreateElementAction(validPayload, adapter)
	return validPayload.ID, nil
}

func DoCreateInteraction(t *testing.T, adapter *adapters.SQLAdapter, profileID uint, elementID uint) (uint, error) {
	var err error
	
	emptyPayload := &Interaction{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &Interaction{
		ElementID: elementID,
		CreatedByProfileID: profileID,
		Scope:     "comment",
		Comment:   "Test comment",
		Rating:    5,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = CreateInteractionAction(validPayload, adapter)
	return validPayload.ID, err
}

func DoCreateKeyword(t *testing.T, adapter *adapters.SQLAdapter, elementID uint) (uint, error) {
	var err error
	
	emptyPayload := &Keyword{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &Keyword{
		Term:     "Test Keyword",
	}	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = CreateKeywordAction(validPayload, adapter)
	return validPayload.ID, err
}

func DoRemoveKeyword(t *testing.T, adapter *adapters.SQLAdapter, keywordID uint) error {
	var err error
	
	emptyPayload := &payloads.DefaultID{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &payloads.DefaultID{
		ID: keywordID,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = DeleteKeywordAction(validPayload, adapter)
	return err
}

func DoDeleteInteraction(t *testing.T, adapter *adapters.SQLAdapter, interactionID uint) error {
	var err error
	
	emptyPayload := &payloads.DefaultID{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &payloads.DefaultID{
		ID: interactionID,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = DeleteInteractionAction(validPayload, adapter)
	return err
}

func DoDeleteElement(t *testing.T, adapter *adapters.SQLAdapter, elementID uint) error {
	var err error
	
	emptyPayload := &payloads.DefaultID{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &payloads.DefaultID{
		ID: elementID,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = DeleteElementAction(validPayload, adapter)
	return err
}

func DoDeleteProfile(t *testing.T, adapter *adapters.SQLAdapter, profileID uint) error {
	var err error
	
	emptyPayload := &payloads.DefaultID{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &payloads.DefaultID{
		ID: profileID,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = DeleteElementAction(validPayload, adapter)
	return err
}

// cleanupTestData removes any existing test data
func cleanupTestData(t *testing.T, adapter *adapters.SQLAdapter) {
	// Clean up in correct order to respect foreign key constraints
	queries := []string{
		"DELETE FROM taxonomy",
		"DELETE FROM interactions",
		"DELETE FROM elements",
		"DELETE FROM profile_relationships",
	}
	
	for _, query := range queries {
		_, err := adapter.ExecContext(context.Background(), query)
		if err != nil {
			t.Logf("Cleanup warning for query '%s': %v", query, err)
		}
	}
}

func TestElement(t *testing.T) {
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
	var elementID uint

	t.Run("CreateProfile", func(t *testing.T) {
		profileID, err := DoCreateProfile(t, adapter)
		require.NoError(t, err, "CreateProfile should succeed")
		require.NotZero(t, profileID, "Profile ID should not be zero")

		t.Run("CreateElement", func(t *testing.T) {
			elementID, err = DoCreateElement(t, adapter, profileID)
			require.NoError(t, err, "CreateElement should succeed")
			require.NotZero(t, elementID, "Element ID should not be zero")
		})

		t.Run("CreateInteraction", func(t *testing.T) {
			interactionID, err := DoCreateInteraction(t, adapter, profileID, elementID)
			require.NoError(t, err, "CreateInteraction should succeed")
			require.NotZero(t, interactionID, "Interaction ID should not be zero")

			t.Run("DeleteInteraction", func(t *testing.T) {
				err = DoDeleteInteraction(t, adapter, interactionID)
				require.NoError(t, err, "DeleteInteraction should succeed")
			})
		})

		t.Run("CreateKeyword", func(t *testing.T) {
			keywordID, err := DoCreateKeyword(t, adapter, elementID)
			require.NoError(t, err, "CreateKeyword should succeed")
			require.NotZero(t, keywordID, "Keyword ID should not be zero")

			t.Run("RemoveKeyword", func(t *testing.T) {
				err = DoRemoveKeyword(t, adapter, keywordID)
				require.NoError(t, err, "RemoveKeyword should succeed")
			})
		})

		t.Run("DeleteElement", func(t *testing.T) {
			err = DoDeleteElement(t, adapter, elementID)
			require.NoError(t, err, "DeleteElement should succeed")
		})

		t.Run("DeleteProfile", func(t *testing.T) {
			err = DoDeleteProfile(t, adapter, profileID)
			require.NoError(t, err, "DeleteProfile should succeed")
		})
	})
}

// TestDatabaseConnection tests basic database connectivity
func TestDatabaseConnection(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"root",
		"localhost",
		3307,
		"opus_test",
	)

	adapter, err := adapters.NewSQLAdapter(dsn)
	require.NoError(t, err, "Database should connect successfully")
	defer adapter.Close()

	// Test simple query
	var result int
	err = adapter.QueryRowContext(context.Background(), "SELECT 1").Scan(&result)
	require.NoError(t, err, "Simple query should work")
	require.Equal(t, 1, result, "Query should return 1")
}

// // TestIndividualComponents allows testing each component separately
// func TestIndividualComponents(t *testing.T) {
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 		"root",
// 		"root",
// 		"localhost",
// 		3307,
// 		"opus_test",
// 	)

// 	adapter, err := adapters.NewSQLAdapter(dsn)
// 	require.NoError(t, err)
// 	defer adapter.Close()

// 	cleanupTestData(t, adapter)

// 	t.Run("ProfileOnly", func(t *testing.T) {
// 		profileID, err := DoCreateProfile(t, adapter)
// 		require.NoError(t, err)
// 		require.NotZero(t, profileID)

// 		err = DoDeleteProfile(t, adapter, profileID)
// 		require.NoError(t, err)
// 	})
// }
