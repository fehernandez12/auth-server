package services

import (
	"auth-server/mapper"
	"auth-server/models"
	"auth-server/repository"
	"context"

	"github.com/google/uuid"
)

type ClientService struct {
	repo *repository.ClientRepository
}

func NewClientService(repo *repository.ClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

func (s *ClientService) GetAll() ([]*models.ClientDto, error) {
	clients, err := s.repo.FindAll(context.Background())
	if err != nil {
		return nil, err
	}
	return mapper.ClientsToClientDtos(clients), nil
}

func (s *ClientService) GetById(id string) (*models.ClientDto, error) {
	client, err := s.repo.FindById(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return mapper.ClientToClientDto(client), nil
}

func (s *ClientService) CreateClient(client *models.ClientDto) (*models.ClientDto, error) {
	clientModel := &models.Client{
		BaseUUIDEntity: models.BaseUUIDEntity{
			ID: uuid.New(),
		},
		ClientName:  client.ClientName,
		RedirectURI: client.RedirectURI,
	}
	clientModel, err := s.repo.Save(context.Background(), clientModel)
	if err != nil {
		return nil, err
	}
	return mapper.ClientToClientDto(clientModel), nil
}
