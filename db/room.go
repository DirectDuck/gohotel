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
	ctx context.Context, Room *types.Room,
) (*types.Room, error) {
	result, err := self.dbColl.InsertOne(ctx, Room)
	if err != nil {
		return nil, err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("Failed to cast %v to id", result.InsertedID)
	}
	return self.GetRoomByID(ctx, insertedID)
}
