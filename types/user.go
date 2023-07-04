package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	IsAdmin           bool               `bson:"isAdmin" json:"-"`
	Password          string             `bson:"-" json:"-"`
	EncryptedPassword string             `bson:"encryptedPassword,omitempty" json:"-"`
}

type LoginUserParams struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"-" json:"password"`
}

type BaseUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type CreateUserParams struct {
	BaseUserParams
	Password string `json:"password"`
}

type UpdateUserParams struct {
	BaseUserParams
}

func NewUserFromCreateParams(params CreateUserParams) (*User, error) {
	return &User{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Password:  params.Password,
	}, nil
}

func NewUserFromUpdateParams(params UpdateUserParams) (*User, error) {
	return &User{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
	}, nil
}
