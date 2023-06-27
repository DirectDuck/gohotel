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

const usersCollectionName = "users"

type Dropper interface {
	Drop(context.Context) error
}

type UserStore interface {
	Dropper

	CreateUser(context.Context, *types.User) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	GetUserByID(context.Context, primitive.ObjectID) (*types.User, error)
	UpdateUserByID(context.Context, primitive.ObjectID, *types.User) (*types.User, error)
	DeleteUserByID(context.Context, primitive.ObjectID) error
}

type MongoUserStore struct {
	db     *mongo.Database
	dbColl *mongo.Collection
}

func NewMongoUserStore(db *mongo.Database) *MongoUserStore {
	return &MongoUserStore{
		db:     db,
		dbColl: db.Collection(usersCollectionName),
	}
}

func (self *MongoUserStore) GetUserByID(ctx context.Context, id primitive.ObjectID) (*types.User, error) {
	user := &types.User{}

	err := self.dbColl.FindOne(
		ctx, bson.M{"_id": id},
	).Decode(user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (self *MongoUserStore) UpdateUserByID(
	ctx context.Context, id primitive.ObjectID, data *types.User,
) (*types.User, error) {

	_, err := self.dbColl.UpdateByID(
		ctx, id, bson.M{"$set": data},
	)
	if err != nil {
		return nil, err
	}

	user, err := self.GetUserByID(ctx, id)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (self *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cursor, err := self.dbColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*types.User

	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (self *MongoUserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	result, err := self.dbColl.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("Failed to cast %v to id", result.InsertedID)
	}
	return self.GetUserByID(ctx, insertedID)
}

func (self *MongoUserStore) DeleteUserByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := self.dbColl.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (self *MongoUserStore) Drop(ctx context.Context) error {
	err := self.db.Drop(ctx)
	if err != nil {
		return err
	}
	return nil
}
