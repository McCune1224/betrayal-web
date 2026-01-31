package handlers

import (
	"net/http"

	"backend/internal/db"
	"backend/internal/logging"
	"github.com/labstack/echo/v5"
)

func HealthHandler(c *echo.Context) error {
	ctx := (*c).Request().Context()
	logger := logging.WithContext(ctx)

	logger.Debug("health_check_requested")

	if err := db.HealthCheck(ctx); err != nil {
		logger.Warn("health_check_failed",
			"error", err,
		)
		return (*c).JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "unhealthy",
			"error":  err.Error(),
		})
	}

	logger.Debug("health_check_passed")
	return (*c).JSON(http.StatusOK, map[string]string{
		"status": "healthy",
	})
}
