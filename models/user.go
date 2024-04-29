package models

import (
	"github.com/google/uuid"
)

type User struct {
	BaseUUIDEntity
	Email                 string `json:"email" gorm:"unique"`
	Username              string `json:"username" gorm:"unique"`
	Password              string `json:"password"`
	Enabled               bool   `json:"enabled"`
	AccountNonLocked      bool   `json:"account_non_locked"`
	AccountNonExpired     bool   `json:"account_non_expired"`
	CredentialsNonExpired bool   `json:"credentials_non_expired"`
	Roles                 []Role `json:"roles" gorm:"many2many:user_roles;"`
}

// NewUser creates a new user from a SignupRequest.
func NewUser(dto *SignupRequest) *User {
	return &User{
		BaseUUIDEntity: BaseUUIDEntity{
			ID: uuid.New(),
		},
		Username:              dto.Username,
		Email:                 dto.Email,
		Password:              dto.Password,
		Enabled:               true,
		AccountNonLocked:      true,
		AccountNonExpired:     true,
		CredentialsNonExpired: true,
	}
}

func (u *User) Authorities() []string {
	var authorities []string
	for _, role := range u.Roles {
		authorities = append(authorities, role.Name)
		for _, permission := range role.Permissions {
			authorities = append(authorities, permission.Name)
		}
	}
	return authorities
}
