package handler

import (
	"io"
	"text/template"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/mccune1224/betrayal/pkg/data"
)

type Handler struct {
	models    data.Models
	templates *EchoTemplates
}

// PageData is the data that is passed to the template to render the page
type PageData struct {
  Title string
  ContentTemplate string
}


func NewHandler(DB *sqlx.DB) *Handler {
	return &Handler{
		models:    data.NewModels(DB),
		templates: NewTemplates(),
	}
}

func (h *Handler) GetTemplates() *EchoTemplates {
  return h.templates
}

// EchoTemplates is a custom renderer for echo wrapped with helper functions
type EchoTemplates struct {
	templates *template.Template
}

func (t *EchoTemplates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}



func NewTemplates() *EchoTemplates {
	return &EchoTemplates{
		templates: template.Must(template.ParseGlob("views//*.html")),
	}
}
