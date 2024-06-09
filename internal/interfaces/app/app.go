package app

import (
	"app/internal/handler/login"
	"app/internal/handler/register"
	"app/internal/infras/database"
	"app/internal/interfaces/app/middlewares"
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
	return a
}

func (s *app) Start() error {
	// Middlwares
	s.engine.Use(gin.Recovery())
	s.engine.Use(middlewares.ErrorHandlerMiddleware())

	usecase := usecases.NewPostgresUsecase(s.dbc)

	// Routers
	routing(s.engine, s.cfg.ApiPrefix, usecase)

	return s.srv.ListenAndServe()
}

func routing(engine *gin.Engine, apiPrefix string, usecase usecases.Usecases) {
	router := engine.Group(apiPrefix)

	router.POST("/register", register.NewHandler(usecase.UserUC).Handle)
	router.POST("/login", login.NewHandler(usecase.UserUC).Handle)
}
