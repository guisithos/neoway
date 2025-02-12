package services

import (
	"errors"
	"testing"
	"time"

	"github.com/guisithos/neoway/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repo para teste
type MockClientRepository struct {
	mock.Mock
}

func (m *MockClientRepository) Create(client *models.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *MockClientRepository) FindByDocument(document string) (*models.Client, error) {
	args := m.Called(document)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Client), args.Error(1)
}

func (m *MockClientRepository) ListClients() ([]models.Client, error) {
	args := m.Called()
	return args.Get(0).([]models.Client), args.Error(1)
}

func (m *MockClientRepository) BlockClient(document string) error {
	args := m.Called(document)
	return args.Error(0)
}

func (m *MockClientRepository) UnblockClient(document string) error {
	args := m.Called(document)
	return args.Error(0)
}

func TestCreateClient(t *testing.T) {
	mockRepo := new(MockClientRepository)
	service := NewClientService(mockRepo)

	tests := []struct {
		name        string
		client      *models.Client
		setupMock   func()
		expectError bool
	}{
		{
			name: "registro feito",
			client: &models.Client{
				Name:     "Alqua Ayala",
				Document: "07521001907",
				Type:     models.PersonType,
			},
			setupMock: func() {
				// Mock FindByDocument para retornar nil (cliente não existe)
				mockRepo.On("FindByDocument", "07521001907").Return(nil, nil)
				// Mock Create para suceder
				mockRepo.On("Create", mock.AnythingOfType("*models.Client")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "erro no repo",
			client: &models.Client{
				Name:     "Alqua Ayala",
				Document: "07521001907",
				Type:     models.PersonType,
			},
			setupMock: func() {
				// Mock FindByDocument para retornar nil (cliente não existe)
				mockRepo.On("FindByDocument", "07521001907").Return(nil, nil)
				// Mock Create para retornar erro
				mockRepo.On("Create", mock.AnythingOfType("*models.Client")).Return(errors.New("db error"))
			},
			expectError: true,
		},
		{
			name: "cliente já existe",
			client: &models.Client{
				Name:     "Alqua Ayala",
				Document: "07521001907",
				Type:     models.PersonType,
			},
			setupMock: func() {
				// Mock FindByDocument para retornar cliente existente
				mockRepo.On("FindByDocument", "07521001907").Return(&models.Client{}, nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock entre testes
			mockRepo = new(MockClientRepository)
			service = NewClientService(mockRepo)

			tt.setupMock()
			err := service.CreateClient(tt.client)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetClientByDocument(t *testing.T) {
	mockRepo := new(MockClientRepository)
	service := NewClientService(mockRepo)

	validClient := &models.Client{
		ID:        1,
		Name:      "Alqua Ayala",
		Document:  "07521001907",
		Type:      models.PersonType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		document       string
		setupMock      func()
		expectedClient *models.Client
		expectError    bool
	}{
		{
			name:     "cliente encontrado",
			document: "07521001907",
			setupMock: func() {
				mockRepo.On("FindByDocument", "07521001907").Return(validClient, nil)
			},
			expectedClient: validClient,
			expectError:    false,
		},
		{
			name:     "cliente não encontrado",
			document: "99999999999",
			setupMock: func() {
				mockRepo.On("FindByDocument", "99999999999").Return(nil, errors.New("não encontrado"))
			},
			expectedClient: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			client, err := service.GetClientByDocument(tt.document)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedClient, client)
			}
		})
	}
}

func TestBlockClient(t *testing.T) {
	mockRepo := new(MockClientRepository)
	service := NewClientService(mockRepo)

	tests := []struct {
		name        string
		document    string
		setupMock   func()
		expectError bool
	}{
		{
			name:     "bloqueado",
			document: "07521001907",
			setupMock: func() {
				mockRepo.On("FindByDocument", "07521001907").Return(&models.Client{
					Document: "07521001907",
					Blocked:  false,
				}, nil)
				mockRepo.On("BlockClient", "07521001907").Return(nil)
			},
			expectError: false,
		},
		{
			name:     "cliente já bloqueado",
			document: "07521001907",
			setupMock: func() {
				mockRepo.On("FindByDocument", "07521001907").Return(&models.Client{
					Document: "07521001907",
					Blocked:  true,
				}, nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := service.BlockClient(tt.document)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListClients(t *testing.T) {
	mockRepo := new(MockClientRepository)
	service := NewClientService(mockRepo)

	tests := []struct {
		name        string
		setupMock   func()
		wantClients []models.Client
		wantError   bool
	}{
		{
			name: "lista clientes com sucesso",
			setupMock: func() {
				mockRepo.On("ListClients").Return([]models.Client{
					{
						Document: "09460028438",
						Name:     "Rula Perry",
						Type:     "PF",
						Blocked:  false,
					},
					{
						Document: "74456136000147",
						Name:     "Zotline Ayala",
						Type:     "PJ",
						Blocked:  false,
					},
				}, nil)
			},
			wantClients: []models.Client{
				{
					Document: "09460028438",
					Name:     "Rula Perry",
					Type:     "PF",
					Blocked:  false,
				},
				{
					Document: "74456136000147",
					Name:     "Zotline Ayala",
					Type:     "PJ",
					Blocked:  false,
				},
			},
			wantError: false,
		},
		{
			name: "erro ao listar clientes",
			setupMock: func() {
				mockRepo.On("ListClients").Return(nil, errors.New("erro ao listar clientes"))
			},
			wantClients: nil,
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			clients, err := service.ListClients()

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, clients)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantClients, clients)
			}
		})
	}
}
