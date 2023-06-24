package main

import (
	"flag"
	"hotel/api"

	"github.com/gofiber/fiber/v2"
)

func main() {
	listenPort := flag.String("port", "8000", "Port to run the API server")
	flag.Parse()

	app := fiber.New()

	apiv1 := app.Group("/api/v1")
	apiv1.Get("/user", api.HandleListUsers)
	apiv1.Get("/user/:id", api.HandleGetUser)

	app.Listen(":" + *listenPort)
}
