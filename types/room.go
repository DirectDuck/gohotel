package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type RoomType int

const (
	SingleRoomType  RoomType = 5
	DoubleRoomType  RoomType = 10
	SeaSideRoomType RoomType = 15
	DeluxeRoomType  RoomType = 20
)

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      RoomType           `bson:"type" json:"type"`
	BasePrice float64            `bson:"basePrice" json:"basePrice"`
	Price     float64            `bson:"price" json:"price"`
	HotelID   primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}
