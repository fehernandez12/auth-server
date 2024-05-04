package models

type SignupRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserDto struct {
	ID       string     `json:"id"`
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Created  string     `json:"created_at"`
	Updated  string     `json:"updated_at"`
	Enabled  bool       `json:"enabled"`
	Roles    []*RoleDto `json:"roles"`
}

type UserRolesRequest struct {
	Roles []string `json:"roles"`
}

type RoleRequest struct {
	Name          string `json:"name"`
	ApplicationId string `json:"application_id"`
}

type RoleDto struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Application *ApplicationDto `json:"application"`
}

type PermissionRequest struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	RoleID string `json:"role_id"`
}

type PermissionDto struct {
	ID   uint     `json:"id"`
	Name string   `json:"name"`
	Role *RoleDto `json:"role"`
}

type ApplicationRequest struct {
	AppName string `json:"name"`
}

type ApplicationDto struct {
	ID      string `json:"id"`
	AppName string `json:"name"`
}

type ClientDto struct {
	ID          string `json:"id"`
	ClientName  string `json:"name"`
	RedirectURI string `json:"redirect_uri"`
}

type TokenRequest struct {
	GrantType string `json:"grant_type"`
	ClientId  string `json:"client_id"`
	Aud       string `json:"aud"`
	Scope     string `json:"scope"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type IntrospectionRequest struct {
	OpaqueToken string `json:"opaque_token"`
}

type IntrospectionResponse struct {
	Active bool `json:"active"`
}

type TokenInfoRequest struct {
	AccessToken string `json:"access_token"`
}

type TokenInfoResponse struct {
	Active bool `json:"active"`
}

type ErrorResponse struct {
	Messages []string `json:"message"`
}
