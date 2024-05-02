package services

import (
	"auth-server/mapper"
	"auth-server/models"
	"auth-server/repository"
	"context"
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
