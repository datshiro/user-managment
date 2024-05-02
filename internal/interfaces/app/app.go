package app

import (
	"app/internal/handler/register"
	"app/internal/infras/database"
	"app/internal/interfaces/app/middlewares"
	"app/internal/interfaces/usecases"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

  // database 
  dbConfig := database.DBConfig{
			Host:               "localhost",
			User:               "postgres",
			Password:           "postgres",
			Port:               "5432",
			DB:                 "cake_db",
			SSLMode:            "disable",
			TimeZone:           "Asia/Shanghai",
  }
  dbc , err := dbConfig.Connect()
  if err != nil {
    log.Fatalf("Failed to make connection to database: %+v", dbConfig)
  }

  usecase := usecases.NewPostgresUsecase(dbc)

	// Routers
	routing(s.engine, s.cfg.ApiPrefix, usecase)

	return s.srv.ListenAndServe()
}

func routing(engine *gin.Engine, apiPrefix string, usecase usecases.Usecases) {
	registerHandler := register.NewHandler(usecase.UserUC)
  router := engine.Group(apiPrefix)

	router.POST("/register", registerHandler.Handle)
}
