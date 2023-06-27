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
	CreateHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	GetHotelByID(context.Context, primitive.ObjectID) (*types.Hotel, error)
}

type MongoHotelStore struct {
	dbSrc  *mongo.Database
	dbColl *mongo.Collection
}

func NewMongoHotelStore(dbSrc *mongo.Database) *MongoHotelStore {
	return &MongoHotelStore{
		dbSrc:  dbSrc,
		dbColl: dbSrc.Collection(dbHotelsCollectionName),
	}
}

func (self *MongoHotelStore) GetHotelByID(
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

func (self *MongoHotelStore) CreateHotel(
	ctx context.Context, hotel *types.Hotel,
) (*types.Hotel, error) {
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
