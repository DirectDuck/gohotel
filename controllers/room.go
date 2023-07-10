package controllers

import (
	"context"
	"fmt"
	roomprices_rpc "hotel/services/roomprices/rpc"
	"hotel/types"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomController struct {
	Store *Store
}

func (self *RoomController) RoomToUnfolded(ctx context.Context, room *types.Room) (*types.RoomUnfolded, error) {
	if room == nil {
		return nil, nil
	}
	result, err := self.Store.DB.Hotels.GetOneByID(ctx, room.HotelID, &types.Hotel{})
	if err != nil {
		return nil, err
	}
	hotel := CastPtrInterface[types.Hotel](result)

	bookingDates, err := self.Store.CT.Bookings.GetOccupiedForRoom(ctx, room.ID)
	if err != nil {
		return nil, err
	}

	return &types.RoomUnfolded{
		Room:        room,
		Hotel:       hotel,
		BookedDates: bookingDates,
	}, nil
}

func (self *RoomController) GetByID(
	ctx context.Context, id primitive.ObjectID,
) (*types.Room, error) {
	result, err := self.Store.DB.Rooms.GetOneByID(ctx, id, &types.Room{})
	if err != nil {
		return nil, err
	}
	return CastPtrInterface[types.Room](result), nil
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
	result, err := self.Store.DB.Rooms.Get(ctx, queryMarshalled, []*types.Room{})
	if err != nil {
		return nil, err
	}
	return CastInterface[[]*types.Room](result), nil
}

func (self *RoomController) Validate(room *types.RoomUnfolded) map[string]string {
	errors := map[string]string{}
	if !room.Type.IsValid() {
		errors["type"] = fmt.Sprintf("Invalid room type")
	}

	return errors
}

func (self *RoomController) Evaluate(ctx context.Context, room *types.RoomUnfolded) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	resp, err := self.Store.RoomPrices.GetRoomPrice(
		ctxWithTimeout, &roomprices_rpc.RoomPriceRequest{
			Type: int64(room.Type),
		},
	)
	if err != nil {
		return err
	}
	room.Price = resp.Price
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
	err = self.Evaluate(ctx, roomUnfolded)
	if err != nil {
		return nil, err
	}
	id, err := self.Store.DB.Rooms.Create(ctx, roomUnfolded.Room)
	if err != nil {
		return nil, err
	}
	created, err := self.GetByID(ctx, id)
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
	err = self.Evaluate(ctx, roomUnfolded)
	if err != nil {
		return nil, err
	}

	err = self.Store.DB.Rooms.UpdateByID(ctx, id, roomUnfolded.Room)
	if err != nil {
		return nil, err
	}
	return self.GetUnfoldedByID(ctx, id)
}

func (self *RoomController) DeleteByID(
	ctx context.Context, id primitive.ObjectID,
) error {
	return self.Store.DB.Rooms.DeleteByID(ctx, id)
}
