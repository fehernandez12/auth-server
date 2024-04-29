package server

import (
	"auth-server/hasher"
	"auth-server/logger"
	"auth-server/models"
	"auth-server/repository"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServerConfig struct {
	Timeout int
	Addr    string
	Secret  []byte
}

type Server struct {
	config                *ServerConfig
	DB                    *gorm.DB
	clientRepository      repository.Repository[models.Client]
	applicationRepository repository.Repository[models.Application]
	userRepository        repository.Repository[models.User]
	roleRepository        repository.Repository[models.Role]
	permissionRepository  repository.Repository[models.Permission]
	logger                *logger.Logger
	hasher                hasher.Hasher
}

func StartServer() error {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	stopper := make(chan struct{})
	go func() {
		<-done
		close(stopper)
	}()
	server, err := newServer()
	if err != nil {
		return err
	}
	return server.start(stopper)
}

func newServer() (*Server, error) {
	s := &Server{
		logger: logger.NewLogger(),
	}
	banner, err := os.ReadFile("resources/banner.txt")
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
	err = db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.Application{},
		&models.Client{},
	)
	if err != nil {
		s.logger.Fatal(err)
	}
	s.DB = db
	s.clientRepository = repository.NewClientRepository(db)
	s.applicationRepository = repository.NewApplicationRepository(db)
	s.userRepository = repository.NewUserRepository(db)
	s.roleRepository = repository.NewRoleRepository(db)
	s.permissionRepository = repository.NewPermissionRepository(db)
	s.hasher = hasher.NewPBKDF2Hasher(200000, s.config.Secret)
	s.logger.WithField("Status", "Application is running")
	return s, nil
}

func (s *Server) start(stop <-chan struct{}) error {
	corsObj := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{
			X_REQUESTED_WITH, CONTENT_TYPE, AUTHORIZATION,
		}),
	)
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
		Secret:  []byte(os.Getenv("AUTH_SERVER_SECRET")),
	}, nil
}
