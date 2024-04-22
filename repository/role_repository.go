package repository

import (
	"auth-server/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (p *RoleRepository) FindAll(ctx context.Context) ([]*models.Role, error) {
	return nil, errors.New("not implemented")
}

func (p *RoleRepository) FindById(ctx context.Context, id string) (*models.Role, error) {
	return nil, errors.New("not implemented")
}

func (p *RoleRepository) FindByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	err := p.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (p *RoleRepository) Save(ctx context.Context, entity interface{}) (*models.Role, error) {
	role := entity.(*models.Role)
	err := p.db.WithContext(ctx).Save(role).Error
	if err != nil {
		return nil, err
	}
	var savedRole models.Role
	p.db.WithContext(ctx).Where("id = ?", role.ID).First(&savedRole)
	return &savedRole, nil
}

func (p *RoleRepository) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}
