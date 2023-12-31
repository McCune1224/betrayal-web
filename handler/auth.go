package handler

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleAuth(c echo.Context) error {
	data, err := h.models.Abilities.GetByFuzzy("fire")
	if err != nil {
		return err
	}

	payload := echo.Map{
		"Ability": data,
	}

	return c.Render(200, "auth.html", payload)
}
