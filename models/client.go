package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	pq "github.com/lib/pq"
)

type Client struct {
	ID            uuid.UUID      `json:"id" gorm:"primaryKey"`
	ClientName    string         `json:"name" gorm:"unique"`
	Email         string         `json:"email" gorm:"unique"`
	AllowedScopes pq.StringArray `json:"allowed_scopes" gorm:"type:text[]"`
}

func (c Client) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Client) HasAllowedScopes(scope, appName string) bool {
	requestedScopes := strings.Split(scope, " ")
	for _, requestedScope := range requestedScopes {
		scopeAllowed := false
		for _, allowedScope := range c.AllowedScopes {
			app := strings.Split(allowedScope, ":")[0]
			permission := strings.Split(allowedScope, ":")[1]
			if app == appName && permission == strings.Split(requestedScope, ":")[1] {
				scopeAllowed = true
				break
			}
		}
		// If a requested scope is not allowed, return false
		if !scopeAllowed {
			return false
		}
	}
	// If all requested scopes are allowed, return true
	return true
}

func (c *Client) CheckApps(appNames []string) error {
	for _, allowedScope := range c.AllowedScopes {
		app := strings.Split(allowedScope, ":")[0]
		if !contains(appNames, app) {
			return fmt.Errorf("invalid app name: %s", app)
		}
	}
	return nil
}

func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}
