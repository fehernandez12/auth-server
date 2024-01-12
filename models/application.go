package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Application struct {
	ID      uuid.UUID `json:"id" gorm:"primaryKey"`
	AppName string    `json:"name" gorm:"unique"`
}

func (a Application) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}
