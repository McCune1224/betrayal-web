package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"backend/internal/db"
	"backend/internal/game"
	"backend/internal/handlers"
	"backend/internal/logging"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Initialize logging first
	logging.InitLogger()
	logger := logging.Logger()

	// Log startup
	logger.Info("application starting",
		"version", "1.0.0",
		"time", time.Now(),
	)

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logging.Fatal("DATABASE_URL not set", fmt.Errorf("environment variable DATABASE_URL is required"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize database
	if err := db.InitDB(ctx, dbURL); err != nil {
		logging.Fatal("Failed to connect to database", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		db.CloseDB(ctx)
	}()

	e := echo.New()

	// Middleware - ORDER MATTERS
	// 1. Recovery first to catch panics
	e.Use(middleware.Recover())
	// 2. Request logging
	e.Use(logging.HTTPMiddleware())
	// 3. CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type"},
	}))

	// Initialize game hub
	hub := game.NewHub()
	go hub.Run()
	logger.Info("game hub initialized")

	// Handlers
	roomHandler := &handlers.RoomHandler{Hub: hub}
	e.POST("/api/rooms", roomHandler.CreateRoom)
	e.POST("/api/rooms/:code/join", roomHandler.JoinRoom)
	e.GET("/ws", roomHandler.UpgradeWebSocket)
	// DB health check
	e.GET("/api/health/db", handlers.HealthHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logging.StartupLog(port, true)
	if err := e.Start(fmt.Sprintf(":%s", port)); err != nil {
		logging.Fatal("Failed to start server", err)
	}
}
