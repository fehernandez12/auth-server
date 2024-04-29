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

type RoleRequest struct {
	Name string `json:"name"`
}

type RoleDto struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
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

type ClientRequest struct {
	ClientName string `json:"name"`
}

type ApplicationRequest struct {
	AppName string `json:"name"`
}
