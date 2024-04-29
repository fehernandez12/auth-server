package repository

import (
	"auth-server/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (p *UserRepository) FindAll(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	err := p.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (p *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := p.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := p.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *UserRepository) FindById(ctx context.Context, id string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (p *UserRepository) Save(ctx context.Context, entity interface{}) (*models.User, error) {
	user := entity.(*models.User)
	err := p.db.WithContext(ctx).Save(user).Error
	if err != nil {
		return nil, err
	}
	var savedUser models.User
	p.db.WithContext(ctx).Where("id = ?", user.ID).First(&savedUser)
	return &savedUser, nil
}

func (p *UserRepository) AddRolesToUser(ctx context.Context, user *models.User, roles []*models.Role) error {
	err := p.db.WithContext(ctx).Model(user).Association("Roles").Append(roles)
	if err != nil {
		return err
	}
	return nil
}

func (p *UserRepository) Delete(ctx context.Context, id string) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&models.User{}).Error
	if err != nil {
		return err
	}
	return nil
}
