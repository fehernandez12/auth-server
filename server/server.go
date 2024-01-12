package server

import (
	"auth-server/models"
	"auth-server/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	TOKEN_ROUTE        = "/token"
	ADMIN_CLIENT_ROUTE = "/admin/client"
)

type ServerConfig struct {
	Timeout int
	Addr    string
}

type Server struct {
	config                *ServerConfig
	DB                    *gorm.DB
	clientRepository      repository.Repository[models.Client]
	applicationRepository repository.Repository[models.Application]
	logger                *Logger
}

func NewServer() (*Server, error) {
	s := &Server{
		logger: NewLogger(),
	}
	banner, err := os.ReadFile("server/banner.txt")
	if err != nil {
		s.logger.Fatal(err)
	}
	fmt.Println(string(banner))
	s.logger.WithField("Status", "Loading server config...")
	config, err := s.readServerConfig()
	if err != nil {
		s.logger.Fatal(err)
	}
	s.config = config
	s.logger.WithField("Status", "Loading database connection...")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		s.logger.Fatal(err)
	}
	s.logger.WithField("Status", "Migrating database...")
	err = db.AutoMigrate(&models.Client{}, &models.Application{})
	if err != nil {
		s.logger.Fatal(err)
	}
	s.DB = db
	s.clientRepository = repository.NewClientRepository(db)
	s.applicationRepository = repository.NewApplicationRepository(db)
	s.logger.WithField("Status", "Application is running")
	return s, nil
}

func (s *Server) Start(stop <-chan struct{}) error {
	corsObj := handlers.CORS(handlers.AllowedOrigins([]string{"*"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}))
	srv := &http.Server{
		Addr:    s.config.Addr,
		Handler: corsObj(s.router()),
	}
	go func() {
		s.logger.WithField("addr", s.config.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal(err)
		}
	}()
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.Timeout)*time.Millisecond)
	defer cancel()
	s.logger.WithField("timeout", s.config.Timeout)
	return srv.Shutdown(ctx)
}

func (s *Server) router() http.Handler {
	router := mux.NewRouter()
	router.Use(s.logger.RequestLoggerMiddleware)
	router.HandleFunc("/health", s.healthHandler).Methods("GET")
	router.HandleFunc(TOKEN_ROUTE, s.HandleToken).Methods("POST")
	router.HandleFunc("/introspect", s.HandleIntrospection).Methods("POST")
	router.HandleFunc("/tokeninfo", s.HandleTokenInfo).Methods("GET")
	// Admin-only routes. They require s.AuthMiddleware.
	authRouter := router.PathPrefix("/admin").Subrouter()
	authRouter.Use(s.AuthMiddleware)
	authRouter.HandleFunc("/client", s.HandleClient).Methods("POST")
	return router
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) readServerConfig() (*ServerConfig, error) {
	addr := os.Getenv("AUTH_SERVER_ADDR")
	if addr == "" {
		s.logger.Fatal(errors.New("addr cannot be blank. Please set the AUTH_SERVER_ADDR environment variable"))
	}
	timeout, err := strconv.ParseInt(os.Getenv("AUTH_SERVER_TIMEOUT"), 10, 64)
	if err != nil {
		s.logger.Fatal(err)
	}
	return &ServerConfig{
		Addr:    addr,
		Timeout: int(timeout),
	}, nil
}

func (s *Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	decoder := json.NewDecoder(r.Body)
	var tokenRequest models.TokenRequest
	err := decoder.Decode(&tokenRequest)
	if err != nil {
		s.HandleError(w, http.StatusBadRequest, TOKEN_ROUTE, err)
	}
	if tokenRequest.GrantType == "client_credentials" {
		ctx := context.Background()
		clientData, appData, err := s.FetchClientAndApplication(ctx, tokenRequest.ClientId, tokenRequest.Aud)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, TOKEN_ROUTE, err)
			return
		}
		s.logger.WithField("clientData", clientData)
		s.logger.WithField("appData", appData)

		if !clientData.HasAllowedScopes(tokenRequest.Scope, appData.AppName) {
			s.HandleError(w, http.StatusUnauthorized, TOKEN_ROUTE, errors.New("client does not have the requested scopes"))
			return
		}

		payload := models.NewPayload(clientData.ID.String(), appData.ID.String(), 1, tokenRequest.Scope)
		jwt, err := models.NewJwt(payload, "JWT")
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, TOKEN_ROUTE, err)
			return
		}
		var tokenResponse models.TokenResponse
		jwtToken, err := jwt.Token()
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, TOKEN_ROUTE, err)
			return
		}
		tokenResponse.AccessToken = jwtToken
		tokenResponse.TokenType = "Bearer"
		response, err := json.Marshal(tokenResponse)
		if err != nil {
			s.HandleError(w, http.StatusInternalServerError, TOKEN_ROUTE, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		s.logger.Info(http.StatusOK, TOKEN_ROUTE, start)
	}
}

func (s *Server) FetchClientAndApplication(ctx context.Context, clientId string, applicationId string) (*models.Client, *models.Application, error) {
	client, err := s.clientRepository.FindById(ctx, clientId)
	if err != nil {
		return nil, nil, err
	}

	application, err := s.applicationRepository.FindById(ctx, applicationId)
	if err != nil {
		return nil, nil, err
	}

	return client, application, nil
}

func (s *Server) HandleIntrospection(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) HandleTokenInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) HandleClient(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var clientRequest models.ClientRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&clientRequest)
	if err != nil {
		s.HandleError(w, http.StatusBadRequest, ADMIN_CLIENT_ROUTE, err)
		return
	}
	client := &models.Client{
		ID:            uuid.New(),
		ClientName:    clientRequest.ClientName,
		Email:         clientRequest.Email,
		AllowedScopes: strings.Split(clientRequest.Scopes, " "),
	}
	result, err := s.clientRepository.Save(context.Background(), client)
	if err != nil {
		s.HandleError(w, http.StatusConflict, ADMIN_CLIENT_ROUTE, err)
		return
	}
	applications, err := s.applicationRepository.FindAll(context.Background())
	if err != nil {
		s.HandleError(w, http.StatusInternalServerError, ADMIN_CLIENT_ROUTE, err)
		return
	}
	appNames := make([]string, len(applications))
	for i, app := range applications {
		appNames[i] = app.AppName
	}
	err = client.CheckApps(appNames)
	if err != nil {
		s.HandleError(w, http.StatusBadRequest, ADMIN_CLIENT_ROUTE, err)
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

func (s *Server) ValidateToken(token string) (*models.Payload, error) {
	parts := models.SplitToken(token)
	header, err := models.ParseHeader(parts[0])
	if err != nil {
		return nil, err
	}
	s.logger.WithField("header", header)
	if header.Alg != os.Getenv("AUTH_SERVER_JWT_ALG") {
		return nil, errors.New("invalid algorithm")
	}
	payload, err := models.ParsePayload(parts[1])
	if err != nil {
		return nil, err
	}
	s.logger.WithField("payload", payload)
	if payload.Exp < time.Now().Unix() {
		return nil, errors.New("token expired")
	}
	jwt, err := models.NewJwt(payload, "JWT")
	if err != nil {
		return nil, err
	}
	valid, err := jwt.CheckSignature(parts[2])
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid signature")
	}
	return payload, nil
}
