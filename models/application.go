package models

import (
	"encoding/json"
)

type Application struct {
	BaseUUIDEntity
	AppName string `json:"name" gorm:"unique"`
}

func (a Application) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}
