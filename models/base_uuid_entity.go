package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseUUIDEntity struct {
	ID        uuid.UUID  `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}
