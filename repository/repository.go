package repository

import (
	"context"
)

type Entity interface{}

type Repository[T Entity] interface {
	FindAll(ctx context.Context) ([]*T, error)
	FindById(ctx context.Context, id string) (*T, error)
	Save(ctx context.Context, entity interface{}) (*T, error)
	Delete(ctx context.Context, id string) error
}
