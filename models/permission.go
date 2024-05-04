package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Name   string `json:"name"`
	RoleID uuid.UUID
	Role   *Role `json:"role"`
}
