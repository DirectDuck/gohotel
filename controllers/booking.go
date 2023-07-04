package controllers

import (
	"context"
	"fmt"
	"hotel/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingController struct {
	Store *Store
}

func (self *BookingController) BookingToUnfolded(ctx context.Context, booking *types.Booking) (*types.BookingUnfolded, error) {
	if booking == nil {
		return nil, nil
	}
	result, err := self.Store.DB.Rooms.GetOneByID(ctx, booking.RoomID, &types.Room{})
	if err != nil {
		return nil, err
	}
	room := CastPtrInterface[types.Room](result)
	result, err = self.Store.DB.Users.GetOneByID(ctx, booking.UserID, &types.User{})
	if err != nil {
		return nil, err
	}
	user := CastPtrInterface[types.User](result)

	return &types.BookingUnfolded{
		Booking: booking,
		Room:    room,
		User:    user,
	}, nil
}

func (self *BookingController) GetByID(
	ctx context.Context, id primitive.ObjectID,
) (*types.Booking, error) {
	result, err := self.Store.DB.Bookings.GetOneByID(ctx, id, &types.Booking{})
	if err != nil {
		return nil, err
	}
	return CastPtrInterface[types.Booking](result), nil
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
	UserID primitive.ObjectID `bson:"userID,omitempty" json:"-"`
	RoomID primitive.ObjectID `bson:"roomID,omitempty" json:"roomID"`
}

func (self *BookingController) Get(
	ctx context.Context, query *BookingGetQueryParams,
) ([]*types.Booking, error) {
	if query == nil {
		query = &BookingGetQueryParams{}
	}
	user, err := GetUserFromContext(self.Store.DB, ctx)
	if !user.IsAdmin {
		query.UserID = user.ID
	}
	if err != nil {
		return nil, err
	}
	result, err := self.Store.DB.Bookings.Get(ctx, query, []*types.Booking{})
	if err != nil {
		return nil, err
	}
	return CastInterface[[]*types.Booking](result), nil
}

func (self *BookingController) GetOccupiedForRoom(
	ctx context.Context, roomID primitive.ObjectID,
) ([]*types.BookingDates, error) {
	query := &BookingGetQueryParams{RoomID: roomID}
	result, err := self.Store.DB.Bookings.Get(ctx, query, []*types.BookingDates{})
	if err != nil {
		return nil, err
	}
	return CastInterface[[]*types.BookingDates](result), nil
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
	userID, err := GetUserIDFromContext(self.Store.DB, ctx)
	if err != nil {
		return nil, err
	}
	booking.UserID = userID
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
	id, err := self.Store.DB.Bookings.Create(ctx, bookingUnfolded.Booking)
	if err != nil {
		return nil, err
	}
	created, err := self.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return self.BookingToUnfolded(ctx, created)
}

func (self *BookingController) UpdateByID(
	ctx context.Context, id primitive.ObjectID, booking *types.Booking,
) (*types.BookingUnfolded, error) {
	userID, err := GetUserIDFromContext(self.Store.DB, ctx)
	if err != nil {
		return nil, err
	}
	booking.UserID = userID
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

	err = self.Store.DB.Bookings.UpdateByID(ctx, id, bookingUnfolded.Booking)
	if err != nil {
		return nil, err
	}
	return self.GetUnfoldedByID(ctx, id)
}

func (self *BookingController) DeleteByID(
	ctx context.Context, id primitive.ObjectID,
) error {
	return self.Store.DB.Bookings.DeleteByID(ctx, id)
}
