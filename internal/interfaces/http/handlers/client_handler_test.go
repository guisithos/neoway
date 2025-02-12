package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/guisithos/neoway/internal/application/services"
	"github.com/guisithos/neoway/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock do ClientService
type MockClientService struct {
	mock.Mock
}

func (m *MockClientService) CreateClient(client *models.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *MockClientService) ListClients() ([]models.Client, error) {
	args := m.Called()
	return args.Get(0).([]models.Client), args.Error(1)
}

func (m *MockClientService) BlockClient(document string) error {
	args := m.Called(document)
	return args.Error(0)
}

func (m *MockClientService) UnblockClient(document string) error {
	args := m.Called(document)
	return args.Error(0)
}

func (m *MockClientService) GetClientByDocument(document string) (*models.Client, error) {
	args := m.Called(document)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Client), args.Error(1)
}

func TestCreateClientHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(mockService *MockClientService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "valida registro",
			requestBody: models.Client{
				Name:     "Nerel Fischer",
				Document: "02440076910",
				Type:     models.PersonType,
			},
			setupMock: func(mockService *MockClientService) {
				mockService.On("CreateClient", mock.AnythingOfType("*models.Client")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "formato CPF inválido",
			requestBody: models.Client{
				Name:     "Test Client",
				Document: "123", // CPF inválido
				Type:     models.PersonType,
			},
			setupMock: func(mockService *MockClientService) {
				// Não precisa mock pois a validação falha antes
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "CPF precisa ter 11 dígitos para o tipo PERSON",
		},
		{
			name: "cliente já existe",
			requestBody: models.Client{
				Name:     "Nerel Fischer",
				Document: "02440076910",
				Type:     models.PersonType,
			},
			setupMock: func(mockService *MockClientService) {
				mockService.On("CreateClient", mock.AnythingOfType("*models.Client")).
					Return(services.ErrClientAlreadyExists)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "cliente já existe",
		},
		{
			name:        "corpo da requisição inválido",
			requestBody: "invalid json",
			setupMock: func(mockService *MockClientService) {
				// Não precisa mock pois o parsing do JSON falha
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "corpo da requisição inválido",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockClientService)
			handler := NewClientHandler(mockService)
			tt.setupMock(mockService)

			// Cria req
			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()

			handler.CreateClient(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestListClientsHandler(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(mockService *MockClientService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "lista clientes",
			setupMock: func(mockService *MockClientService) {
				clients := []models.Client{
					{
						Name:     "Nerel Fischer",
						Document: "02440076910",
						Type:     models.PersonType,
					},
					{
						Name:     "Gionon Elliott",
						Document: "06514032940",
						Type:     models.PersonType,
					},
				}
				mockService.On("ListClients").Return(clients, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "empty list",
			setupMock: func(mockService *MockClientService) {
				mockService.On("ListClients").Return([]models.Client{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "internal error",
			setupMock: func(mockService *MockClientService) {
				mockService.On("ListClients").Return([]models.Client{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "erro ao listar clientes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockClientService)
			handler := NewClientHandler(mockService)
			tt.setupMock(mockService)

			// Cria req
			req := httptest.NewRequest(http.MethodGet, "/clients", nil)
			rec := httptest.NewRecorder()

			handler.ListClients(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
			mockService.AssertExpectations(t)
		})
	}
}
