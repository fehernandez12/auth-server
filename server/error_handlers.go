package server

import (
	"auth-server/models"
	"encoding/json"
	"net/http"
)

func (s *Server) HandleError(w http.ResponseWriter, statusCode int, route string, cause error) {
	var errorResponse models.ErrorResponse
	errorResponse.Messages = append(errorResponse.Messages, cause.Error())
	response, err := json.Marshal(errorResponse)
	if err != nil {
		s.HandleError(w, http.StatusInternalServerError, route, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
	s.logger.Error(statusCode, route, cause)
}
