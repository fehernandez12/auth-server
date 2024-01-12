package server

import (
	"errors"
	"net/http"
)

// AuthMiddleware is a middleware that checks if the request is authenticated.
// If the request is authenticated, the middleware will call the next handler.
// If the request is not authenticated, the middleware will return a 401 Unauthorized response.
//
// The request is verified via the Authorization header.
// This header must contain a JWT Bearer token, which is validated with the server itself.
func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			s.HandleError(w, http.StatusUnauthorized, "AuthMiddleware", errors.New("missing authorization header"))
			return
		}
		token := auth[7:]
		if token == "" {
			s.HandleError(w, http.StatusUnauthorized, "AuthMiddleware", errors.New("missing bearer token"))
			return
		}
		_, err := s.ValidateToken(token)
		if err != nil {
			s.HandleError(w, http.StatusUnauthorized, "AuthMiddleware", err)
			return
		}
		next.ServeHTTP(w, r)
	})
}
