package api

import (
	"github.com/angel-one/fd-core/commons/flags"
	"github.com/angel-one/fd-core/constants"
	goActuator "github.com/angel-one/go-actuator"
	"github.com/gin-gonic/gin"
)

var (
	actuatorHandler = goActuator.GetActuatorHandler(&goActuator.Config{
		Env:     flags.Env(),
		Name:    constants.ApplicationName,
		Port:    flags.Port(),
		Version: constants.V1,
		Endpoints: []int{
			goActuator.Info,
			goActuator.Ping,
			goActuator.Metrics,
		},
	})
)

func Actuator(ctx *gin.Context) {
	actuatorHandler(ctx.Writer, ctx.Request)
}
