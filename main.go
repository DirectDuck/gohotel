package main

import (
	"flag"
	"hotel/api"
	"hotel/db"

	"github.com/gofiber/fiber/v2"
)

const (
	dburi  = "mongodb://admin:admin@localhost:27017"
	dbName = "hotel-reservation"
)

func main() {
	listenPort := flag.String("port", "8000", "Port to run the API server")
	flag.Parse()

	app := fiber.New(
		fiber.Config{
			ErrorHandler: api.HandleAPIError,
		},
	)

	dbSrc := db.GetDatabase()

	userHandler := api.NewUserHandler(
		db.NewMongoUserStore(dbSrc),
	)

	apiv1 := app.Group("/api/v1")
	apiv1.Post("/user", userHandler.HandleCreateUser)
	apiv1.Get("/user", userHandler.HandleListUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Put("/user/:id", userHandler.HandleUpdateUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	app.Listen(":" + *listenPort)
}
