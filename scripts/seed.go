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
	hotelStore := db.NewMongoHotelStore(dbSrc)
	hotel, err := hotelStore.CreateHotel(
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

	roomStore := db.NewMongoRoomStore(dbSrc)

	for _, room := range rooms {
		_, err := roomStore.CreateRoom(
			context.TODO(), room,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Seeding done")
}
