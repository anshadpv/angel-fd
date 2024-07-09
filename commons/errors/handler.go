package errors

import (
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/goerr"
	"github.com/gin-gonic/gin"
)

type ErrResponse struct {
	Message string `json:"message,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

func Wrap(fderr *Error, err error) {
	fderr.Err = err
}

func Throw(ctx *gin.Context, err error) {
	response := ErrResponse{
		Reason: err.Error(),
	}
	log.Error(ctx).Msg(goerr.Stack(err))
	ctx.AbortWithStatusJSON(goerr.Code(err), response)
}

func Message(ctx *gin.Context, err error) {
	response := ErrResponse{
		Message: err.Error(),
	}
	log.Debug(ctx).Msg(err.Error())
	ctx.AbortWithStatusJSON(goerr.Code(err), response)
}
