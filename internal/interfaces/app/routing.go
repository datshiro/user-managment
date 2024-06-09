package app

import (
	"app/internal/handler/v1/auth"
	"app/internal/interfaces/usecases"

	"github.com/gin-gonic/gin"
)

func routing(engine *gin.Engine, apiPrefix string, usecase usecases.Usecases) {
	router := engine.Group(apiPrefix)

	{
		authRouter := router.Group("auth")
		authHandler := auth.NewHandler(usecase.UserUC)

		authRouter.POST("/register", authHandler.HandleRegister)
		authRouter.POST("/login", authHandler.HandleLogin)

	}
}
