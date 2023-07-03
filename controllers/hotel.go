package controllers

import (
	"context"
	"fmt"
	"hotel/db"
	"hotel/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	minHotelNameLen     = 2
	minHotelLocationLen = 2
)

type HotelController struct {
	Store *db.Store
}

func (self *HotelController) HotelToWIthRooms(ctx context.Context, hotel *types.Hotel) (*types.HotelWithRooms, error) {
	roomQuery, err := bson.Marshal(RoomGetQueryParams{HotelID: hotel.ID})
	rooms, err := self.Store.Rooms.Get(ctx, roomQuery)
	if err != nil {
		return nil, err
	}

	return &types.HotelWithRooms{
		Hotel: hotel,
		Rooms: rooms,
	}, nil
}

func (self *HotelController) GetByID(
	ctx context.Context, id primitive.ObjectID,
) (*types.Hotel, error) {
	return self.Store.Hotels.GetByID(ctx, id)
}

func (self *HotelController) GetWithRoomsByID(
	ctx context.Context, id primitive.ObjectID,
) (*types.HotelWithRooms, error) {
	hotel, err := self.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return self.HotelToWIthRooms(ctx, hotel)
}

func (self *HotelController) Get(ctx context.Context) ([]*types.Hotel, error) {
	return self.Store.Hotels.Get(ctx)
}

func (self *HotelController) Validate(hotel *types.Hotel) map[string]string {
	errors := map[string]string{}
	if len(hotel.Name) < minHotelNameLen {
		errors["firstName"] = fmt.Sprintf(
			"Name length should be at least %d characters", minHotelNameLen,
		)
	}

	if len(hotel.Location) < minHotelLocationLen {
		errors["lastName"] = fmt.Sprintf(
			"Location length should be at least %d characters", minHotelLocationLen,
		)
	}

	return errors
}

func (self *HotelController) Evaluate(hotel *types.Hotel) error {
	return nil
}

func (self *HotelController) Create(
	ctx context.Context, hotel *types.Hotel,
) (*types.HotelWithRooms, error) {
	errs := self.Validate(hotel)
	if len(errs) != 0 {
		return nil, ValidationError{Fields: errs}
	}
	err := self.Evaluate(hotel)
	if err != nil {
		return nil, err
	}
	id, err := self.Store.Hotels.Create(ctx, hotel)
	if err != nil {
		return nil, err
	}
	created, err := self.GetByID(ctx, id)
	return self.HotelToWIthRooms(ctx, created)
}

func (self *HotelController) UpdateByID(
	ctx context.Context, id primitive.ObjectID, hotel *types.Hotel,
) (*types.HotelWithRooms, error) {
	errs := self.Validate(hotel)
	if len(errs) != 0 {
		return nil, ValidationError{Fields: errs}
	}
	err := self.Evaluate(hotel)
	if err != nil {
		return nil, err
	}

	err = self.Store.Hotels.UpdateByID(ctx, id, hotel)
	if err != nil {
		return nil, err
	}
	updated, err := self.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return self.HotelToWIthRooms(ctx, updated)
}

func (self *HotelController) DeleteByID(
	ctx context.Context, id primitive.ObjectID,
) error {
	return self.Store.Hotels.DeleteByID(ctx, id)
}
