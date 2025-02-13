package repository

import (
	"errors"
	"log"

	"github.com/guisithos/neoway/internal/domain/models"
	"gorm.io/gorm"
)

type ClientRepository interface {
	Create(client *models.Client) error
	FindByDocument(document string) (*models.Client, error)
	FindByName(name string) ([]models.Client, error)
	ListClients() ([]models.Client, error)
	BlockClient(document string) error
	UnblockClient(document string) error
}

type clientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) Create(client *models.Client) error {
	return r.db.Create(client).Error
}

func (r *clientRepository) FindByDocument(document string) (*models.Client, error) {
	var client models.Client
	err := r.db.Where("document = ?", document).First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *clientRepository) FindByName(name string) ([]models.Client, error) {
	var clients []models.Client
	err := r.db.Where("name ILIKE ?", "%"+name+"%").Find(&clients).Error
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (r *clientRepository) ListClients() ([]models.Client, error) {
	log.Println("Método ListClients do repo chamado")

	var clients []models.Client
	result := r.db.Order("name ASC").Find(&clients)
	if result.Error != nil {
		log.Printf("erro no banco de dados: %v", result.Error)
		return nil, result.Error
	}

	log.Printf("query banco de dados retornou %d linhas", result.RowsAffected)

	if result.RowsAffected == 0 {
		log.Println("nenhum registro encontrado no banco de dados")
		return []models.Client{}, nil
	}

	// Log primeiro cliente como exemplo
	if len(clients) > 0 {
		log.Printf("exemplo de cliente: ID=%d, Nome=%s", clients[0].ID, clients[0].Name)
	}

	return clients, nil
}

func (r *clientRepository) BlockClient(document string) error {
	result := r.db.Model(&models.Client{}).
		Where("document = ?", document).
		Update("blocked", true)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("cliente n encontrado")
	}
	return nil
}

func (r *clientRepository) UnblockClient(document string) error {
	result := r.db.Model(&models.Client{}).
		Where("document = ?", document).
		Update("blocked", false)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("cliente n encontrado")
	}
	return nil
}
