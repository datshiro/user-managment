package app

import (
	"app/internal/infras/database"
	"app/internal/infras/app/middlewares"
	"app/internal/interfaces/usecases"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App interface {
	Start() error
	Stop(ctx context.Context) error
}

type app struct {
	engine *gin.Engine
	srv    *http.Server
	cfg    Opts
	dbc    *gorm.DB
}

func (s *app) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func NewApp(opts ...OptFunc) App {
	engine := gin.Default()
	o := defaultOpts()

	for _, optFunc := range opts {
		optFunc(&o)
	}

	a := &app{
		engine: engine,
		srv: &http.Server{
			Handler: engine,
			Addr:    fmt.Sprintf(":%d", o.Port),
		},
	}

	if o.IsConnectDatabase {
		a.dbc = database.NewPostgresConnection()
	}

	// Routers
	usecase := usecases.NewPostgresUsecase(a.dbc)
	routing(engine, o.ApiPrefix, usecase)
	return a
}

func (s *app) Start() error {
	// Middlwares
	s.engine.Use(gin.Recovery())
	s.engine.Use(middlewares.ErrorHandlerMiddleware())

	return s.srv.ListenAndServe()
}

