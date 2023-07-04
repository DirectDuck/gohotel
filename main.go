package main

import (
	"flag"
	"hotel/api"
	"hotel/controllers"
	"hotel/db"
	"log"

	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

const (
	dburi  = "mongodb://admin:admin@localhost:27017"
	dbName = "hotel-reservation"
)

func main() {
	listenPort := flag.String("port", "8000", "Port to run the API server")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	app := fiber.New(
		fiber.Config{
			ErrorHandler: api.HandleAPIError,
		},
	)

	apiv1 := app.Group("/api/v1")

	CTStore := controllers.NewStore(db.GetDatabase())

	userHandler := api.NewUserHandler(
		&controllers.UserController{Store: CTStore},
	)

	apiv1.Post("/login", userHandler.HandleLogin)

	secret := os.Getenv("JWT_SECRET")
	app.Use(jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte(secret)},
		TokenLookup: "header:Authorization",
	}))

	apiv1.Post("/user", userHandler.HandleCreateUser)
	apiv1.Get("/user", userHandler.HandleListUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Put("/user/:id", userHandler.HandleUpdateUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	hotelHandler := api.NewHotelHandler(
		&controllers.HotelController{Store: CTStore},
	)

	apiv1.Post("/hotel", hotelHandler.HandleCreateHotel)
	apiv1.Get("/hotel", hotelHandler.HandleListHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Put("/hotel/:id", hotelHandler.HandleUpdateHotel)
	apiv1.Delete("/hotel/:id", hotelHandler.HandleDeleteHotel)

	roomHandler := api.NewRoomHandler(
		&controllers.RoomController{Store: CTStore},
	)

	apiv1.Post("/room", roomHandler.HandleCreateRoom)
	apiv1.Get("/room", roomHandler.HandleListRooms)
	apiv1.Get("/room/:id", roomHandler.HandleGetRoom)
	apiv1.Put("/room/:id", roomHandler.HandleUpdateRoom)
	apiv1.Delete("/room/:id", roomHandler.HandleDeleteRoom)

	bookingHandler := api.NewBookingHandler(
		&controllers.BookingController{Store: CTStore},
	)

	apiv1.Post("/booking", bookingHandler.HandleCreateBooking)
	apiv1.Get("/booking", bookingHandler.HandleListBookings)
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Put("/booking/:id", bookingHandler.HandleUpdateBooking)
	apiv1.Delete("/booking/:id", bookingHandler.HandleDeleteBooking)

	app.Listen(":" + *listenPort)
}
