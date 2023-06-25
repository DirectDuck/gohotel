package db

import (
	"context"
	"hotel/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore interface {
	GetUserByID(context.Context, string) (*types.User, error)
}

type MongoUserStore struct {
	dbColl mongo.Collection
}

func NewMongoUserStore(dbColl *mongo.Collection) *MongoUserStore {
	return &MongoUserStore{
		dbColl: *dbColl,
	}
}

func (self *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	user := &types.User{}
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	err = self.dbColl.FindOne(
		ctx, bson.M{"_id": oid},
	).Decode(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
