package repository

import (
	"auth-server/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{
		db: db,
	}
}

func (p *PermissionRepository) FindAll(ctx context.Context) ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := p.db.WithContext(ctx).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (p *PermissionRepository) FindById(ctx context.Context, id string) (*models.Permission, error) {
	return nil, errors.New("not implemented")
}

func (p *PermissionRepository) FindByName(ctx context.Context, name string) (*models.Permission, error) {
	var permission models.Permission
	err := p.db.WithContext(ctx).Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (p *PermissionRepository) Save(ctx context.Context, entity interface{}) (*models.Permission, error) {
	permission := entity.(*models.Permission)
	err := p.db.WithContext(ctx).Save(permission).Error
	if err != nil {
		return nil, err
	}
	var savedPermission models.Permission
	p.db.WithContext(ctx).Where("id = ?", permission.ID).First(&savedPermission)
	return &savedPermission, nil
}

func (p *PermissionRepository) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}
