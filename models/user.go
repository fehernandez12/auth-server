package models

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
