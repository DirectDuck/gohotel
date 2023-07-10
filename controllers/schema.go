package controllers

import (
	"hotel/db"
	roomprices_rpc "hotel/services/roomprices/rpc"

	"google.golang.org/grpc"
)

type Controllers struct {
	Users    *UserController
	Hotels   *HotelController
	Rooms    *RoomController
	Bookings *BookingController
}

type Store struct {
	DB         *db.DB
	CT         *Controllers
	RoomPrices roomprices_rpc.RoomPricesServiceClient
}

func NewStore(DB *db.DB, roompricesConn *grpc.ClientConn) *Store {
	store := &Store{
		DB:         DB,
		CT:         &Controllers{},
		RoomPrices: roomprices_rpc.NewRoomPricesServiceClient(roompricesConn),
	}
	store.CT.Users = &UserController{store}
	store.CT.Hotels = &HotelController{store}
	store.CT.Rooms = &RoomController{store}
	store.CT.Bookings = &BookingController{store}
	return store
}
