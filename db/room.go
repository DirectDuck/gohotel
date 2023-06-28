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

const dbRoomsCollectionName = "rooms"

type RoomStore interface {
	CreateRoom(context.Context, *types.Room) (*types.Room, error)
	GetRoomByID(context.Context, primitive.ObjectID) (*types.Room, error)
	GetRooms(context.Context) ([]*types.Room, error)
	GetRoomsForHotel(context.Context, primitive.ObjectID) ([]*types.Room, error)
	UpdateRoomByID(context.Context, primitive.ObjectID, *types.Room) (*types.Room, error)
	DeleteRoomByID(context.Context, primitive.ObjectID) error
}

type MongoRoomStore struct {
	dbSrc  *mongo.Database
	dbColl *mongo.Collection
}

func NewMongoRoomStore(dbSrc *mongo.Database) *MongoRoomStore {
	return &MongoRoomStore{
		dbSrc:  dbSrc,
		dbColl: dbSrc.Collection(dbRoomsCollectionName),
	}
}

func (self *MongoRoomStore) GetRoomByID(
	ctx context.Context, id primitive.ObjectID,
) (*types.Room, error) {
	room := &types.Room{}

	err := self.dbColl.FindOne(
		ctx, bson.M{"_id": id},
	).Decode(room)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return room, nil
}

func (self *MongoRoomStore) CreateRoom(
	ctx context.Context, room *types.Room,
) (*types.Room, error) {
	result, err := self.dbColl.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("Failed to cast %v to id", result.InsertedID)
	}
	return self.GetRoomByID(ctx, insertedID)
}

func (self *MongoRoomStore) GetRooms(ctx context.Context) ([]*types.Room, error) {
	cursor, err := self.dbColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room

	err = cursor.All(ctx, &rooms)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (self *MongoRoomStore) GetRoomsForHotel(
	ctx context.Context, hotelID primitive.ObjectID,
) ([]*types.Room, error) {
	cursor, err := self.dbColl.Find(ctx, bson.M{
		"hotelID": hotelID,
	})
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room

	err = cursor.All(ctx, &rooms)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (self *MongoRoomStore) UpdateRoomByID(
	ctx context.Context, id primitive.ObjectID, data *types.Room,
) (*types.Room, error) {

	_, err := self.dbColl.UpdateByID(
		ctx, id, bson.M{"$set": data},
	)
	if err != nil {
		return nil, err
	}

	room, err := self.GetRoomByID(ctx, id)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return room, nil
}

func (self *MongoRoomStore) DeleteRoomByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := self.dbColl.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
