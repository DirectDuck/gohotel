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
	Password          string             `bson:"-" json:"-"`
	EncryptedPassword string             `bson:"encryptedPassword,omitempty" json:"-"`
}

func (self *User) Validate(dbBefore *User) map[string]string {
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

	if dbBefore == nil {
		if len(self.Password) < minUserPasswordLen {
			errors["password"] = fmt.Sprintf(
				"Password length should be at least %d characters", minUserPasswordLen,
			)
		}
	}

	return errors
}

func (self *User) Evaluate(dbBefore *User) error {
	if dbBefore == nil {
		encryptedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(self.Password), bcryptCost,
		)
		if err != nil {
			return err
		}
		self.EncryptedPassword = string(encryptedPassword)
	}
	return nil
}

type LoginUserParams struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"-" json:"password"`
}

func (self *User) CheckPasswordValid(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(self.EncryptedPassword), []byte(password)) == nil
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
