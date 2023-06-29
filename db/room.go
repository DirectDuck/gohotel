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
	Create(context.Context, *types.Room) (*types.Room, error)
	GetByID(context.Context, primitive.ObjectID) (*types.Room, error)
	Get(context.Context, *ListRoomsQueryParams) ([]*types.Room, error)
	GetForHotel(context.Context, primitive.ObjectID) ([]*types.Room, error)
	UpdateByID(context.Context, primitive.ObjectID, *types.Room) (*types.Room, error)
	DeleteByID(context.Context, primitive.ObjectID) error
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
) (*types.Room, error) {
	result, err := self.dbColl.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("Failed to cast %v to id", result.InsertedID)
	}
	return self.GetByID(ctx, insertedID)
}

type ListRoomsQueryParams struct {
	HotelID primitive.ObjectID `bson:"hotelID,omitempty" json:"hotelID"`
}

func (self *MongoRoomStore) Get(
	ctx context.Context, query *ListRoomsQueryParams,
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

func (self *MongoRoomStore) GetForHotel(
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

func (self *MongoRoomStore) UpdateByID(
	ctx context.Context, id primitive.ObjectID, data *types.Room,
) (*types.Room, error) {

	_, err := self.dbColl.UpdateByID(
		ctx, id, bson.M{"$set": data},
	)
	if err != nil {
		return nil, err
	}

	room, err := self.GetByID(ctx, id)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return room, nil
}

func (self *MongoRoomStore) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := self.dbColl.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
