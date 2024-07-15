package routes

import (
	"context"
	"errors"
	"net/http"

	"github.com/angel-one/fd-core/api/middleware"
	"github.com/angel-one/fd-core/commons/config"

	"github.com/angel-one/fd-core/api"
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func createRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.New()
	router.Use(middlewares...)
	router.Use(gin.Recovery())
	return router
}

func DefaultRouter(ctx context.Context) *gin.Engine {
	authKey := config.Default().Secrets[constants.JWTSymmetricKey]
	allowedOrigins := config.Default().GetStringD(constants.ApplicationConfig, constants.AllowedOrginsForCorsConfigKey, constants.AllowedOriginsForCorsDefault)
	if allowedOrigins == "" {
		log.Fatal(ctx).Err(errors.New(constants.WhitelistedHostIsNotSet)).Msg(constants.WhitelistedHostIsNotSet)
	}
	whitelistedHeaderHosts := config.Default().GetStringD(constants.ApplicationConfig, constants.WhitelistedHostsHeader, constants.WhitelistedHostsHeaderDefault)
	if allowedOrigins == "" {
		log.Fatal(ctx).Err(errors.New(constants.WhitelistedHostIsNotSet)).Msg(constants.WhitelistedHostIsNotSet)
	}

	router := createRouter(
		middleware.CORS(middleware.CORSMiddlewareOptions{AllowedMethods: constants.AllowedMethodsForCors, AllowedOrigins: allowedOrigins}),
		middleware.SecurityHeader(),
		middleware.HeaderChecks(middleware.HeadersMiddlewareOptions{WhitelistedHosts: whitelistedHeaderHosts, ExcludedURI: []string{constants.ActuatorInfo, constants.SwaggerRoute}}),
		middleware.Auth([]byte(authKey)),
		middleware.Logger(),
	)

	router.GET(constants.SwaggerRoute, ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET(constants.ActuatorRoute, api.Actuator)

	// dedicate group for vendors webhooks
	InitWebhookRoutes(router.Group(constants.Webhook))

	// dedicate group for all version1 APIs
	v1Group := router.Group(constants.V1)

	// init all routes
	InitHomepageRoute(v1Group)
	InitTokenRoute(v1Group)
	InitPortfolioRoute(v1Group)
	initPlansV1Group(v1Group)
	initFAQ(v1Group)
	InitComparePageRoute(v1Group)
	InitJobsRoute(v1Group)
	InitFsiDetails(v1Group)
	InitHomeInfoRoute(v1Group)
	InitUpswingRoute(v1Group)

	// init invalid routes
	initNoRoute(router)

	return router
}

func initNoRoute(router *gin.Engine) {
	router.NoRoute(handleUnknown)
}

func handleUnknown(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Unsupported URI"})
}
