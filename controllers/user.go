package controllers

import (
	"context"
	"fmt"
	"hotel/types"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12

	minUserFirstNameLen = 2
	minUserLastNameLen  = 2
	minUserPasswordLen  = 7
)

type UserController struct {
	Store *Store
}

func (self *UserController) CheckPasswordValid(user *types.User, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password)) == nil
}

func (self *UserController) Login(
	ctx context.Context, params *types.LoginUserParams,
) (string, *types.User, error) {
	query, err := bson.Marshal(bson.M{"email": params.Email})
	if err != nil {
		return "", nil, err
	}
	result, err := self.Store.DB.Users.GetOne(ctx, query, &types.User{})
	if err != nil {
		return "", nil, err
	}
	user := CastPtrInterface[types.User](result)

	if user == nil || !self.CheckPasswordValid(user, params.Password) {
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

func (self *UserController) GetByID(
	ctx context.Context, id primitive.ObjectID,
) (*types.User, error) {
	result, err := self.Store.DB.Users.GetOneByID(ctx, id, &types.User{})
	if err != nil {
		return nil, err
	}
	return CastPtrInterface[types.User](result), nil
}

func (self *UserController) Get(ctx context.Context) ([]*types.User, error) {
	result, err := self.Store.DB.Users.Get(ctx, bson.M{}, []*types.User{})
	if err != nil {
		return nil, err
	}
	return CastInterface[[]*types.User](result), nil
}

func (self *UserController) Validate(user *types.User, userBefore *types.User) map[string]string {
	errors := map[string]string{}
	if len(user.FirstName) < minUserFirstNameLen {
		errors["firstName"] = fmt.Sprintf(
			"First name length should be at least %d characters", minUserFirstNameLen,
		)
	}

	if len(user.LastName) < minUserLastNameLen {
		errors["lastName"] = fmt.Sprintf(
			"Last name length should be at least %d characters", minUserLastNameLen,
		)
	}

	if !IsEmailValid(user.Email) {
		errors["email"] = fmt.Sprintf(
			"Email \"%s\" is invalid", user.Email,
		)
	}

	if userBefore == nil {
		if len(user.Password) < minUserPasswordLen {
			errors["password"] = fmt.Sprintf(
				"Password length should be at least %d characters", minUserPasswordLen,
			)
		}
	}
	return errors
}

func (self *UserController) Evaluate(user *types.User, userBefore *types.User) error {
	if userBefore == nil {
		encryptedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(user.Password), bcryptCost,
		)
		if err != nil {
			return err
		}
		user.EncryptedPassword = string(encryptedPassword)
	}
	return nil
}

func (self *UserController) Create(
	ctx context.Context, user *types.User,
) (*types.User, error) {
	errs := self.Validate(user, nil)
	if len(errs) != 0 {
		return nil, ValidationError{Fields: errs}
	}
	err := self.Evaluate(user, nil)
	if err != nil {
		return nil, err
	}
	id, err := self.Store.DB.Users.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return self.GetByID(ctx, id)
}

func (self *UserController) UpdateByID(
	ctx context.Context, id primitive.ObjectID, user *types.User,
) (*types.User, error) {
	userBefore, err := self.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	errs := self.Validate(user, userBefore)
	if len(errs) != 0 {
		return nil, ValidationError{Fields: errs}
	}
	err = self.Evaluate(user, userBefore)
	if err != nil {
		return nil, err
	}

	err = self.Store.DB.Users.UpdateByID(ctx, id, user)
	if err != nil {
		return nil, err
	}
	return self.GetByID(ctx, id)
}

func (self *UserController) DeleteByID(
	ctx context.Context, id primitive.ObjectID,
) error {
	return self.Store.DB.Users.DeleteByID(ctx, id)
}
