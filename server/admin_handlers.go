package server

import (
	"auth-server/models"
	"auth-server/repository"
	"auth-server/services"
	"encoding/json"
	"net/http"
	"time"
)

// HandleUser handles user creation and retrieval. When called via POST,
// it creates a new user, expecting it as a SignupRequest.
// When called via GET, it retrieves all users.
func (s *Server) HandleUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	repo := s.userRepository.(*repository.UserRepository)
	service := services.NewUserService(repo)
	var response []byte
	switch r.Method {
	case http.MethodGet:
		result, err := service.GetAll()
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_REGISTER_ROUTE, err)
			return
		}
		response, err = json.Marshal(result)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_REGISTER_ROUTE, err)
			return
		}
	case http.MethodPost:
		var signupRequest models.SignupRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&signupRequest)
		if err != nil {
			s.HandleError(w, http.StatusBadRequest, ADMIN_REGISTER_ROUTE, err)
			return
		}
		password, err := s.hasher.GenerateFromPassword(signupRequest.Password)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_REGISTER_ROUTE, err)
			return
		}
		signupRequest.Password = password
		result, err := service.CreateUser(&signupRequest)
		if err != nil {
			s.HandleError(w, http.StatusConflict, ADMIN_REGISTER_ROUTE, err)
			return
		}
		response, err = json.Marshal(result)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_REGISTER_ROUTE, err)
			return
		}
	}
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	status := s.getStatusCode(r.Method)
	w.WriteHeader(status)
	w.Write(response)
	s.logger.Info(status, ADMIN_REGISTER_ROUTE, start)
}

func (s *Server) HandleRole(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	repo := s.roleRepository.(*repository.RoleRepository)
	service := services.NewRoleService(repo)
	var response []byte
	switch r.Method {
	case http.MethodGet:
		result, err := service.GetAll()
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_ROLE_ROUTE, err)
			return
		}
		response, err = json.Marshal(result)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_ROLE_ROUTE, err)
			return
		}
	case http.MethodPost:
		var Request models.RoleRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&Request)
		if err != nil {
			s.HandleError(w, http.StatusBadRequest, ADMIN_ROLE_ROUTE, err)
			return
		}
	}
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	status := s.getStatusCode(r.Method)
	w.WriteHeader(status)
	w.Write(response)
	s.logger.Info(status, ADMIN_ROLE_ROUTE, start)
}
func (s *Server) HandlePermission(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	repo := s.permissionRepository.(*repository.PermissionRepository)
	service := services.NewPermissionService(repo)
	var response []byte
	switch r.Method {
	case http.MethodGet:
		result, err := service.GetAll()
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_PERMISSION_ROUTE, err)
			return
		}
		response, err = json.Marshal(result)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_PERMISSION_ROUTE, err)
			return
		}
	case http.MethodPost:
		var ApplicationRequest models.ApplicationRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&ApplicationRequest)
		if err != nil {
			s.HandleError(w, http.StatusBadRequest, ADMIN_PERMISSION_ROUTE, err)
			return
		}

	}
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	status := s.getStatusCode(r.Method)
	w.WriteHeader(status)
	w.Write(response)
	s.logger.Info(status, ADMIN_PERMISSION_ROUTE, start)
}

func (s *Server) HandleClient(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	repo := s.userRepository.(*repository.UserRepository)
	service := services.NewUserService(repo)
	var response []byte
	switch r.Method {
	case http.MethodGet:
		result, err := service.GetAll()
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_CLIENT_ROUTE, err)
			return
		}
		response, err = json.Marshal(result)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_CLIENT_ROUTE, err)
			return
		}
	case http.MethodPost:
		var ClientRequest models.ClientRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&ClientRequest)
		if err != nil {
			s.HandleError(w, http.StatusBadRequest, ADMIN_CLIENT_ROUTE, err)
			return
		}

	}
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	status := s.getStatusCode(r.Method)
	w.WriteHeader(status)
	w.Write(response)
	s.logger.Info(status, ADMIN_CLIENT_ROUTE, start)
}
func (s *Server) HandleApplication(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	repo := s.userRepository.(*repository.UserRepository)
	service := services.NewUserService(repo)
	var response []byte
	switch r.Method {
	case http.MethodGet:
		result, err := service.GetAll()
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_APPLICATION_ROUTE, err)
			return
		}
		response, err = json.Marshal(result)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_APPLICATION_ROUTE, err)
			return
		}
	case http.MethodPost:
		var ApplicationRequest models.ApplicationRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&ApplicationRequest)
		if err != nil {
			s.HandleError(w, http.StatusBadRequest, ADMIN_APPLICATION_ROUTE, err)
			return
		}
	}
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	status := s.getStatusCode(r.Method)
	w.WriteHeader(status)
	w.Write(response)
	s.logger.Info(status, ADMIN_APPLICATION_ROUTE, start)
}
