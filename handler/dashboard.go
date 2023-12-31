package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleDashboard(e echo.Context) error {
	cookie, err := e.Cookie("token")
	if err != nil {
		return e.JSON(401, echo.Map{"error": "Unauthorized"})
	}

	discordClient, err := discordgo.New("Bearer " + cookie.Value)
	if err != nil {
		return e.JSON(500, err.Error())
	}
	// get user info
	user, err := discordClient.User("@me")
	if err != nil {
		return e.JSON(500, err.Error())
	}

	userAvatarURL := "https://cdn.discordapp.com/avatars/" + user.ID + "/" + user.Avatar + ".png"


  data := echo.Map{
    "Username": user.Username,
    "Avatar": userAvatarURL,
  }

	return e.Render(200, "dashboard.html", data)
}
