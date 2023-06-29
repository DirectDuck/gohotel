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
	hotel, err := dbSrc.Store.Hotels.Create(
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
			Type:      types.SingleRoomType,
			BasePrice: 150,
			HotelID:   hotel.ID,
		},
		{
			Type:      types.DoubleRoomType,
			BasePrice: 200,
			HotelID:   hotel.ID,
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
