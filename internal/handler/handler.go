package handler

import "github.com/gin-gonic/gin"

type Handler interface {
	Handle(*gin.Context)
}

type Request interface {
	Bind(*gin.Context) error
	Validate() error
}
