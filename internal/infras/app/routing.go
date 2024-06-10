package app

import (
	"app/internal/handler/v1/auth"
	"app/internal/interfaces/usecases"

	"github.com/gin-gonic/gin"
)

func routing(engine *gin.Engine, apiPrefix string, usecase usecases.Usecases) {
	router := engine.Group(apiPrefix)

	{
		authHandler := auth.NewHandler(usecase.UserUC)

    // /auth/*
		authRouter := router.Group("auth")
		authRouter.POST("/register", authHandler.HandleRegister)
		authRouter.POST("/login", authHandler.HandleLogin)
	}
}
