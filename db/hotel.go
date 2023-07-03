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
	Create(context.Context, *types.Hotel) (primitive.ObjectID, error)
	Get(context.Context) ([]*types.Hotel, error)
	GetByID(context.Context, primitive.ObjectID) (*types.Hotel, error)
	UpdateByID(context.Context, primitive.ObjectID, *types.Hotel) error
	DeleteByID(context.Context, primitive.ObjectID) error
}

type MongoHotelStore struct {
	db     *MongoDB
	dbColl *mongo.Collection
}

func NewMongoHotelStore(dbSrc *MongoDB) *MongoHotelStore {
	return &MongoHotelStore{
		db:     dbSrc,
		dbColl: dbSrc.Collection(dbHotelsCollectionName),
	}
}

func (self *MongoHotelStore) Create(
	ctx context.Context, hotel *types.Hotel,
) (primitive.ObjectID, error) {
	result, err := self.dbColl.InsertOne(ctx, hotel)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.ObjectID{}, fmt.Errorf("Failed to cast %v to id", result.InsertedID)
	}
	return insertedID, err
}

func (self *MongoHotelStore) Get(ctx context.Context) ([]*types.Hotel, error) {
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

func (self *MongoHotelStore) GetByID(
	ctx context.Context, id primitive.ObjectID,
) (*types.Hotel, error) {
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

	return hotel, nil
}

func (self *MongoHotelStore) UpdateByID(
	ctx context.Context, id primitive.ObjectID, hotel *types.Hotel,
) error {
	_, err := self.dbColl.UpdateByID(
		ctx, id, bson.M{"$set": hotel},
	)
	if err != nil {
		return err
	}
	return nil
}

func (self *MongoHotelStore) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := self.dbColl.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
