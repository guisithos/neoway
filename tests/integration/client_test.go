package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/guisithos/neoway/internal/application/services"
	"github.com/guisithos/neoway/internal/domain/models"
	"github.com/guisithos/neoway/internal/infrastructure/repository"
	"github.com/guisithos/neoway/internal/interfaces/http/handlers"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	// Use environment variables or fallback to default test values
	host := getEnv("TEST_DB_HOST", "localhost")
	user := getEnv("TEST_DB_USER", "postgres")
	password := getEnv("TEST_DB_PASSWORD", "postgres123")
	dbname := getEnv("TEST_DB_NAME", "neoway_test")
	port := getEnv("TEST_DB_PORT", "5433")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	fmt.Printf("Connecting to database with DSN: %s\n", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to test database: %v", err))
	}

	// Migrate tables
	err = db.AutoMigrate(&models.Client{}, &models.RequestMetrics{})
	if err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	// Clean database before each test
	db.Exec("TRUNCATE TABLE clients RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE request_metrics RESTART IDENTITY CASCADE")

	return db
}

// Helper function to get environment variables
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func TestClientCreationFlow(t *testing.T) {
	// Setup test database
	db := setupTestDB()

	// Setup real repositories and services
	clientRepo := repository.NewClientRepository(db)
	clientService := services.NewClientService(clientRepo)
	clientHandler := handlers.NewClientHandler(clientService)

	tests := []struct {
		name           string
		client         models.Client
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful creation and retrieval",
			client: models.Client{
				Name:     "Integration Test Client",
				Document: "07521001907",
				Type:     models.PersonType,
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, rec.Code)

				var response models.Client
				err := json.NewDecoder(rec.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, "Integration Test Client", response.Name)
				assert.Equal(t, "07521001907", response.Document)
			},
		},
		{
			name: "business client creation",
			client: models.Client{
				Name:     "Business Client",
				Document: "12345678901234",
				Type:     models.BusinessType,
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, rec.Code)

				var response models.Client
				err := json.NewDecoder(rec.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, "Business Client", response.Name)
				assert.Equal(t, "12345678901234", response.Document)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create client
			clientJSON, err := json.Marshal(tt.client)
			assert.NoError(t, err)

			createReq := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewBuffer(clientJSON))
			createReq.Header.Set("Content-Type", "application/json")
			createRec := httptest.NewRecorder()

			// Execute creation request
			clientHandler.CreateClient(createRec, createReq)

			// Check creation response
			tt.checkResponse(t, createRec)

			// Retrieve client
			getReq := httptest.NewRequest(http.MethodGet, "/clients/"+tt.client.Document, nil)
			getRec := httptest.NewRecorder()

			// Execute retrieval request
			clientHandler.GetClientByDocument(getRec, getReq)

			// Verify retrieval
			assert.Equal(t, http.StatusOK, getRec.Code)
			var retrievedClient models.Client
			err = json.NewDecoder(getRec.Body).Decode(&retrievedClient)
			assert.NoError(t, err)
			assert.Equal(t, tt.client.Name, retrievedClient.Name)
			assert.Equal(t, tt.client.Document, retrievedClient.Document)
			assert.Equal(t, tt.client.Type, retrievedClient.Type)
		})
	}
}

func TestClientBlockingFlow(t *testing.T) {
	db := setupTestDB()
	clientRepo := repository.NewClientRepository(db)
	clientService := services.NewClientService(clientRepo)
	clientHandler := handlers.NewClientHandler(clientService)

	// First create a client
	client := models.Client{
		Name:     "Block Test Client",
		Document: "07521001907",
		Type:     models.PersonType,
	}

	// Create the client first
	clientJSON, _ := json.Marshal(client)
	createReq := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewBuffer(clientJSON))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	clientHandler.CreateClient(createRec, createReq)
	assert.Equal(t, http.StatusCreated, createRec.Code)

	// Test blocking flow
	t.Run("block client", func(t *testing.T) {
		blockReq := httptest.NewRequest(http.MethodPost, "/clients/"+client.Document+"/block", nil)
		blockRec := httptest.NewRecorder()
		clientHandler.BlockClient(blockRec, blockReq)
		assert.Equal(t, http.StatusOK, blockRec.Code)

		// Verify client is blocked
		getReq := httptest.NewRequest(http.MethodGet, "/clients/"+client.Document, nil)
		getRec := httptest.NewRecorder()
		clientHandler.GetClientByDocument(getRec, getReq)

		var blockedClient models.Client
		json.NewDecoder(getRec.Body).Decode(&blockedClient)
		assert.True(t, blockedClient.Blocked)
	})
}

func TestListClientsFlow(t *testing.T) {
	db := setupTestDB()
	clientRepo := repository.NewClientRepository(db)
	clientService := services.NewClientService(clientRepo)
	clientHandler := handlers.NewClientHandler(clientService)

	// Create multiple clients
	clients := []models.Client{
		{Name: "Client 1", Document: "07521001907", Type: models.PersonType},
		{Name: "Client 2", Document: "12345678901234", Type: models.BusinessType},
	}

	// Create all clients
	for _, client := range clients {
		clientJSON, _ := json.Marshal(client)
		req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewBuffer(clientJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		clientHandler.CreateClient(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	// Test listing
	t.Run("list all clients", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/clients", nil)
		rec := httptest.NewRecorder()
		clientHandler.ListClients(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var listedClients []models.Client
		err := json.NewDecoder(rec.Body).Decode(&listedClients)
		assert.NoError(t, err)
		assert.Len(t, listedClients, len(clients))
	})
}
