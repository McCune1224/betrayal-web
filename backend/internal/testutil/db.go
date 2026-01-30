package testutil

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var testPool *pgxpool.Pool

func init() {
	godotenv.Load("../../.env")
}

func SetupTestDB(t testing.TB) *pgxpool.Pool {
	t.Helper()

	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping database test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(testDBURL)
	require.NoError(t, err, "failed to parse test database config")

	config.MaxConns = 5
	config.MinConns = 1

	pool, err := pgxpool.NewWithConfig(ctx, config)
	require.NoError(t, err, "failed to create test database pool")

	err = pool.Ping(ctx)
	require.NoError(t, err, "failed to ping test database")

	testPool = pool
	return pool
}

func TeardownTestDB(t testing.TB, pool *pgxpool.Pool) {
	t.Helper()
	if pool != nil {
		pool.Close()
	}
}

func TruncateTables(t testing.TB, pool *pgxpool.Pool, tables ...string) {
	t.Helper()

	if len(tables) == 0 {
		tables = []string{"actions", "players", "rooms", "roles"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, table := range tables {
		_, err := pool.Exec(ctx, "TRUNCATE TABLE "+table+" CASCADE")
		require.NoError(t, err, "failed to truncate table %s", table)
	}
}

func GetTestPool() *pgxpool.Pool {
	return testPool
}
