package handler

import (
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func (h *Handler) HandleAuth(c echo.Context) error {
	return c.JSON(200, "todo")
}

func (h *Handler) HandleAuthCallback(c echo.Context) error {
	code := c.QueryParam("code")

	if code == "" {
		return c.JSON(400, "code not found")
	}

	oAuthClient := NewDiscordOauth()

	token, err := oAuthClient.Exchange(c.Request().Context(), code)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name:    "token",
		Value:   token.AccessToken,
		Expires: token.Expiry,
		Path:    "/",
		Secure:  true,
	})
	c.SetCookie(&http.Cookie{
		Name:   "refresh_token",
		Value:  token.RefreshToken,
		Path:   "/",
		Secure: true,
	})
	// Send to dashboard
	return c.Redirect(302, "/dash")
}

func NewDiscordOauth() *oauth2.Config {
	cfg := &oauth2.Config{
		ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		Scopes:       []string{"identify", "email", "connections"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
		RedirectURL: os.Getenv("REDIRECT_URL"),
	}

	return cfg
}
