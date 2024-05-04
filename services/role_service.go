package services

import (
	"auth-server/mapper"
	"auth-server/models"
	"auth-server/repository"
	"context"

	"github.com/google/uuid"
)

type RoleService struct {
	repo *repository.RoleRepository
}

func NewRoleService(repo *repository.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) GetAll() ([]*models.RoleDto, error) {
	roles, err := s.repo.FindAll(context.Background())
	if err != nil {
		return nil, err
	}
	return mapper.RolesToRoleDtos(roles), nil
}

func (s *RoleService) GetById(id string) (*models.RoleDto, error) {
	roleId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	role, err := s.repo.FindById(context.Background(), roleId.String())
	if err != nil {
		return nil, err
	}
	return mapper.RoleToRoleDto(role), nil
}

func (s *RoleService) GetRoleById(id string) (*models.Role, error) {
	roleId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	role, err := s.repo.FindById(context.Background(), roleId.String())
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) CreateRole(role *models.RoleRequest, app *models.ApplicationDto) (*models.RoleDto, error) {
	appId, err := uuid.Parse(app.ID)
	if err != nil {
		return nil, err
	}
	roleModel := &models.Role{
		BaseUUIDEntity: models.BaseUUIDEntity{
			ID: uuid.New(),
		},
		Name:          role.Name,
		ApplicationID: appId,
	}
	roleModel, err = s.repo.Save(context.Background(), roleModel)
	if err != nil {
		return nil, err
	}
	return mapper.RoleToRoleDto(roleModel), nil
}
