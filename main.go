package main

import (
	"errors"
	"fmt"

	"github.com/angel-one/fd-core/api/routes"
	"github.com/angel-one/fd-core/business/jobs"
	"github.com/angel-one/fd-core/commons/cache"
	"github.com/angel-one/fd-core/commons/config"
	"github.com/angel-one/fd-core/commons/context"
	"github.com/angel-one/fd-core/commons/database"
	"github.com/angel-one/fd-core/commons/flags"
	"github.com/angel-one/fd-core/commons/httpclient"
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/fd-core/constants"
	"github.com/angel-one/fd-core/factory"

	_ "github.com/angel-one/fd-core/docs"
)

// @title FD Backend/API Service
// @version 1.0
// @description Fixed Deposit Backend & API service
// @termsOfService https://swagger.io/terms/

// @contact.name Madhan Ganesh
// @contact.team NewBusinessUnit (NBU)
// @contact.email

var ctx = context.Background("boot")

// custom hook - pre init
func init() {
	log.Info(ctx).Msg("bootstrapping fd-core service...")
}

// config initialization
func init() {
	var err error
	configNames := []string{constants.ApplicationConfig, constants.LoggerConfig, constants.DatabaseConfig, constants.HTTPClientConfig, constants.RedisConfig}
	if flags.Mode() == "test" {
		err = config.InitTestMode(fmt.Sprintf("%s/%s", flags.BaseConfigPath(), flags.Env()), configNames...)
	} else {
		err = config.InitReleaseMode(configNames...)
	}
	if err != nil {
		log.Fatal(ctx).Err(err).Stack().Msg("failed to initialize configs")
	}
	log.Info(ctx).Msg("configuration initialized")
}

func init() {
	authKey := config.Default().Secrets[constants.JWTSymmetricKey]
	if authKey == "" {
		log.Fatal(ctx).Err(errors.New(constants.ErrAuthKeyEnvNotSet)).Msg(constants.ErrAuthKeyEnvNotSet)
	}

	log.Info(ctx).Msgf("starting with environment: '%s' ; mode: %s", flags.Env(), flags.Mode())
}

// logger initialization
func init() {
	loglevel, err := config.Default().GetString(constants.LoggerConfig, constants.LogLevelKey)
	if err != nil {
		log.Fatal(ctx).Err(err).Stack().Msg("failed to initialize logger")
	}
	log.InitLogger(log.Level(loglevel))
	log.Info(ctx).Msg("logger initialized")
}

// performs db initialization
func init() {
	err := database.Init(ctx, config.Default(), constants.DatabaseConfig)
	if err != nil {
		log.Fatal(ctx).Err(err).Stack().Msg("failed to initialize database")
	}
	log.Info(ctx).Msg("database initialized")
}

// performs redis initialization
func init() {
	err := cache.Init(ctx, config.Default(), constants.RedisConfig)
	if err != nil {
		log.Fatal(ctx).Err(err).Stack().Msg("failed to initialize redis cluster cache client")
	}
	log.Info(ctx).Msg("redis cluster initialized")
}

// perform http clients
func init() {
	clientConfigKeys := []string{constants.UpSwingGenerateToken, constants.UpSwingPCIRegistration, constants.UpSwingNetWorth, constants.UpswingDataIngestion, constants.ProfileServerConfig, constants.UpswingPendingJourney}
	if err := httpclient.Init(ctx, config.Default(), constants.HTTPClientConfig, clientConfigKeys); err != nil {
		log.Fatal(ctx).Err(err).Stack().Msg("failed to initialize http client")
	}

	log.Info(ctx).Msg("httpclient initalized")
}

// init factory
func init() {
	factory.Init(ctx)
	log.Info(ctx).Msg("fd-core factory initialized...")
}

// start cron jobs
func init() {
	disabled := config.Default().GetBoolD(constants.ApplicationConfig, "jobsDisabled", false)
	if disabled {
		log.Info(ctx).Msg("jobs are marked not to run , its state is disabled, skipping startin it")
	} else {
		jobs.StartJobs()
		log.Info(ctx).Msg("started all cron jobs")
	}
}

// custom hook - post init
func init() {
	log.Info(ctx).Msg("fd-core service initialized...")
}

// main function
func main() {
	router := routes.DefaultRouter(ctx)
	log.Info(ctx).Msgf("starting server and listening to port: %d", flags.Port())
	err := router.Run(fmt.Sprintf(":%d", flags.Port()))
	if err != nil {
		log.Fatal(ctx).Err(err).Msg("server startup failed")
	}
	log.Info(ctx).Msg("api server started...")
}
