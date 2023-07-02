package db

import (
	"context"
	"errors"
	"fmt"
	"hotel/types"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const dbUsersCollectionName = "users"

type UserStore interface {
	Login(context.Context, *types.LoginUserParams) (string, *types.User, error)
	Create(context.Context, *types.User) (*types.User, error)
	Get(context.Context) ([]*types.User, error)
	GetByID(context.Context, primitive.ObjectID) (*types.User, error)
	UpdateByID(context.Context, primitive.ObjectID, *types.User) (*types.User, error)
	DeleteByID(context.Context, primitive.ObjectID) error
}

type MongoUserStore struct {
	db     *MongoDB
	dbColl *mongo.Collection
}

func NewMongoUserStore(dbSrc *MongoDB) *MongoUserStore {
	return &MongoUserStore{
		db:     dbSrc,
		dbColl: dbSrc.Collection(dbUsersCollectionName),
	}
}

func (self *MongoUserStore) Login(
	ctx context.Context, params *types.LoginUserParams,
) (string, *types.User, error) {
	user := &types.User{}

	err := self.dbColl.FindOne(
		ctx, bson.M{"email": params.Email},
	).Decode(user)

	if err != nil {
		return "", nil, err
	}

	if !user.CheckPasswordValid(params.Password) {
		return "", nil, fmt.Errorf("Invalid credentials")
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", nil, fmt.Errorf("Failed to sign token: %s", err.Error())
	}

	return tokenStr, user, nil
}

func (self *MongoUserStore) GetByID(ctx context.Context, id primitive.ObjectID) (*types.User, error) {
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

func (self *MongoUserStore) UpdateByID(
	ctx context.Context, id primitive.ObjectID, user *types.User,
) (*types.User, error) {
	userBefore, err := self.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	errs := user.Validate(userBefore)
	if len(errs) != 0 {
		return nil, ValidationError{Fields: errs}
	}
	err = user.Evaluate(userBefore)
	if err != nil {
		return nil, err
	}

	_, err = self.dbColl.UpdateByID(
		ctx, id, bson.M{"$set": user},
	)
	if err != nil {
		return nil, err
	}

	updatedUser, err := self.GetByID(ctx, id)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return updatedUser, nil
}

func (self *MongoUserStore) Get(ctx context.Context) ([]*types.User, error) {
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

func (self *MongoUserStore) Create(ctx context.Context, user *types.User) (*types.User, error) {
	errs := user.Validate(nil)
	if len(errs) != 0 {
		return nil, ValidationError{Fields: errs}
	}
	err := user.Evaluate(nil)
	if err != nil {
		return nil, err
	}
	result, err := self.dbColl.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("Failed to cast %v to id", result.InsertedID)
	}
	return self.GetByID(ctx, insertedID)
}

func (self *MongoUserStore) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := self.dbColl.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
