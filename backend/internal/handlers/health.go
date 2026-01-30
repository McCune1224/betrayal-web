package handlers

import (
	"net/http"

	"backend/internal/db"
	"github.com/labstack/echo/v5"
)

func HealthHandler(c *echo.Context) error {
	ctx := c.Request().Context()

	if err := db.HealthCheck(ctx); err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "unhealthy",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
	})
}
