package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// query is the handler for the /query endpoint.
func query(c echo.Context) error {
	query := c.QueryParam("q")
	return c.String(http.StatusOK, "query: "+query)
}

// ingest is the handler for the /data endpoint.
func ingest(c echo.Context) error {
	table := c.QueryParam("table")
	return c.String(http.StatusCreated, "table: "+table)
}
