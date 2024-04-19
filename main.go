package main

import (
	"github.com/labstack/echo/v4"

	"qad/pkg/web"
)

func main() {
	e := echo.New()
	web.ConfigureRoutes(e)
	e.Logger.Fatal(e.Start(":1337"))
}
