package repository

import (
	"auth-server/models"
	"context"
	"log"

	"gorm.io/gorm"
)

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

func (p *ClientRepository) FindAll(ctx context.Context) ([]*models.Client, error) {
	var clients []*models.Client
	err := p.db.WithContext(ctx).Find(&clients).Error
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (p *ClientRepository) FindById(ctx context.Context, id string) (*models.Client, error) {
	var client models.Client
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&client).Error
	log.Println(err)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (p *ClientRepository) Save(ctx context.Context, entity interface{}) (*models.Client, error) {
	client := entity.(*models.Client)
	err := p.db.WithContext(ctx).Save(client).Error
	if err != nil {
		return nil, err
	}
	p.db.WithContext(ctx).Where("id = ?", client.ID).First(&client)
	return client, nil
}

func (p *ClientRepository) Delete(ctx context.Context, id string) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Client{}).Error
	if err != nil {
		return err
	}
	return nil
}
