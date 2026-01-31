// Package logging provides structured logging for the backend.
// It wraps Go's slog package to provide:
// - Environment-based log format (JSON for prod, text for dev)
// - Configurable log levels
// - Request ID tracking
// - WebSocket client logging helpers
package logging

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"
)

// contextKey is a private type for context keys to avoid collisions
type contextKey string

const requestIDKey contextKey = "request_id"

var (
	// logger is the global logger instance
	logger *slog.Logger

	// defaultLogger is used before InitLogger is called
	defaultLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
)

// Config holds logging configuration
type Config struct {
	Level     slog.Level
	Format    string // "json" or "text"
	AddSource bool
}

// InitLogger initializes the global logger with configuration from environment.
// This should be called once at application startup.
func InitLogger() {
	config := loadConfig()

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: config.AddSource,
	}

	if config.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)
}

// loadConfig reads logging configuration from environment variables
func loadConfig() Config {
	config := Config{
		Level:  slog.LevelInfo,
		Format: "text",
	}

	// Parse LOG_LEVEL
	if levelStr := os.Getenv("LOG_LEVEL"); levelStr != "" {
		switch strings.ToLower(levelStr) {
		case "debug":
			config.Level = slog.LevelDebug
		case "info":
			config.Level = slog.LevelInfo
		case "warn":
			config.Level = slog.LevelWarn
		case "error":
			config.Level = slog.LevelError
		}
	}

	// Parse LOG_FORMAT
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		config.Format = strings.ToLower(format)
	}

	// Add source in debug mode
	if config.Level == slog.LevelDebug {
		config.AddSource = true
	}

	return config
}

// Logger returns the global logger instance.
// Returns a default logger if InitLogger hasn't been called.
func Logger() *slog.Logger {
	if logger != nil {
		return logger
	}
	return defaultLogger
}

// WithRequestID creates a new context with the given request ID.
// The request ID will be included in all logs using this context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// RequestID extracts the request ID from the context.
// Returns empty string if no request ID is set.
func RequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// WithContext returns a logger with fields from the context.
// Currently adds request_id if present in context.
func WithContext(ctx context.Context) *slog.Logger {
	l := Logger()
	if ctx == nil {
		return l
	}

	if requestID := RequestID(ctx); requestID != "" {
		return l.With("request_id", requestID)
	}
	return l
}

// WSLogger creates a logger for WebSocket client operations.
// Includes room_code and player_id fields.
func WSLogger(roomCode, playerID string) *slog.Logger {
	return Logger().With(
		"room_code", roomCode,
		"player_id", playerID,
	)
}

// WSLoggerWithContext creates a WebSocket logger from a context.
func WSLoggerWithContext(ctx context.Context, roomCode, playerID string) *slog.Logger {
	l := WithContext(ctx)
	return l.With(
		"room_code", roomCode,
		"player_id", playerID,
	)
}

// RoomLogger creates a logger for room operations.
func RoomLogger(roomCode string) *slog.Logger {
	return Logger().With("room_code", roomCode)
}

// StartupLog logs application startup information.
func StartupLog(port string, dbConnected bool) {
	Logger().Info("server starting",
		slog.String("port", port),
		slog.Bool("database_connected", dbConnected),
		slog.Time("startup_time", time.Now()),
	)
}

// ShutdownLog logs application shutdown.
func ShutdownLog(reason string) {
	Logger().Info("server shutting down",
		slog.String("reason", reason),
		slog.Time("shutdown_time", time.Now()),
	)
}

// Fatal logs a fatal error and exits the application.
func Fatal(msg string, err error) {
	Logger().Error(msg, slog.String("error", err.Error()))
	os.Exit(1)
}

// FatalWithContext logs a fatal error with context and exits.
func FatalWithContext(ctx context.Context, msg string, err error) {
	WithContext(ctx).Error(msg, slog.String("error", err.Error()))
	os.Exit(1)
}
