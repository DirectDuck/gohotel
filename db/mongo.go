package db

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStore struct {
	Coll *mongo.Collection
}

func (self *MongoStore) Create(
	ctx context.Context, value interface{},
) (primitive.ObjectID, error) {
	result, err := self.Coll.InsertOne(ctx, value)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.ObjectID{}, fmt.Errorf("Failed to cast %v to id", result.InsertedID)
	}
	return insertedID, nil
}

func (self *MongoStore) Get(ctx context.Context, query interface{}, castTo interface{}) (interface{}, error) {
	cursor, err := self.Coll.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	objs := castTo

	err = cursor.All(ctx, &objs)
	if err != nil {
		return nil, err
	}

	return objs, nil
}

func (self *MongoStore) GetCount(ctx context.Context, query interface{}) (int64, error) {
	return self.Coll.CountDocuments(ctx, query)
}

func (self *MongoStore) GetOne(
	ctx context.Context, query interface{}, castTo interface{},
) (interface{}, error) {
	obj := castTo

	err := self.Coll.FindOne(ctx, query).Decode(obj)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return obj, nil
}

func (self *MongoStore) GetOneByID(
	ctx context.Context, id primitive.ObjectID, castTo interface{},
) (interface{}, error) {
	obj := castTo

	err := self.Coll.FindOne(ctx, bson.M{"_id": id}).Decode(obj)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return obj, nil
}

func (self *MongoStore) UpdateByID(
	ctx context.Context, id primitive.ObjectID, objs interface{},
) error {
	_, err := self.Coll.UpdateByID(
		ctx, id, bson.M{"$set": objs},
	)
	if err != nil {
		return err
	}

	return nil
}

func (self *MongoStore) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := self.Coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
