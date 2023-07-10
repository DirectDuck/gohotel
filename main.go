package main

import (
	"context"
	"hotel/api"
	"hotel/controllers"
	"hotel/db"
	"log"
	"time"

	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getRoompricesConn() *grpc.ClientConn {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	roompricesConn, err := grpc.DialContext(
		ctx, os.Getenv("ROOMPRICES_LISTEN_URL"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("Failed to connect to Roomprices service: %s\n", err.Error())
	}
	return roompricesConn
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	app := fiber.New(
		fiber.Config{
			ErrorHandler: api.HandleAPIError,
		},
	)

	roompricesConn := getRoompricesConn()
	defer roompricesConn.Close()

	CTStore := controllers.NewStore(db.GetDatabase(), roompricesConn)

	userHandler := api.NewUserHandler(
		&controllers.UserController{Store: CTStore},
	)

	apiv1 := app.Group("/api/v1")
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

	app.Listen(os.Getenv("APP_LISTEN_URL"))
}
