package services

import (
	"auth-server/mapper"
	"auth-server/models"
	"auth-server/repository"
	"context"

	"github.com/google/uuid"
)

type ApplicationService struct {
	repo *repository.ApplicationRepository
}

func NewApplicationService(repo *repository.ApplicationRepository) *ApplicationService {
	return &ApplicationService{repo: repo}
}

func (s *ApplicationService) GetAll() ([]*models.ApplicationDto, error) {
	apps, err := s.repo.FindAll(context.Background())
	if err != nil {
		return nil, err
	}
	return mapper.ApplicationsToApplicationDtos(apps), nil
}

func (s *ApplicationService) GetById(id string) (*models.ApplicationDto, error) {
	app, err := s.repo.FindById(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return mapper.ApplicationToApplicationDto(app), nil
}

func (s *ApplicationService) CreateApplication(app *models.ApplicationDto) (*models.ApplicationDto, error) {
	appModel := &models.Application{
		BaseUUIDEntity: models.BaseUUIDEntity{
			ID: uuid.New(),
		},
		AppName: app.AppName,
	}
	appModel, err := s.repo.Save(context.Background(), appModel)
	if err != nil {
		return nil, err
	}
	return mapper.ApplicationToApplicationDto(appModel), nil
}
