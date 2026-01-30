package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"backend/internal/db"
	"backend/internal/game"
	"backend/internal/handlers"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.InitDB(ctx, dbURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		db.CloseDB(ctx)
	}()

	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type"},
	}))

	// Initialize game hub
	hub := game.NewHub()
	go hub.Run()

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

	log.Printf("Server running on :%s\n", port)
	if err := e.Start(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
