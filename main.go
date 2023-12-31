package main

import (
	"log"
	"os"

	"github.com/McCune1224/betrayal-web/endpoints"
	"github.com/McCune1224/betrayal-web/handler"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	app := echo.New()
	// Connect to DB
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("error opening database,", err)
	}

	app.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: "${status} | ${latency_human} | ${method} | ${uri} | ${error} \n",
		},
	))

  app.Use(middleware.Recover())

  //trailling slash
  app.Pre(middleware.RemoveTrailingSlash())

	handler := handler.NewHandler(db)
	app.Renderer = handler.GetTemplates()
	app.Static("/static", "static")
	endpoints.AttachRoutes(app, handler)

	log.Fatal(app.Start(":" + os.Getenv("PORT")))
}
