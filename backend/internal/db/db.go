package db

import (
	"context"
	"fmt"
	"time"

	"backend/internal/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func InitDB(ctx context.Context, connString string) error {
	logger := logging.Logger()

	logger.Info("initializing_database_connection")

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		logger.Error("failed_to_parse_database_config",
			"error", err,
		)
		return fmt.Errorf("unable to parse database config: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = 5 * time.Minute

	logger.Debug("database_config",
		"max_conns", config.MaxConns,
		"min_conns", config.MinConns,
		"max_conn_lifetime", config.MaxConnLifetime,
	)

	pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Error("failed_to_create_connection_pool",
			"error", err,
		)
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		logger.Error("failed_to_ping_database",
			"error", err,
		)
		return fmt.Errorf("unable to ping database: %w", err)
	}

	logger.Info("database_connection_established",
		"max_conns", config.MaxConns,
		"min_conns", config.MinConns,
	)
	return nil
}

func CloseDB(ctx context.Context) error {
	logger := logging.Logger()

	if pool != nil {
		pool.Close()
		logger.Info("database_connection_closed")
	}
	return nil
}

func GetPool() *pgxpool.Pool {
	return pool
}

func HealthCheck(ctx context.Context) error {
	logger := logging.Logger()

	if pool == nil {
		logger.Warn("health_check_failed",
			"reason", "pool_not_initialized",
		)
		return fmt.Errorf("database pool not initialized")
	}

	err := pool.Ping(ctx)
	if err != nil {
		logger.Error("health_check_failed",
			"reason", "ping_failed",
			"error", err,
		)
		return err
	}

	logger.Debug("health_check_passed")
	return nil
}
