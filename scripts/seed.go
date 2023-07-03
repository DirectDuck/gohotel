package main

import (
	"context"
	"fmt"
	"hotel/db"
	"hotel/types"
	"log"
)

func main() {
	dbSrc := db.GetDatabase()
	hotelID, err := dbSrc.Store.Hotels.Create(
		context.TODO(),
		&types.Hotel{
			Name:     "Hotel 1",
			Location: "Berlin",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	rooms := []*types.Room{
		{
			Type:    types.SingleRoomType,
			HotelID: hotelID,
		},
		{
			Type:    types.DoubleRoomType,
			HotelID: hotelID,
		},
	}

	for _, room := range rooms {
		_, err := dbSrc.Store.Rooms.Create(
			context.TODO(), room,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Seeding done")
}
