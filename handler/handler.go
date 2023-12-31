package handler

import (
	"io"
	"log"
	"path/filepath"
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
  dirs := []string{
    "views/*.html",
    "views/dashboards/*.html",
  }
  files := []string{}
  for _, dir := range dirs {
    ff, err := filepath.Glob(dir)
    if err != nil {
      log.Fatal(err)
    }
    files = append(files, ff...)
  }

	return &EchoTemplates{
		templates: template.Must(template.ParseFiles(files...)),
	}
}
