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
	Create(context.Context, *types.Room) (primitive.ObjectID, error)
	GetByID(context.Context, primitive.ObjectID) (*types.Room, error)
	Get(context.Context, []byte) ([]*types.Room, error)
	UpdateByID(context.Context, primitive.ObjectID, *types.Room) error
	DeleteByID(context.Context, primitive.ObjectID) error
}

type MongoRoomStore struct {
	db     *MongoDB
	dbColl *mongo.Collection
}

func NewMongoRoomStore(dbSrc *MongoDB) *MongoRoomStore {
	return &MongoRoomStore{
		db:     dbSrc,
		dbColl: dbSrc.Collection(dbRoomsCollectionName),
	}
}

func (self *MongoRoomStore) GetByID(
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

func (self *MongoRoomStore) Create(
	ctx context.Context, room *types.Room,
) (primitive.ObjectID, error) {
	result, err := self.dbColl.InsertOne(ctx, room)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.ObjectID{}, fmt.Errorf("Failed to cast %v to id", result.InsertedID)
	}
	return insertedID, nil
}

func (self *MongoRoomStore) Get(
	ctx context.Context, query []byte,
) ([]*types.Room, error) {
	cursor, err := self.dbColl.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	rooms := []*types.Room{}

	err = cursor.All(ctx, &rooms)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (self *MongoRoomStore) UpdateByID(
	ctx context.Context, id primitive.ObjectID, room *types.Room,
) error {
	_, err := self.dbColl.UpdateByID(
		ctx, id, bson.M{"$set": room},
	)
	if err != nil {
		return err
	}
	return nil
}

func (self *MongoRoomStore) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := self.dbColl.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
