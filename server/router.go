package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route constants
const (
	TOKEN_ROUTE              = "/token/"
	ADMIN_CLIENT_ROUTE       = "/admin/client/"
	ADMIN_USER_ROUTE         = "/admin/user/"
	ADMIN_USER_DETAILS_ROUTE = "/admin/user/{username}/"
	ADMIN_USER_ROLES_ROUTE   = "/admin/user/:username/roles/"
	ADMIN_ROLE_ROUTE         = "/admin/role/"
	ADMIN_PERMISSION_ROUTE   = "/admin/permission/"
	ADMIN_APPLICATION_ROUTE  = "/admin/application/"
)

func (s *Server) router() http.Handler {
	// Base Router
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.Use(s.logger.RequestLoggerMiddleware)

	// Admin Router
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(s.AuthMiddleware)
	adminRouter.HandleFunc("/user/", s.HandleUser).Methods(http.MethodGet, http.MethodPost)
	adminRouter.HandleFunc("/user/{username}/", s.HandleUserDetails).Methods(http.MethodGet)
	adminRouter.HandleFunc("/user/{username}/roles/", s.HandleUserRoles).Methods(http.MethodGet, http.MethodPost, http.MethodPatch)
	adminRouter.HandleFunc("/role/", s.HandleRole).Methods(http.MethodGet, http.MethodPost)
	adminRouter.HandleFunc("/permission/", s.HandlePermission).Methods(http.MethodGet, http.MethodPost)
	adminRouter.HandleFunc("/client/", s.HandleClient).Methods(http.MethodGet, http.MethodPost)
	adminRouter.HandleFunc("/application/", s.HandleApplication).Methods(http.MethodGet, http.MethodPost)

	// Public Router
	publicRouter := router.PathPrefix("/public").Subrouter()
	publicRouter.HandleFunc("/health/", s.healthHandler).Methods("GET")

	// OAuth2 Router
	// oauth2Router := router.PathPrefix("/oauth2").Subrouter()
	// oauth2Router.HandleFunc(TOKEN_ROUTE, s.HandleToken).Methods("POST")
	// oauth2Router.HandleFunc("/introspect", s.HandleIntrospection).Methods("POST")
	// oauth2Router.HandleFunc("/tokeninfo", s.HandleTokenInfo).Methods("GET")

	// Admin-only routes. They require s.AuthMiddleware.
	return router
}
