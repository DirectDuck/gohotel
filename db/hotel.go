package db

import (
	"context"
	"errors"
	"fmt"
	"hotel/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const dbHotelsCollectionName = "hotels"

type HotelStore interface {
	CreateHotel(context.Context, *types.Hotel) (*types.HotelWithRooms, error)
	GetHotelByID(context.Context, primitive.ObjectID) (*types.HotelWithRooms, error)
	GetHotels(context.Context) ([]*types.Hotel, error)
	UpdateHotelByID(context.Context, primitive.ObjectID, *types.Hotel) (*types.HotelWithRooms, error)
	DeleteHotelByID(context.Context, primitive.ObjectID) error
}

type MongoHotelStore struct {
	dbSrc     *mongo.Database
	dbColl    *mongo.Collection
	roomStore RoomStore
}

func NewMongoHotelStore(dbSrc *mongo.Database) *MongoHotelStore {
	return &MongoHotelStore{
		dbSrc:     dbSrc,
		dbColl:    dbSrc.Collection(dbHotelsCollectionName),
		roomStore: NewMongoRoomStore(dbSrc),
	}
}

func (self *MongoHotelStore) GetHotelByID(
	ctx context.Context, id primitive.ObjectID,
) (*types.HotelWithRooms, error) {
	hotel := &types.Hotel{}

	err := self.dbColl.FindOne(
		ctx, bson.M{"_id": id},
	).Decode(hotel)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	rooms, err := self.roomStore.GetRoomsForHotel(ctx, id)
	if err != nil {
		return nil, err
	}

	return &types.HotelWithRooms{
		Hotel: hotel,
		Rooms: rooms,
	}, nil
}

func (self *MongoHotelStore) CreateHotel(
	ctx context.Context, hotel *types.Hotel,
) (*types.HotelWithRooms, error) {
	result, err := self.dbColl.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("Failed to cast %v to id", result.InsertedID)
	}
	return self.GetHotelByID(ctx, insertedID)
}

func (self *MongoHotelStore) GetHotels(ctx context.Context) ([]*types.Hotel, error) {
	cursor, err := self.dbColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel

	err = cursor.All(ctx, &hotels)
	if err != nil {
		return nil, err
	}

	return hotels, nil
}

func (self *MongoHotelStore) UpdateHotelByID(
	ctx context.Context, id primitive.ObjectID, data *types.Hotel,
) (*types.HotelWithRooms, error) {

	_, err := self.dbColl.UpdateByID(
		ctx, id, bson.M{"$set": data},
	)
	if err != nil {
		return nil, err
	}

	hotel, err := self.GetHotelByID(ctx, id)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return hotel, nil
}

func (self *MongoHotelStore) DeleteHotelByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := self.dbColl.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
