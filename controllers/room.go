package controllers

import (
	"context"
	"fmt"
	"hotel/db"
	"hotel/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomController struct {
	Store *db.Store
}

func (self *RoomController) RoomToUnfolded(ctx context.Context, room *types.Room) (*types.RoomUnfolded, error) {
	hotel, err := self.Store.Hotels.GetByID(ctx, room.HotelID)
	if err != nil {
		return nil, err
	}

	bookingQuery, err := bson.Marshal(BookingGetQueryParams{RoomID: room.ID})
	if err != nil {
		return nil, err
	}

	bookings, err := self.Store.Bookings.Get(ctx, bookingQuery)
	if err != nil {
		return nil, err
	}

	return &types.RoomUnfolded{
		Room:     room,
		Hotel:    hotel,
		Bookings: bookings,
	}, nil
}

func (self *RoomController) GetByID(
	ctx context.Context, id primitive.ObjectID,
) (*types.Room, error) {
	return self.Store.Rooms.GetByID(ctx, id)
}

func (self *RoomController) GetUnfoldedByID(
	ctx context.Context, id primitive.ObjectID,
) (*types.RoomUnfolded, error) {
	room, err := self.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return self.RoomToUnfolded(ctx, room)
}

type RoomGetQueryParams struct {
	HotelID primitive.ObjectID `bson:"hotelID,omitempty" json:"hotelID"`
}

func (self *RoomController) Get(
	ctx context.Context, query *RoomGetQueryParams,
) ([]*types.Room, error) {
	if query == nil {
		query = &RoomGetQueryParams{}
	}
	queryMarshalled, err := bson.Marshal(query)
	if err != nil {
		return nil, err
	}
	return self.Store.Rooms.Get(ctx, queryMarshalled)
}

func (self *RoomController) Validate(room *types.RoomUnfolded) map[string]string {
	errors := map[string]string{}
	if room.Price < 0 {
		errors["price"] = fmt.Sprintf("Price can't be less than 0")
	}

	if !room.Type.IsValid() {
		errors["type"] = fmt.Sprintf("Invalid room type")
	}

	return errors
}

func (self *RoomController) Evaluate(room *types.RoomUnfolded) error {
	return nil
}

func (self *RoomController) Create(
	ctx context.Context, room *types.Room,
) (*types.RoomUnfolded, error) {
	roomUnfolded, err := self.RoomToUnfolded(ctx, room)
	if err != nil {
		return nil, err
	}
	errs := self.Validate(roomUnfolded)
	if len(errs) != 0 {
		return nil, ValidationError{Fields: errs}
	}
	err = self.Evaluate(roomUnfolded)
	if err != nil {
		return nil, err
	}
	id, err := self.Store.Rooms.Create(ctx, roomUnfolded.Room)
	if err != nil {
		return nil, err
	}
	created, err := self.Store.Rooms.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return self.RoomToUnfolded(ctx, created)
}

func (self *RoomController) UpdateByID(
	ctx context.Context, id primitive.ObjectID, room *types.Room,
) (*types.RoomUnfolded, error) {
	roomUnfolded, err := self.RoomToUnfolded(ctx, room)
	if err != nil {
		return nil, err
	}

	errs := self.Validate(roomUnfolded)
	if len(errs) != 0 {
		return nil, ValidationError{Fields: errs}
	}
	err = self.Evaluate(roomUnfolded)
	if err != nil {
		return nil, err
	}

	err = self.Store.Rooms.UpdateByID(ctx, id, roomUnfolded.Room)
	if err != nil {
		return nil, err
	}
	return self.GetUnfoldedByID(ctx, id)
}

func (self *RoomController) DeleteByID(
	ctx context.Context, id primitive.ObjectID,
) error {
	return self.Store.Rooms.DeleteByID(ctx, id)
}
