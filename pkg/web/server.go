package web

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewEcho returns an Echo server configured with
// routes, middleware, etc
func NewEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	ConfigureRoutes(e)

	t := &Template{
		templates: template.Must(template.New("t").Funcs(template.FuncMap{
			"ReverseURL": e.Reverse,
		}).ParseGlob("pkg/web/templates/*.html")),
	}
	e.Renderer = t

	return e
}

// ConfigureRoutes sets up our desired URL routes on the given
// Echo server instance.
func ConfigureRoutes(e *echo.Echo) {
	e.GET("/", home).Name = "home"

	e.GET("/query", query).Name = "query"

	e.POST("/data", ingest).Name = "ingest"
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
