package endpoints

import (
	"github.com/McCune1224/betrayal-web/handler"
	"github.com/labstack/echo/v4"
)

func AttachRoutes(app *echo.Echo, h *handler.Handler) {
	app.GET("/", func(c echo.Context) error {
		token, _ := c.Cookie("token")

		if token != nil {
			return c.Redirect(302, "/dash")
		}

		oAuthClient := handler.NewDiscordOauth()

		data := echo.Map{
			"DiscordURL": oAuthClient.AuthCodeURL(""),
		}

		return c.Render(200, "index.html", data)
	})
	auth := app.Group("/auth")
	auth.GET("/", h.HandleAuth)
	auth.GET("/redirect", h.HandleAuthCallback)

	dashboard := app.Group("/dash")
	dashboard.GET("", h.HandleDashboard)
	dashboard.GET("/inventories", h.HandleInventories)
}
