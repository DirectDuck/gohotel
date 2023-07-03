package types

import (
	"cloud.google.com/go/civil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RoomID    primitive.ObjectID `bson:"roomID" json:"roomID"`
	UserID    primitive.ObjectID `bson:"userID" json:"userID"`
	DateFrom  civil.Date         `bson:"dateFrom" json:"dateFrom"`
	DateTo    civil.Date         `bson:"dateTo" json:"dateTo"`
	TotalCost float64            `bson:"totalCost" json:"totalCost"`
}

type BookingUnfolded struct {
	*Booking
	Room *Room `bson:"-" json:"room"`
	User *User `bson:"-" json:"user"`
}

type BaseBookingParams struct {
	RoomID   primitive.ObjectID `json:"roomID"`
	UserID   primitive.ObjectID `json:"userID"`
	DateFrom civil.Date         `json:"dateFrom"`
	DateTo   civil.Date         `json:"dateTo"`
}

type CreateBookingParams struct {
	BaseBookingParams
}

type UpdateBookingParams struct {
	BaseBookingParams
}

func NewBookingFromCreateParams(params CreateBookingParams) (*Booking, error) {
	return &Booking{
		RoomID:   params.RoomID,
		UserID:   params.UserID,
		DateFrom: params.DateFrom,
		DateTo:   params.DateTo,
	}, nil
}

func NewBookingFromUpdateParams(params UpdateBookingParams) (*Booking, error) {
	return &Booking{
		RoomID:   params.RoomID,
		UserID:   params.UserID,
		DateFrom: params.DateFrom,
		DateTo:   params.DateTo,
	}, nil
}
