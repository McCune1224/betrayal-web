package logging

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"time"
)

// responseRecorder wraps echo.Context to capture status code
type responseRecorder struct {
	echo.Context
	status int
}

func (r *responseRecorder) JSON(code int, i interface{}) error {
	r.status = code
	return r.Context.JSON(code, i)
}

func (r *responseRecorder) String(code int, s string) error {
	r.status = code
	return r.Context.String(code, s)
}

func (r *responseRecorder) HTML(code int, html string) error {
	r.status = code
	return r.Context.HTML(code, html)
}

func (r *responseRecorder) NoContent(code int) error {
	r.status = code
	return r.Context.NoContent(code)
}

// HTTPMiddleware returns an Echo middleware that logs all HTTP requests.
// It adds request ID tracking and logs method, path, status, duration, and errors.
func HTTPMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()

			// Generate or extract request ID
			requestID := (*c).Request().Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// Add request ID to context and response header
			ctx := WithRequestID((*c).Request().Context(), requestID)
			(*c).SetRequest((*c).Request().WithContext(ctx))
			(*c).Response().Header().Set("X-Request-ID", requestID)

			// Get logger with context
			logger := WithContext(ctx)

			// Log request start at debug level
			logger.Debug("http_request_started",
				"method", (*c).Request().Method,
				"path", (*c).Request().URL.Path,
				"remote_addr", (*c).Request().RemoteAddr,
			)

			// Wrap context to capture status code
			recorder := &responseRecorder{
				Context: *c,
				status:  200, // Default status
			}

			// Process request
			err := next(&recorder.Context)

			// Calculate duration
			duration := time.Since(start)
			status := recorder.status

			// Build log fields
			logArgs := []interface{}{
				"method", (*c).Request().Method,
				"path", (*c).Request().URL.Path,
				"status", status,
				"duration_ms", duration.Milliseconds(),
				"remote_addr", (*c).Request().RemoteAddr,
			}

			// Add error if present
			if err != nil {
				logArgs = append(logArgs, "error", err.Error())
			}

			// Log at appropriate level based on status
			switch {
			case status >= 500:
				logger.Error("http_request", logArgs...)
			case status >= 400:
				logger.Warn("http_request", logArgs...)
			default:
				logger.Info("http_request", logArgs...)
			}

			return err
		}
	}
}
