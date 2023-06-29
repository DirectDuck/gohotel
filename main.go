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

	userHandler := api.NewUserHandler(dbSrc.Store)

	apiv1 := app.Group("/api/v1")
	apiv1.Post("/user", userHandler.HandleCreateUser)
	apiv1.Get("/user", userHandler.HandleListUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Put("/user/:id", userHandler.HandleUpdateUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	hotelHandler := api.NewHotelHandler(dbSrc.Store)

	apiv1.Post("/hotel", hotelHandler.HandleCreateHotel)
	apiv1.Get("/hotel", hotelHandler.HandleListHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Put("/hotel/:id", hotelHandler.HandleUpdateHotel)
	apiv1.Delete("/hotel/:id", hotelHandler.HandleDeleteHotel)

	roomHandler := api.NewRoomHandler(dbSrc.Store)

	apiv1.Post("/room", roomHandler.HandleCreateRoom)
	apiv1.Get("/room", roomHandler.HandleListRooms)
	apiv1.Get("/room/:id", roomHandler.HandleGetRoom)
	apiv1.Put("/room/:id", roomHandler.HandleUpdateRoom)
	apiv1.Delete("/room/:id", roomHandler.HandleDeleteRoom)

	app.Listen(":" + *listenPort)
}
