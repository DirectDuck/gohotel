package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomType int

const (
	SingleRoomType  RoomType = 5
	DoubleRoomType  RoomType = 10
	SeaSideRoomType RoomType = 15
	DeluxeRoomType  RoomType = 20
)

func (self RoomType) IsValid() bool {
	switch self {
	case
		SingleRoomType, DoubleRoomType,
		SeaSideRoomType, DeluxeRoomType:
		return true
	}
	return false
}

type Room struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type    RoomType           `bson:"type" json:"type"`
	Price   float64            `bson:"price" json:"price"`
	HotelID primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}

type RoomUnfolded struct {
	*Room
	Hotel       *Hotel          `bson:"-" json:"hotel"`
	BookedDates []*BookingDates `bson:"-" json:"bookedDates"`
}

type BaseRoomParams struct {
	Type    RoomType           `json:"type"`
	Price   float64            `json:"price"`
	HotelID primitive.ObjectID `json:"hotelID"`
}

type CreateRoomParams struct {
	BaseRoomParams
}

type UpdateRoomParams struct {
	BaseRoomParams
}

func NewRoomFromCreateParams(params CreateRoomParams) (*Room, error) {
	return &Room{
		Type:    params.Type,
		Price:   params.Price,
		HotelID: params.HotelID,
	}, nil
}

func NewRoomFromUpdateParams(params UpdateRoomParams) (*Room, error) {
	return &Room{
		Type:    params.Type,
		Price:   params.Price,
		HotelID: params.HotelID,
	}, nil
}
