package controllers

import (
	"context"
	"fmt"
	"hotel/db"
	"hotel/types"
	"regexp"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ValidationError struct {
	Fields map[string]string
}

func (self ValidationError) Error() string {
	errStr := "ValidationError:"
	for k, v := range self.Fields {
		errStr += fmt.Sprintf("\n%s: %s", k, v)
	}
	return errStr
}

func GetUserIDFromContext(dbStore *db.DB, ctx context.Context) (primitive.ObjectID, error) {
	ctxVal := ctx.Value("user")
	if ctxVal == nil {
		return primitive.ObjectID{}, nil
	}
	claims := ctxVal.(*jwt.Token).Claims.(jwt.MapClaims)
	idStr := claims["id"]
	if idStr == nil {
		return primitive.ObjectID{}, nil
	}
	id, err := primitive.ObjectIDFromHex(idStr.(string))
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return id, err
}

func GetUserFromContext(dbStore *db.DB, ctx context.Context) (*types.User, error) {
	id, err := GetUserIDFromContext(dbStore, ctx)
	if err != nil {
		return nil, err
	}
	user, err := dbStore.Users.GetOneByID(ctx, id, &types.User{})
	if err != nil {
		return nil, err
	}
	return user.(*types.User), nil
}

func IsEmailValid(e string) bool {
	// Sourced from https://stackoverflow.com/a/67686133
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func CastPtrInterface[T any](i interface{}) *T {
	casted, ok := i.(*T)
	if !ok {
		return nil
	}
	return casted
}

func CastInterface[T any](i interface{}) T {
	casted, ok := i.(T)
	if !ok {
		return *new(T)
	}
	return casted
}
