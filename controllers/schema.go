package controllers

import "hotel/db"

type Controllers struct {
	Users    *UserController
	Hotels   *HotelController
	Rooms    *RoomController
	Bookings *BookingController
}

type Store struct {
	DB *db.DB
	CT *Controllers
}

func NewStore(DB *db.DB) *Store {
	store := &Store{
		DB: DB,
		CT: &Controllers{},
	}
	store.CT.Users = &UserController{store}
	store.CT.Hotels = &HotelController{store}
	store.CT.Rooms = &RoomController{store}
	store.CT.Bookings = &BookingController{store}
	return store
}
