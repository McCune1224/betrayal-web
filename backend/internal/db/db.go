package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func InitDB(ctx context.Context, connString string) error {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return fmt.Errorf("unable to parse database config: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = 5 * time.Minute

	pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Database connection established successfully")
	return nil
}

func CloseDB(ctx context.Context) error {
	if pool != nil {
		pool.Close()
		log.Println("Database connection closed")
	}
	return nil
}

func GetPool() *pgxpool.Pool {
	return pool
}

func HealthCheck(ctx context.Context) error {
	if pool == nil {
		return fmt.Errorf("database pool not initialized")
	}
	return pool.Ping(ctx)
}
