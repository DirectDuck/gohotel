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
	Create(context.Context, *types.Booking) (*types.BookingUnfolded, error)
	GetByID(context.Context, primitive.ObjectID) (*types.Booking, error)
	GetUnfoldedByID(context.Context, primitive.ObjectID) (*types.BookingUnfolded, error)
	Get(context.Context) ([]*types.Booking, error)
	UpdateByID(context.Context, primitive.ObjectID, *types.Booking) (*types.BookingUnfolded, error)
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

func BookingToUnfolded(ctx context.Context, booking *types.Booking, store *Store) (*types.BookingUnfolded, error) {
	room, err := store.Rooms.GetByID(ctx, booking.RoomID)
	if err != nil {
		return nil, err
	}
	user, err := store.Users.GetByID(ctx, booking.UserID)
	if err != nil {
		return nil, err
	}

	return &types.BookingUnfolded{
		Booking: booking,
		Room:    room,
		User:    user,
	}, nil
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

func (self *MongoBookingStore) GetUnfoldedByID(
	ctx context.Context, id primitive.ObjectID,
) (*types.BookingUnfolded, error) {
	booking, err := self.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return BookingToUnfolded(ctx, booking, self.db.Store)
}

func (self *MongoBookingStore) Create(
	ctx context.Context, booking *types.Booking,
) (*types.BookingUnfolded, error) {
	bookingUnfolded, err := BookingToUnfolded(
		ctx, booking, self.db.Store,
	)
	if err != nil {
		return nil, err
	}
	errs := bookingUnfolded.Validate(nil)
	if len(errs) != 0 {
		return nil, ValidationError{Fields: errs}
	}
	err = bookingUnfolded.Evaluate(nil)
	if err != nil {
		return nil, err
	}

	result, err := self.dbColl.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("Failed to cast %v to id", result.InsertedID)
	}
	return self.GetUnfoldedByID(ctx, insertedID)
}

func (self *MongoBookingStore) Get(ctx context.Context) ([]*types.Booking, error) {
	cursor, err := self.dbColl.Find(ctx, bson.M{})
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
) (*types.BookingUnfolded, error) {
	bookingUnfolded, err := BookingToUnfolded(
		ctx, booking, self.db.Store,
	)
	if err != nil {
		return nil, err
	}

	beforeUnfolded, err := self.GetUnfoldedByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	errs := bookingUnfolded.Validate(beforeUnfolded)
	if len(errs) != 0 {
		return nil, ValidationError{Fields: errs}
	}
	err = bookingUnfolded.Evaluate(beforeUnfolded)
	if err != nil {
		return nil, err
	}

	_, err = self.dbColl.UpdateByID(
		ctx, id, bson.M{"$set": bookingUnfolded.Booking},
	)
	if err != nil {
		return nil, err
	}

	updatedUnfolded, err := self.GetUnfoldedByID(ctx, id)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return updatedUnfolded, nil
}

func (self *MongoBookingStore) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := self.dbColl.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
