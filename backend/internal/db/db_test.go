package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitDB(t *testing.T) {
	t.Run("fails with invalid connection string", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := InitDB(ctx, "invalid://connection/string")

		assert.Error(t, err, "should return error for invalid connection string")
	})
}

func TestCloseDB(t *testing.T) {
	t.Run("does not panic when pool is nil", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		assert.NotPanics(t, func() {
			CloseDB(ctx)
		}, "should not panic when closing nil pool")
	})
}

func TestGetPool(t *testing.T) {
	t.Run("returns nil when not initialized", func(t *testing.T) {
		pool := GetPool()
		assert.Nil(t, pool, "pool should be nil when not initialized")
	})
}

func TestHealthCheck(t *testing.T) {
	t.Run("fails when pool is nil", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := HealthCheck(ctx)

		assert.Error(t, err, "should return error when pool is nil")
		assert.Contains(t, err.Error(), "not initialized")
	})
}
