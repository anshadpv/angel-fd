package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/angel-one/fd-core/commons/log"
	_ "github.com/lib/pq"
)

type dbConfig struct {
	DriverName            string
	URL                   string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdleTime time.Duration
}

func InitDBPool(ctx context.Context, config dbConfig) (*sql.DB, error) {
	// open the database
	var err error
	var pool *sql.DB
	pool, err = sql.Open(
		config.DriverName,
		config.URL,
	)
	if err != nil {
		return nil, err
	}

	// set the configurations
	pool.SetMaxOpenConns(config.MaxOpenConnections)
	pool.SetMaxIdleConns(config.MaxIdleConnections)
	pool.SetConnMaxIdleTime(config.ConnectionMaxIdleTime)
	pool.SetConnMaxLifetime(config.ConnectionMaxLifetime)

	ctx, stop := context.WithCancel(ctx)
	defer stop()
	Ping(ctx, pool)

	return pool, nil
}

func Ping(ctx context.Context, pool *sql.DB) {

	// The pq driver does not support ctx timeout. So the timeout is a no-op.
	// But it's useful for other drivers
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// The pq driver does not implement Pinger interface.
	if err := pool.PingContext(ctx); err != nil {
		log.Fatal(ctx).Err(err).Stack().Msg("unable to ping database")
	}
}

func Close(pool *sql.DB) error {
	err := pool.Close()
	if err != nil {
		return err
	}
	return nil
}

func GetConnectionsInUse(pool *sql.DB) int {
	return pool.Stats().InUse
}
