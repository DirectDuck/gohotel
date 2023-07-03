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

const dbBookingsCollectionName = "bookings"

type BookingStore interface {
	Create(context.Context, *types.Booking) (primitive.ObjectID, error)
	GetByID(context.Context, primitive.ObjectID) (*types.Booking, error)
	Get(context.Context, []byte) ([]*types.Booking, error)
	UpdateByID(context.Context, primitive.ObjectID, *types.Booking) error
	DeleteByID(context.Context, primitive.ObjectID) error
}

type MongoBookingStore struct {
	db     *MongoDB
	dbColl *mongo.Collection
}

func NewMongoBookingStore(dbSrc *MongoDB) *MongoBookingStore {
	return &MongoBookingStore{
		db:     dbSrc,
		dbColl: dbSrc.Collection(dbBookingsCollectionName),
	}
}

func (self *MongoBookingStore) GetByID(
	ctx context.Context, id primitive.ObjectID,
) (*types.Booking, error) {
	booking := &types.Booking{}

	err := self.dbColl.FindOne(
		ctx, bson.M{"_id": id},
	).Decode(booking)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return booking, nil
}

func (self *MongoBookingStore) Create(
	ctx context.Context, booking *types.Booking,
) (primitive.ObjectID, error) {
	result, err := self.dbColl.InsertOne(ctx, booking)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.ObjectID{}, fmt.Errorf("Failed to cast %v to id", result.InsertedID)
	}
	return insertedID, nil
}

func (self *MongoBookingStore) Get(ctx context.Context, query []byte) ([]*types.Booking, error) {
	cursor, err := self.dbColl.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	bookings := []*types.Booking{}

	err = cursor.All(ctx, &bookings)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (self *MongoBookingStore) UpdateByID(
	ctx context.Context, id primitive.ObjectID, booking *types.Booking,
) error {
	_, err := self.dbColl.UpdateByID(
		ctx, id, bson.M{"$set": booking},
	)
	if err != nil {
		return err
	}

	return nil
}

func (self *MongoBookingStore) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := self.dbColl.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
