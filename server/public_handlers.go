package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	resp, err := json.Marshal(map[string]string{"status": "OK"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	w.Write(resp)
}
