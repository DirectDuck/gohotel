package types

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12

	minUserFirstNameLen = 2
	minUserLastNameLen  = 2
	minUserPasswordLen  = 7
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}

type BaseUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (self *BaseUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(self.FirstName) < minUserFirstNameLen {
		errors["firstName"] = fmt.Sprintf(
			"First name length should be at least %d characters", minUserFirstNameLen,
		)
	}

	if len(self.LastName) < minUserLastNameLen {
		errors["lastName"] = fmt.Sprintf(
			"Last name length should be at least %d characters", minUserLastNameLen,
		)
	}

	if !IsEmailValid(self.Email) {
		errors["email"] = fmt.Sprintf(
			"Email \"%s\" is invalid", self.Email,
		)
	}
	return errors
}

type CreateUserParams struct {
	BaseUserParams
	Password string `json:"password"`
}

func (self *CreateUserParams) Validate() map[string]string {
	errors := self.BaseUserParams.Validate()

	if len(self.Password) < minUserPasswordLen {
		errors["password"] = fmt.Sprintf(
			"Password length should be at least %d characters", minUserPasswordLen,
		)
	}

	return errors
}

type UpdateUserParams struct {
	BaseUserParams
}

func (self *UpdateUserParams) Validate() map[string]string {
	errors := self.BaseUserParams.Validate()
	return errors
}

func NewUserFromCreateParams(params CreateUserParams) (*User, error) {
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

func NewUserFromUpdateParams(params UpdateUserParams) (*User, error) {
	return &User{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
	}, nil
}
