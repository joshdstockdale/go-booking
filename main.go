package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/joshdstockdale/go-booking/api"
)

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "Listen address of the API Server")
	flag.Parse()

	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/user", api.HandleGetUsers)
	apiv1.Get("/user/:id", api.HandleGetUser)
	app.Listen(*listenAddr)
}
