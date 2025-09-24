package templater

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer реализует echo.Renderer
type TemplateRenderer struct {
    Templates *template.Template
}

// Render метод для рендеринга шаблона
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.Templates.ExecuteTemplate(w, name, data)
}