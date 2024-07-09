package middleware

import (
	"time"

	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		reqID := requestID(ctx)

		ctx.Set(constants.LogRequestTimeKey, start)
		ctx.Set(constants.LogPathKey, ctx.Request.URL.Path)
		ctx.Set(constants.HeaderUserID, ctx.Request.Header.Get(constants.HeaderUserID))
		ctx.Set(constants.HeaderRequestID, reqID)
		ctx.Set(constants.LogClientIPKey, ctx.ClientIP())
		ctx.Set(constants.LogAppVersionKey, ctx.Request.Header.Get(constants.HeaderAppVersion))
		ctx.Set(constants.LogMACAddressKey, ctx.Request.Header.Get(constants.HeaderMACAddress))
		ctx.Set(constants.LogOSKey, ctx.Request.Header.Get(constants.HeaderOS))
		ctx.Set(constants.LogSourceKey, ctx.Request.Header.Get(constants.HeaderSource))
		ctx.Set(constants.HeaderAcceptLanguage, ctx.Request.Header.Get(constants.HeaderAcceptLanguage))

		//Language Headers en-us

		// Process request
		ctx.Next()
	}
}

func requestID(c *gin.Context) string {
	reqID := c.Request.Header.Get(constants.HeaderRequestID)
	if reqID == "" {
		reqID = uuid.NewString()
	}
	return reqID
}
