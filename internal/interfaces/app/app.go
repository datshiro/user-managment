package app

import (
	"app/internal/handler/login"
	"app/internal/handler/register"
	"app/internal/infras/database"
	"app/internal/interfaces/app/middlewares"
	"app/internal/interfaces/usecases"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server interface {
	Start() error
	Stop(ctx context.Context) error
}

type server struct {
	engine *gin.Engine
	srv    *http.Server
	cfg    Config
}

func (s *server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func NewApp(cfg Config) Server {
	engine := gin.Default()
	srv := &http.Server{
		Handler: engine,
		Addr:    fmt.Sprintf(":%d", cfg.Port),
	}
	return &server{engine: engine, cfg: cfg, srv: srv}
}

func (s *server) Start() error {
	// Middlwares
	s.engine.Use(gin.Recovery())
	s.engine.Use(middlewares.ErrorHandlerMiddleware())

	dbc, err := initDatabase(s.cfg.DbConfig)
	if err != nil {
		log.Fatalf("Failed to make connection to database: %s |  %+v", err, s.cfg.DbConfig)
	}

	usecase := usecases.NewPostgresUsecase(dbc)

	// Routers
	routing(s.engine, s.cfg.ApiPrefix, usecase)

	return s.srv.ListenAndServe()
}

func routing(engine *gin.Engine, apiPrefix string, usecase usecases.Usecases) {
	router := engine.Group(apiPrefix)

	router.POST("/register", register.NewHandler(usecase.UserUC).Handle)
	router.POST("/login", login.NewHandler(usecase.UserUC).Handle)
}

func initDatabase(cfg database.DBConfig) (*gorm.DB, error) {
	// database
	dbConfig := database.DBConfig{
		Host:     cfg.Host,
		User:     cfg.User,
		Password: cfg.Password,
		Port:     cfg.Port,
		DB:       cfg.DB,
		SSLMode:  cfg.SSLMode,
		TimeZone: cfg.TimeZone,
	}
	return dbConfig.NewPostgresConnection()
}
