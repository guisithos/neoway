package services

import (
	"errors"

	"github.com/guisithos/neoway/internal/domain/models"
	"github.com/guisithos/neoway/internal/infrastructure/repository"
	"github.com/guisithos/neoway/util"
)

type ClientService interface {
	CreateClient(client *models.Client) error
	ListClients() ([]models.Client, error)
	BlockClient(document string) error
	UnblockClient(document string) error
	GetClientByDocument(document string) (*models.Client, error)
	GetClientsByName(name string) ([]models.Client, error)
}

type clientService struct {
	repo repository.ClientRepository
}

func NewClientService(repo repository.ClientRepository) ClientService {
	return &clientService{repo: repo}
}

func (s *clientService) CreateClient(client *models.Client) error {
	// Validar o formato do documento
	if err := validateDocument(client.Document, client.Type); err != nil {
		return err
	}

	// Verificar se o cliente já existe
	existing, _ := s.repo.FindByDocument(client.Document)
	if existing != nil {
		return ErrClientAlreadyExists
	}

	return s.repo.Create(client)
}

func validateDocument(document string, clientType models.ClientType) error {
	if clientType == models.PersonType {
		if !util.IsCPF(document) {
			return errors.New("CPF inválido")
		}
	} else if clientType == models.BusinessType {
		if !util.IsCNPJ(document) {
			return errors.New("CNPJ inválido")
		}
	} else {
		return errors.New("tipo de cliente inválido")
	}
	return nil
}

func (s *clientService) ListClients() ([]models.Client, error) {
	return s.repo.ListClients()
}

func (s *clientService) BlockClient(document string) error {
	// Validar se o documento existe
	existing, err := s.repo.FindByDocument(document)
	if err != nil {
		return errors.New("cliente não encontrado")
	}

	if existing.Blocked {
		return errors.New("cliente já está bloqueado")
	}

	return s.repo.BlockClient(document)
}

func (s *clientService) UnblockClient(document string) error {
	// Validar se o documento existe
	existing, err := s.repo.FindByDocument(document)
	if err != nil {
		return errors.New("cliente não encontrado")
	}

	if !existing.Blocked {
		return errors.New("cliente não está bloqueado")
	}

	return s.repo.UnblockClient(document)
}

func (s *clientService) GetClientByDocument(document string) (*models.Client, error) {
	client, err := s.repo.FindByDocument(document)
	if err != nil {
		return nil, errors.New("cliente não encontrado")
	}
	return client, nil
}

func (s *clientService) GetClientsByName(name string) ([]models.Client, error) {
	if name == "" {
		return nil, errors.New("nome não pode estar vazio")
	}
	
	clients, err := s.repo.FindByName(name)
	if err != nil {
		return nil, err
	}
	
	if len(clients) == 0 {
		return nil, errors.New("nenhum cliente encontrado com esse nome")
	}
	
	return clients, nil
}
