package models

import "github.com/google/uuid"

type Role struct {
	BaseUUIDEntity
	Name          string `json:"name"`
	ApplicationID uuid.UUID
	Application   *Application `json:"application"`
}
