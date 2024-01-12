package repository

import (
	"auth-server/models"
	"context"

	"gorm.io/gorm"
)

type ApplicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{
		db: db,
	}
}

func (p *ApplicationRepository) FindAll(ctx context.Context) ([]*models.Application, error) {
	var applications []*models.Application
	err := p.db.WithContext(ctx).Find(&applications).Error
	if err != nil {
		return nil, err
	}
	return applications, nil
}

func (p *ApplicationRepository) FindById(ctx context.Context, id string) (*models.Application, error) {
	var application models.Application
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&application).Error
	if err != nil {
		return nil, err
	}
	return &application, nil
}

func (p *ApplicationRepository) Save(ctx context.Context, entity interface{}) (*models.Application, error) {
	application := entity.(*models.Application)
	err := p.db.WithContext(ctx).Save(application).Error
	if err != nil {
		return nil, err
	}
	return application, nil
}

func (p *ApplicationRepository) Delete(ctx context.Context, id string) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Application{}).Error
	if err != nil {
		return err
	}
	return nil
}
