package models

import (
	"encoding/json"
)

type Client struct {
	BaseUUIDEntity
	ClientName string `json:"name" gorm:"unique"`
}

func (c Client) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}
