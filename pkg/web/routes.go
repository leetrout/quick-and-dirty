package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type RouteMap map[string]string

// Cache the reversed map
var _reversedRouteMap = RouteMap{}

// Reverse allows a reverse lookup from URL to name.
func (rm *RouteMap) Reverse(url string) string {
	// Memoize the map
	if len(_reversedRouteMap) == 0 {
		for name, url := range *rm {
			_reversedRouteMap[url] = name
		}
	}

	return _reversedRouteMap[url]
}

// RouteLookup maps friendly names to URL paths.
var RouteLookup = RouteMap{
	"home":   "/",
	"query":  "/query",
	"ingest": "/data",
}

// ConfigureRoutes sets up our desired URL routes on the given
// Echo server instance.
func ConfigureRoutes(e *echo.Echo) {
	e.GET(RouteLookup["home"], func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET(RouteLookup["query"], query)

	e.POST(RouteLookup["ingest"], ingest)
}
