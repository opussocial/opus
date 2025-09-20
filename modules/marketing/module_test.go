package marketing

import (
	"context"
	"fmt"
	"time"
	"log"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/pedrokoblitz/opus-go/internal/adapters"
	"gitlab.com/pedrokoblitz/opus-go/internal/payloads"
	"gitlab.com/pedrokoblitz/opus-go/internal/quality"
)

func DoCreateLead(t *testing.T, adapter *adapters.SQLAdapter) (uint, error) {
	var err error
	
	emptyPayload := &Lead{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &Lead{
		Name:  "Test Lead",
		Email: "test@example.com",
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = CreateLeadAction(validPayload, adapter)
	return validPayload.ID, err
}

func DoCreateOperation(t *testing.T, adapter *adapters.SQLAdapter) (uint, error) {
	var err error
	
	emptyPayload := &Operation{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &Operation{
		Name:       "Test Operation",
		Slug:       "test-operation",
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = CreateOperationAction(validPayload, adapter)
	return validPayload.ID, err
}

func DoCreateVisit(t *testing.T, adapter *adapters.SQLAdapter, operationID uint, leadID uint) (uint, error) {
	var err error
	
	emptyPayload := &Visit{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)


	validPayload := &Visit{
		VisitDate: time.Now(),
		PageUrl: "/",
		UtmSource: "whatever",
		UtmMedium: "mobile",
		UtmCampaign: "operation-slug",
		LeadID:      leadID,
		ReferrerUrl: "google.com",
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = CreateVisitAction(validPayload, adapter)
	return validPayload.ID, err
}

func DoDeleteVisit(t *testing.T, adapter *adapters.SQLAdapter, visitID uint) error {
	var err error
	
	emptyPayload := &payloads.DefaultID{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &payloads.DefaultID{
		ID: visitID,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = DeleteVisitAction(validPayload, adapter)
	return err
}

func DoDeleteOperation(t *testing.T, adapter *adapters.SQLAdapter, operationID uint) error {
	var err error
	
	emptyPayload := &payloads.DefaultID{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &payloads.DefaultID{
		ID: operationID,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = DeleteOperationAction(validPayload, adapter)
	return nil
}

func DoDeleteLead(t *testing.T, adapter *adapters.SQLAdapter, leadID uint) error {
	var err error
	
	emptyPayload := &payloads.DefaultID{}
	err = emptyPayload.Validate()
	require.ErrorIs(t, err, quality.ErrValidation)

	validPayload := &payloads.DefaultID{
		ID: leadID,
	}
	
	err = validPayload.Validate()
	require.NoError(t, err)
	
	err = DeleteLeadAction(validPayload, adapter)
	return nil
}

// cleanupTestData removes any existing test data
func cleanupTestData(t *testing.T, adapter *adapters.SQLAdapter) {
	// Clean up in correct order to respect foreign key constraints
	queries := []string{
		"DELETE FROM visits",
		"DELETE FROM operations",
		"DELETE FROM leads",
	}
	
	for _, query := range queries {
		_, err := adapter.ExecContext(context.Background(), query)
		if err != nil {
			t.Logf("Cleanup warning for query '%s': %v", query, err)
		}
	}
}

func TestMarketing(t *testing.T) {
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

	t.Run("CreateLead", func(t *testing.T) {
		leadID, err := DoCreateLead(t, adapter)
		require.NoError(t, err, "CreateLead should succeed")
		require.NotZero(t, leadID, "Lead ID should not be zero")

		t.Run("CreateOperation", func(t *testing.T) {
			operationID, err := DoCreateOperation(t, adapter)
			require.NoError(t, err, "CreateOperation should succeed")
			require.NotZero(t, operationID, "Operation ID should not be zero")

			t.Run("CreateVisit", func(t *testing.T) {
				visitID, err := DoCreateVisit(t, adapter, operationID, leadID)
				require.NoError(t, err, "CreateVisit should succeed")
				require.NotZero(t, visitID, "Visit ID should not be zero")

				t.Run("DeleteVisit", func(t *testing.T) {
					err = DoDeleteVisit(t, adapter, visitID)
					require.NoError(t, err, "DeleteVisit should succeed")
				})
			})

			t.Run("DeleteOperation", func(t *testing.T) {
				err = DoDeleteOperation(t, adapter, operationID)
				require.NoError(t, err, "DeleteOperation should succeed")
			})
		})

		t.Run("DeleteLead", func(t *testing.T) {
			err = DoDeleteLead(t, adapter, leadID)
			require.NoError(t, err, "DeleteLead should succeed")
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

// TestIndividualComponents allows testing each component separately
func TestIndividualComponents(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"root",
		"localhost",
		3307,
		"opus_test",
	)

	adapter, err := adapters.NewSQLAdapter(dsn)
	require.NoError(t, err)
	defer adapter.Close()

	cleanupTestData(t, adapter)

	t.Run("LeadOnly", func(t *testing.T) {
		leadID, err := DoCreateLead(t, adapter)
		require.NoError(t, err)
		require.NotZero(t, leadID)

		err = DoDeleteLead(t, adapter, leadID)
		require.NoError(t, err)
	})
}
