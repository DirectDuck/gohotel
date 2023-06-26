package types

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12

	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params *CreateUserParams) Validate() error {
	if len(params.FirstName) < minFirstNameLen {
		return fmt.Errorf("First name length should be at least %d characters", minFirstNameLen)
	}

	if len(params.LastName) < minLastNameLen {
		return fmt.Errorf("Last name length should be at least %d characters", minLastNameLen)
	}

	if len(params.Password) < minPasswordLen {
		return fmt.Errorf("Password length should be at least %d characters", minPasswordLen)
	}

	if !IsEmailValid(params.Email) {
		return fmt.Errorf("Email \"%s\" is invalid", params.Email)
	}

	return nil
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(params.Password), bcryptCost,
	)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encryptedPassword),
	}, nil
}
