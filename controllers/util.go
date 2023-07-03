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

func GetUserFromContext(store *db.Store, ctx context.Context) (*types.User, error) {
	ctxVal := ctx.Value("user")
	if ctxVal == nil {
		return nil, nil
	}
	claims := ctxVal.(*jwt.Token).Claims.(jwt.MapClaims)
	idStr := claims["id"]
	if idStr == nil {
		return nil, nil
	}
	id, err := primitive.ObjectIDFromHex(idStr.(string))
	if err != nil {
		return nil, err
	}
	return store.Users.GetByID(ctx, id)
}

func IsEmailValid(e string) bool {
	// Sourced from https://stackoverflow.com/a/67686133
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
