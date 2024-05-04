package services

import (
	"auth-server/mapper"
	"auth-server/models"
	"auth-server/repository"
	"context"

	"github.com/google/uuid"
)

type Permissionservice struct {
	repo *repository.PermissionRepository
}

func NewPermissionService(repo *repository.PermissionRepository) *Permissionservice {
	return &Permissionservice{repo: repo}
}

func (s *Permissionservice) GetAll() ([]*models.PermissionDto, error) {
	permissions, err := s.repo.FindAll(context.Background())
	if err != nil {
		return nil, err
	}
	return mapper.PermissionsToPermissionDtos(permissions), nil
}

func (s *Permissionservice) CreatePermission(permission *models.PermissionRequest) (*models.PermissionDto, error) {
	roleId, err := uuid.Parse(permission.RoleID)
	if err != nil {
		return nil, err
	}
	permissionModel := &models.Permission{
		Name:   permission.Name,
		RoleID: roleId,
	}
	permissionModel, err = s.repo.Save(context.Background(), permissionModel)
	if err != nil {
		return nil, err
	}
	return mapper.PermissionToPermissionDto(permissionModel), nil
}
