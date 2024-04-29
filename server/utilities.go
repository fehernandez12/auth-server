package server

import "net/http"

// getStatusCode returns the appropriate successful status code
// for the provided HTTP method.
func (s *Server) getStatusCode(method string) int {
	switch method {
	case http.MethodGet:
		return http.StatusOK
	case http.MethodPost:
		return http.StatusCreated
	case http.MethodPut:
		return http.StatusAccepted
	case http.MethodDelete:
		return http.StatusNoContent
	default:
		return http.StatusOK
	}
}
