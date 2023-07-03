package controllers

import (
	"context"
	"fmt"
	"hotel/db"
	"hotel/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingController struct {
	Store *db.Store
}

func (self *BookingController) BookingToUnfolded(ctx context.Context, booking *types.Booking) (*types.BookingUnfolded, error) {
	room, err := self.Store.Rooms.GetByID(ctx, booking.RoomID)
	if err != nil {
		return nil, err
	}
	user, err := self.Store.Users.GetByID(ctx, booking.UserID)
	if err != nil {
		return nil, err
	}

	return &types.BookingUnfolded{
		Booking: booking,
		Room:    room,
		User:    user,
	}, nil
}

func (self *BookingController) GetByID(
	ctx context.Context, id primitive.ObjectID,
) (*types.Booking, error) {
	return self.Store.Bookings.GetByID(ctx, id)
}

func (self *BookingController) GetUnfoldedByID(
	ctx context.Context, id primitive.ObjectID,
) (*types.BookingUnfolded, error) {
	booking, err := self.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return self.BookingToUnfolded(ctx, booking)
}

type BookingGetQueryParams struct {
	UserID primitive.ObjectID `bson:"userID,omitempty" json:"userID"`
	RoomID primitive.ObjectID `bson:"roomID,omitempty" json:"roomID"`
}

func (self *BookingController) Get(
	ctx context.Context, query *BookingGetQueryParams,
) ([]*types.Booking, error) {
	if query == nil {
		query = &BookingGetQueryParams{}
	}
	queryMarshalled, err := bson.Marshal(query)
	if err != nil {
		return nil, err
	}
	return self.Store.Bookings.Get(ctx, queryMarshalled)
}

func (self *BookingController) Validate(booking *types.BookingUnfolded) map[string]string {
	errors := map[string]string{}
	if booking.Room == nil {
		errors["roomID"] = fmt.Sprintf("Room not found")
	}
	if booking.User == nil {
		errors["userID"] = fmt.Sprintf("User not found")
	}
	if booking.DateTo.Before(booking.DateFrom) {
		errors["dateTo"] = fmt.Sprintf("Date to can't be less than date from")
	}
	return errors
}

func (self *BookingController) Evaluate(booking *types.BookingUnfolded) error {
	booking.TotalCost = booking.Room.Price * float64(booking.DateTo.DaysSince(booking.DateFrom))
	return nil
}

func (self *BookingController) Create(
	ctx context.Context, booking *types.Booking,
) (*types.BookingUnfolded, error) {
	bookingUnfolded, err := self.BookingToUnfolded(ctx, booking)
	if err != nil {
		return nil, err
	}
	errs := self.Validate(bookingUnfolded)
	if len(errs) != 0 {
		return nil, ValidationError{Fields: errs}
	}
	err = self.Evaluate(bookingUnfolded)
	if err != nil {
		return nil, err
	}
	id, err := self.Store.Bookings.Create(ctx, bookingUnfolded.Booking)
	if err != nil {
		return nil, err
	}
	created, err := self.Store.Bookings.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return self.BookingToUnfolded(ctx, created)
}

func (self *BookingController) UpdateByID(
	ctx context.Context, id primitive.ObjectID, booking *types.Booking,
) (*types.BookingUnfolded, error) {
	bookingUnfolded, err := self.BookingToUnfolded(ctx, booking)
	if err != nil {
		return nil, err
	}

	errs := self.Validate(bookingUnfolded)
	if len(errs) != 0 {
		return nil, ValidationError{Fields: errs}
	}
	err = self.Evaluate(bookingUnfolded)
	if err != nil {
		return nil, err
	}

	err = self.Store.Bookings.UpdateByID(ctx, id, bookingUnfolded.Booking)
	if err != nil {
		return nil, err
	}
	return self.GetUnfoldedByID(ctx, id)
}

func (self *BookingController) DeleteByID(
	ctx context.Context, id primitive.ObjectID,
) error {
	return self.Store.Bookings.DeleteByID(ctx, id)
}
