package services

import (
	"auth-server/mapper"
	"auth-server/models"
	"auth-server/repository"
	"context"
	"errors"
)

type UserService struct {
	repo *repository.UserRepository
}

// NewUserService creates a new instance of UserService with the provided UserRepository.
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser creates a new user with the provided user details.
func (s *UserService) CreateUser(data *models.SignupRequest) (*models.UserDto, error) {
	_user, err := s.repo.FindByEmail(context.Background(), data.Email)
	if err != nil {
		return nil, err
	}
	if _user != nil {
		return nil, errors.New("user with email already exists")
	}

	_user, err = s.repo.FindByUsername(context.Background(), data.Username)
	if err != nil {
		return nil, err
	}
	if _user != nil {
		return nil, errors.New("user with username already exists")
	}

	user := models.NewUser(data)

	result, err := s.repo.Save(context.Background(), user)
	if err != nil {
		return nil, err
	}

	return mapper.UserToUserDto(result), nil
}

// GetAll returns all users.
func (s *UserService) GetAll() ([]*models.UserDto, error) {
	users, err := s.repo.FindAll(context.Background())
	if err != nil {
		return nil, err
	}
	return mapper.UsersToUserDtos(users), nil
}

// GetById returns a user by the provided id.
func (s *UserService) GetById(id string) (*models.User, error) {
	return s.repo.FindById(context.Background(), id)
}

func (s *UserService) GetByUsername(username string) (*models.User, error) {
	return s.repo.FindByUsername(context.Background(), username)
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	return s.repo.FindByEmail(context.Background(), email)
}
