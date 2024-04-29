package server

import (
	"auth-server/models"
	"auth-server/repository"
	"auth-server/services"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
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
		user := models.NewUser(&signupRequest)
		result, err := service.CreateUser(user)
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
	var roleRequest models.RoleRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&roleRequest)
	if err != nil {
		s.HandleError(w, http.StatusBadRequest, ADMIN_ROLE_ROUTE, err)
		return
	}
	role := &models.Role{
		Name: roleRequest.Name,
	}
	repo := s.roleRepository.(*repository.RoleRepository)
	_role, _ := repo.FindByName(context.Background(), role.Name)
	if _role != nil {
		s.HandleError(w, http.StatusConflict, ADMIN_ROLE_ROUTE, errors.New("role already exists"))
		return
	}
	result, err := repo.Save(context.Background(), role)
	if err != nil {
		s.HandleError(w, http.StatusConflict, ADMIN_ROLE_ROUTE, err)
		return
	}
	response, err := json.Marshal(result)
	if err != nil {
		s.HandleError(w, http.StatusInternalServerError, ADMIN_ROLE_ROUTE, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
	s.logger.Info(http.StatusCreated, ADMIN_ROLE_ROUTE, start)
}

func (s *Server) HandlePermission(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method == http.MethodGet {
		repo := s.permissionRepository.(*repository.PermissionRepository)
		result, err := repo.FindAll(context.Background())
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_PERMISSION_ROUTE, err)
			return
		}
		response, err := json.Marshal(result)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_PERMISSION_ROUTE, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		s.logger.Info(http.StatusOK, ADMIN_PERMISSION_ROUTE, start)
		return
	}
}
func (s *Server) HandleApplication(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var ApplicationRequest models.ApplicationRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&ApplicationRequest)
	if err != nil {
		s.HandleError(w, http.StatusBadRequest, ADMIN_APPLICATION_ROUTE, err)
		return
	}
	App := &models.Application{
		BaseUUIDEntity: models.BaseUUIDEntity{
			ID: uuid.New(),
		},
		AppName: ApplicationRequest.AppName,
	}
	repo := s.applicationRepository.(*repository.ApplicationRepository)

	result, err := repo.Save(context.Background(), App)
	if err != nil {
		s.HandleError(w, http.StatusConflict, ADMIN_APPLICATION_ROUTE, err)
		return
	}
	response, err := json.Marshal(result)
	if err != nil {
		s.HandleError(w, http.StatusInternalServerError, ADMIN_APPLICATION_ROUTE, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
	s.logger.Info(http.StatusCreated, ADMIN_APPLICATION_ROUTE, start)
}
func (s *Server) HandleClient(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var ClientRequest models.ClientRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&ClientRequest)
	if err != nil {
		s.HandleError(w, http.StatusBadRequest, ADMIN_CLIENT_ROUTE, err)
		return
	}
	client := &models.Client{
		BaseUUIDEntity: models.BaseUUIDEntity{
			ID: uuid.New(),
		},
		ClientName: ClientRequest.ClientName,
	}
	repo := s.clientRepository.(*repository.ClientRepository)

	result, err := repo.Save(context.Background(), client)
	if err != nil {
		s.HandleError(w, http.StatusConflict, ADMIN_CLIENT_ROUTE, err)
		return
	}
	response, err := json.Marshal(result)
	if err != nil {
		s.HandleError(w, http.StatusInternalServerError, ADMIN_CLIENT_ROUTE, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
	s.logger.Info(http.StatusCreated, ADMIN_CLIENT_ROUTE, start)
}
