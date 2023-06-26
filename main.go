package main

import (
	"context"
	"flag"
	"hotel/api"
	"hotel/db"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dburi  = "mongodb://admin:admin@localhost:27017"
	dbName = "hotel-reservation"
)

func main() {
	dbClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	listenPort := flag.String("port", "8000", "Port to run the API server")
	flag.Parse()

	app := fiber.New(
		fiber.Config{
			ErrorHandler: api.HandleAPIError,
		},
	)

	userHandler := api.NewUserHandler(
		db.NewMongoUserStore(dbClient.Database(dbName)),
	)

	apiv1 := app.Group("/api/v1")
	apiv1.Post("/user", userHandler.HandleCreateUser)
	apiv1.Get("/user", userHandler.HandleListUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Put("/user/:id", userHandler.HandleUpdateUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	app.Listen(":" + *listenPort)
}
