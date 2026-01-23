package handlers

import (
	"betrayal-web/internal"
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// HealthHandler returns HTTP 200 if DB is alive
func HealthHandler(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
	defer cancel()
	if internal.Conn == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "db not connected"})
	}
	// Use the Pingable interface for testability (see internal/db.go)
	if err := internal.Conn.Ping(ctx); err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "db error", "error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
