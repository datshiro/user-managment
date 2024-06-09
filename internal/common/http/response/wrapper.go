package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Wrapper interface {
	JSON(data interface{})
	Error(code int, err error)
}

func New(c *gin.Context, opts ...OptFunc) Wrapper {
	o := defaultOpts()

	for _, opt := range opts {
		opt(&o)
	}
	return &wrapper{ctx: c, Opts: o}
}

type wrapper struct {
	ctx *gin.Context
	Opts
	// Extend logger later
}

func (w *wrapper) JSON(data interface{}) {
	var code  = http.StatusOK
	if w.Code != 0 {
		code = w.Code
	}
	if gin.IsDebugging() {
		w.ctx.IndentedJSON(code, ResponseObject{Success: true, Data: data})
	} else {
		w.ctx.JSON(code, ResponseObject{Success: true, Data: data})
	}
}

func (w *wrapper) Error(code int, err error) {
	message := err.Error()
	w.ctx.JSON(code, ResponseObject{Success: false, Message: message})
}
