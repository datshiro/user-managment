package app

import (
	"app/internal/handler/register"
	"app/internal/interfaces/app/middlewares"
	"context"
	"fmt"
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

	// Routers
	routing(s.engine, s.cfg.ApiPrefix)

	return s.srv.ListenAndServe()
}

func routing(engine *gin.Engine, apiPrefix string) {
	registerHandler := register.NewHandler()
  router := engine.Group(apiPrefix)

	router.POST("/register", registerHandler.Handle)
}
