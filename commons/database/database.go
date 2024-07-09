package database

import (
	"context"
	"database/sql"
	"time"

	fderr "github.com/angel-one/fd-core/commons/errors"
	"github.com/angel-one/fd-core/constants"

	"github.com/angel-one/fd-core/commons/config"
	"github.com/angel-one/fd-core/commons/log"
)

var ConfigNotFound = fderr.New().Code("DB-01").Msg("configs not found from config file").Build()

var postgresDB *sql.DB

var (
	DBDefaultMaxOpenConnections             = 25
	DBDefaultMaxIdleConnections             = 25
	DBDefaultConnectionMaxLifetimeInSeconds = 90
	DBDefaultConnectionMaxIdleTimeInSeconds = 30
	DBDefaultConnecTimeout                  = 10
)

func GetDBPool(postgres bool) *sql.DB {
	return postgresDB
}

func Init(ctx context.Context, cfg *config.Client, key string) error {
	var err error
	if postgresDB, err = initDBPools(ctx, cfg, key, constants.PostgresDBKey); err != nil {
		return err
	}
	return nil
}

func initDBPools(ctx context.Context, cfg *config.Client, key string, poolType string) (*sql.DB, error) {
	var err error
	var dbPool *sql.DB
	var config map[string]interface{}

	if config, err = cfg.GetMap(key, poolType); err != nil {
		log.Fatal(ctx).Err(err).Msgf("error during load of key: %s; pooType: %s configs", key, poolType)
		return nil, err
	}
	if len(config) == 0 {
		log.Fatal(ctx).Err(err).Msgf("missing config for key: %s; pooType: %s configs", key, poolType)
		return nil, err
	}
	dbPool, err = initConnection(ctx, cfg, config)
	if err != nil {
		log.Fatal(ctx).Err(err).Msgf("failed to initlaize %s db pool", poolType)
		return nil, err
	}
	log.Info(ctx).Msgf("successfully initialized db-pool for type %s", poolType)
	return dbPool, nil
}

func initConnection(ctx context.Context, cfg *config.Client, configMap map[string]interface{}) (*sql.DB, error) {
	var err error
	var pool *sql.DB
	driverName, _ := cfg.GetStringFromMap(configMap, constants.DBDriver, "")
	url, err := cfg.GetStringWithSecretsFromMap(configMap, constants.DBURL, "")
	if err != nil {
		log.Fatal(ctx).Err(err).Msg("error getting database url")
	}

	maxOpenConnections, _ := cfg.GetIntFromMap(configMap, constants.DBMaxOpenConnections, DBDefaultMaxOpenConnections)
	maxIdleConnections, _ := cfg.GetIntFromMap(configMap, constants.DBMaxIdleConnections, DBDefaultMaxIdleConnections)
	connectionMaxLifetime, _ := cfg.GetIntFromMap(configMap, constants.DBConnectionMaxLifetimeInSeconds, DBDefaultConnectionMaxLifetimeInSeconds)
	connectionMaxIdleTime, _ := cfg.GetIntFromMap(configMap, constants.DBConnectionMaxIdleTimeInSeconds, DBDefaultConnectionMaxIdleTimeInSeconds)

	pool, err = InitDBPool(ctx, dbConfig{
		DriverName:            driverName,
		URL:                   url,
		MaxOpenConnections:    int(maxOpenConnections),
		MaxIdleConnections:    int(maxIdleConnections),
		ConnectionMaxLifetime: time.Second * time.Duration(connectionMaxLifetime),
		ConnectionMaxIdleTime: time.Second * time.Duration(connectionMaxIdleTime),
	})

	return pool, err
}
