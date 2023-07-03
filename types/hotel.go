package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hotel struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	Location string             `bson:"location" json:"location"`
}

type HotelWithRooms struct {
	*Hotel
	Rooms []*Room `bson:"-" json:"rooms"`
}

type BaseHotelParams struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type CreateHotelParams struct {
	BaseHotelParams
}

type UpdateHotelParams struct {
	BaseHotelParams
}

func NewHotelFromCreateParams(params CreateHotelParams) (*Hotel, error) {
	return &Hotel{
		Name:     params.Name,
		Location: params.Location,
	}, nil
}

func NewHotelFromUpdateParams(params UpdateHotelParams) (*Hotel, error) {
	return &Hotel{
		Name:     params.Name,
		Location: params.Location,
	}, nil
}
