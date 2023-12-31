package endpoints

import (
	"github.com/McCune1224/betrayal-web/handler"
	"github.com/labstack/echo/v4"
)

func AttachRoutes(app *echo.Echo, handler *handler.Handler) {
	app.GET("/", func(c echo.Context) error {
		return c.Render(200, "index.html", nil)
	})
	// auth := app.Group("/auth")
	//
	// auth.GET("/", handler.HandleAuth)
}
