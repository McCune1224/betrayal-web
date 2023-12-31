package handler

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo/v4"
	"github.com/mccune1224/betrayal/pkg/data"
)

func (h *Handler) HandleDashboard(e echo.Context) error {
	cookie, err := e.Cookie("token")
	if err != nil {
		return e.Redirect(302, "/")
	}
	disc, err := discordgo.New("Bearer " + cookie.Value)
	if err != nil {
		return e.JSON(500, err.Error())
	}

	// get user info
	user, err := disc.User("@me")
	if err != nil {
		return e.JSON(500, err.Error())
	}

	userAvatarURL := "https://cdn.discordapp.com/avatars/" + user.ID + "/" + user.Avatar + ".png"

	data := echo.Map{
		"Username": user.Username,
		"Avatar":   userAvatarURL,
	}

	return e.Render(200, "dashboard.html", data)
}

type UserWrapper struct {
  discordgo.User
  IconURL string
}

type InventoryData struct {
	Inventory data.Inventory
	User      UserWrapper
}

func (h *Handler) HandleInventories(e echo.Context) error {
	_, err := e.Cookie("token")
	if err != nil {
		return e.Redirect(302, "/")
	}
	bot, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		return e.JSON(500, "Failed to create discord session"+err.Error())
	}

	inventories, err := h.models.Inventories.GetAll()
	if err != nil {
		return e.JSON(500, err.Error())
	}

	invDatas := make([]InventoryData, len(inventories))

	for i, inv := range inventories {
		user, err := bot.User(inv.DiscordID)
		if err != nil {
			return e.JSON(500, "Failed to fetch user "+err.Error())
		}

    u := UserWrapper{*user, user.AvatarURL("")}

		invDatas[i] = InventoryData{
			User:      u,
			Inventory: inv,
		}
	}
	// return e.JSON(200, invDatas)
	return e.Render(200, "inventories.html", invDatas)
}
