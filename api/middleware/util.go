package middleware

import (
	"strings"

	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
)

func isExcludedPath(ctx *gin.Context) bool {
	return strings.Contains(ctx.FullPath(), constants.ActuatorRoute) || strings.Contains(ctx.FullPath(), constants.SwaggerRoute)
}
