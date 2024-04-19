package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// NewEcho returns an Echo server configured with
// routes, middleware, etc
func NewEcho() *echo.Echo {
	e := echo.New()
	ConfigureRoutes(e)
	return e
}

// ConfigureRoutes sets up our desired URL routes on the given
// Echo server instance.
func ConfigureRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}).Name = "home"

	e.GET("/query", query).Name = "query"

	e.POST("/data", ingest).Name = "ingest"
}
