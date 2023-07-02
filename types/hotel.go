package types

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	minHotelNameLen     = 2
	minHotelLocationLen = 2
)

type Hotel struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	Location string             `bson:"location" json:"location"`
}

func (self *Hotel) Validate(dbBefore *Hotel) map[string]string {
	errors := map[string]string{}
	if len(self.Name) < minHotelNameLen {
		errors["firstName"] = fmt.Sprintf(
			"Name length should be at least %d characters", minHotelNameLen,
		)
	}

	if len(self.Location) < minHotelLocationLen {
		errors["lastName"] = fmt.Sprintf(
			"Location length should be at least %d characters", minHotelLocationLen,
		)
	}

	return errors
}

func (self *Hotel) Evaluate(dbBefore *Hotel) error {
	return nil
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
