package main

import (
	"os"
	"os/signal"
	"syscall"

	"qad/pkg/data"
	"qad/pkg/web"
)

func main() {

	e := web.NewEcho()

	// Create a db connection
	db := data.Connect(data.ConnectWithDSN("foobar.duckdb"))
	defer db.Connection.Close()

	// Capture ctrl+c
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		os.Exit(1)
	}()

	// Run the webserver
	e.Logger.Fatal(e.Start(":1337"))
}
