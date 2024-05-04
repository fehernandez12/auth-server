package server

import (
	"auth-server/mapper"
	"auth-server/models"
	"auth-server/repository"
	"auth-server/services"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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
			s.HandleError(w, http.StatusInternalServerError, ADMIN_USER_ROUTE, err)
			return
		}
		response, err = json.Marshal(result)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_USER_ROUTE, err)
			return
		}
	case http.MethodPost:
		var signupRequest models.SignupRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&signupRequest)
		if err != nil {
			s.HandleError(w, http.StatusBadRequest, ADMIN_USER_ROUTE, err)
			return
		}
		password, err := s.hasher.GenerateFromPassword(signupRequest.Password)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_USER_ROUTE, err)
			return
		}
		signupRequest.Password = password
		result, err := service.CreateUser(&signupRequest)
		if err != nil {
			s.HandleError(w, http.StatusConflict, ADMIN_USER_ROUTE, err)
			return
		}
		response, err = json.Marshal(result)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_USER_ROUTE, err)
			return
		}
	}

	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	status := s.getStatusCode(r.Method)
	w.WriteHeader(status)
	w.Write(response)
	s.logger.Info(status, ADMIN_USER_ROUTE, start)
}

// HandleUserDetails retrieves user details by username.
func (s *Server) HandleUserDetails(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	vars := mux.Vars(r)
	repo := s.userRepository.(*repository.UserRepository)
	service := services.NewUserService(repo)
	var response []byte

	user, err := service.GetByUsername(vars["username"])
	if user == nil && err == nil {
		s.HandleError(w, http.StatusNotFound, ADMIN_USER_DETAILS_ROUTE, err)
		return
	}
	if err != nil {
		s.HandleError(w, http.StatusInternalServerError, ADMIN_USER_DETAILS_ROUTE, err)
		return
	}

	result := mapper.UserToUserDto(user)

	response, err = json.Marshal(result)
	if err != nil {
		s.HandleError(w, http.StatusInternalServerError, ADMIN_USER_DETAILS_ROUTE, err)
		return
	}

	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	status := s.getStatusCode(r.Method)
	w.WriteHeader(status)
	w.Write(response)
	s.logger.Info(status, ADMIN_USER_DETAILS_ROUTE, start)
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
		var data models.RoleRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)
		if err != nil {
			s.HandleError(w, http.StatusBadRequest, ADMIN_ROLE_ROUTE, err)
			return
		}

		appRepo := s.applicationRepository.(*repository.ApplicationRepository)
		appService := services.NewApplicationService(appRepo)

		app, err := appService.GetById(data.ApplicationId)
		if err != nil {
			s.HandleError(w, http.StatusNotFound, ADMIN_ROLE_ROUTE, err)
			return
		}

		result, err := service.CreateRole(&data, app)
		if err != nil {
			s.HandleError(w, http.StatusConflict, ADMIN_ROLE_ROUTE, err)
			return
		}
		response, err = json.Marshal(result)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_ROLE_ROUTE, err)
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
		var permissionRequest models.PermissionRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&permissionRequest)
		if err != nil {
			s.HandleError(w, http.StatusBadRequest, ADMIN_PERMISSION_ROUTE, err)
			return
		}

		result, err := service.CreatePermission(&permissionRequest)
		if err != nil {
			s.HandleError(w, http.StatusConflict, ADMIN_PERMISSION_ROUTE, err)
			return
		}

		response, err = json.Marshal(result)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_PERMISSION_ROUTE, err)
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
	repo := s.clientRepository.(*repository.ClientRepository)
	service := services.NewClientService(repo)
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
		var clientRequest models.ClientDto
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&clientRequest)
		if err != nil {
			s.HandleError(w, http.StatusBadRequest, ADMIN_CLIENT_ROUTE, err)
			return
		}

		client, err := service.CreateClient(&clientRequest)
		if err != nil {
			s.HandleError(w, http.StatusConflict, ADMIN_CLIENT_ROUTE, err)
			return
		}

		response, err = json.Marshal(client)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_CLIENT_ROUTE, err)
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
	repo := s.applicationRepository.(*repository.ApplicationRepository)
	service := services.NewApplicationService(repo)
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
		var applicationRequest models.ApplicationRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&applicationRequest)
		if err != nil {
			s.HandleError(w, http.StatusBadRequest, ADMIN_APPLICATION_ROUTE, err)
			return
		}

		app, err := service.CreateApplication(&models.ApplicationDto{AppName: applicationRequest.AppName})
		if err != nil {
			s.HandleError(w, http.StatusConflict, ADMIN_APPLICATION_ROUTE, err)
			return
		}

		response, err = json.Marshal(app)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_APPLICATION_ROUTE, err)
			return
		}
	}

	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	status := s.getStatusCode(r.Method)
	w.WriteHeader(status)
	w.Write(response)
	s.logger.Info(status, ADMIN_APPLICATION_ROUTE, start)
}

// HandleUserRoles handles user roles retrieval and assignment. This handler works with
// GET, POST, and PATCH methods. When called via GET, it retrieves all roles assigned to a user.
// When called via POST, it assigns roles to a user. When called via PATCH, it updates the roles
// assigned to a user.
func (s *Server) HandleUserRoles(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	vars := mux.Vars(r)
	repo := s.userRepository.(*repository.UserRepository)
	service := services.NewUserService(repo)
	var response []byte

	user, err := service.GetByUsername(vars["username"])
	if user == nil && err == nil {
		s.HandleError(w, http.StatusNotFound, ADMIN_USER_ROLES_ROUTE, err)
		return
	}
	if err != nil {
		s.HandleError(w, http.StatusInternalServerError, ADMIN_USER_ROLES_ROUTE, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		result, err := service.GetUserRoles(user)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_USER_ROLES_ROUTE, err)
			return
		}
		response, err = json.Marshal(result)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_USER_ROLES_ROUTE, err)
			return
		}

	case http.MethodPost:
		result, err := s.handleRoleAssignment(r, user, service)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_USER_ROLES_ROUTE, err)
			return
		}

		response, err = json.Marshal(result)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_USER_ROLES_ROUTE, err)
			return
		}

	case http.MethodPatch:
		result, err := s.handleRoleUpdate(r, user, service)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_USER_ROLES_ROUTE, err)
			return
		}

		response, err = json.Marshal(result)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, ADMIN_USER_ROLES_ROUTE, err)
			return
		}
	}

	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	status := s.getStatusCode(r.Method)
	w.WriteHeader(status)
	w.Write(response)
	s.logger.Info(status, ADMIN_APPLICATION_ROUTE, start)
}

func (s *Server) handleRoleAssignment(r *http.Request, user *models.User, service *services.UserService) ([]*models.RoleDto, error) {
	var rolesRequest models.UserRolesRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rolesRequest)
	if err != nil {
		return nil, err
	}

	roleRepo := s.roleRepository.(*repository.RoleRepository)
	roleService := services.NewRoleService(roleRepo)

	var roleEntities []*models.Role
	for _, roleId := range rolesRequest.Roles {
		role, err := roleService.GetRoleById(roleId)
		if err != nil {
			return nil, err
		}

		entity, err := roleService.GetRoleById(role.ID.String())
		if err != nil {
			return nil, err
		}
		roleEntities = append(roleEntities, entity)
	}

	result, err := service.AssignRolesToUser(user, roleEntities)
	if err != nil {
		return nil, err
	}

	return result.Roles, nil
}

func (s *Server) handleRoleUpdate(r *http.Request, user *models.User, service *services.UserService) ([]*models.RoleDto, error) {
	var rolesRequest models.UserRolesRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rolesRequest)
	if err != nil {
		return nil, err
	}

	roleRepo := s.roleRepository.(*repository.RoleRepository)
	roleService := services.NewRoleService(roleRepo)

	var roleEntities []*models.Role
	for _, role := range rolesRequest.Roles {
		role, err := roleService.GetRoleById(role)
		if err != nil {
			return nil, err
		}

		entity, err := roleService.GetRoleById(role.ID.String())
		if err != nil {
			return nil, err
		}
		roleEntities = append(roleEntities, entity)
	}

	result, err := service.AddRolesToUser(user, roleEntities)
	if err != nil {
		return nil, err
	}

	return result.Roles, nil
}
