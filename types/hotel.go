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

type BaseHotelParams struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

func (self *BaseHotelParams) Validate() map[string]string {
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

type CreateHotelParams struct {
	BaseHotelParams
}

func (self *CreateHotelParams) Validate() map[string]string {
	errors := self.BaseHotelParams.Validate()
	return errors
}

type UpdateHotelParams struct {
	BaseHotelParams
}

func (self *UpdateHotelParams) Validate() map[string]string {
	errors := self.BaseHotelParams.Validate()
	return errors
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
