package services

import (
	"auth-server/mapper"
	"auth-server/models"
	"auth-server/repository"
	"context"
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
