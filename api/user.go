package api

import (
	"hotel/db"
	"hotel/types"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	store *db.Store
}

func NewUserHandler(store *db.Store) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (self *UserHandler) HandleLogin(ctx *fiber.Ctx) error {
	var params types.LoginUserParams
	err := ctx.BodyParser(&params)
	if err != nil {
		return err
	}

	token, user, err := self.store.Users.Login(ctx.Context(), &params)
	if err != nil {
		log.Printf("Auth failed: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Auth failed")
	}

	// Using just fiber.Map is intentional, for learning purposes
	// otherwise it would've been a struct
	return ctx.JSON(fiber.Map{
		"token": token,
		"user":  user,
	})
}

func (self *UserHandler) HandleListUsers(ctx *fiber.Ctx) error {
	users, err := self.store.Users.Get(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(users)
}

func (self *UserHandler) HandleGetUser(ctx *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		return err
	}

	user, err := self.store.Users.GetByID(ctx.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if user == nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return ctx.JSON(user)
}

func (self *UserHandler) HandleCreateUser(ctx *fiber.Ctx) error {
	var params types.CreateUserParams
	err := ctx.BodyParser(&params)
	if err != nil {
		return err
	}

	user, err := types.NewUserFromCreateParams(params)
	if err != nil {
		return err
	}

	createdUser, err := self.store.Users.Create(ctx.Context(), user)
	if err != nil {
		validationError, ok := err.(db.ValidationError)
		if ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationError.Fields)
		}
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(createdUser)
}

func (self *UserHandler) HandleUpdateUser(ctx *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		return err
	}

	var params types.UpdateUserParams
	err = ctx.BodyParser(&params)
	if err != nil {
		return err
	}

	data, err := types.NewUserFromUpdateParams(params)
	if err != nil {
		return err
	}

	updatedUser, err := self.store.Users.UpdateByID(ctx.Context(), id, data)
	if err != nil {
		validationError, ok := err.(db.ValidationError)
		if ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationError.Fields)
		}
		return err
	}
	if updatedUser == nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return ctx.JSON(updatedUser)
}

func (self *UserHandler) HandleDeleteUser(ctx *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		return err
	}

	err = self.store.Users.DeleteByID(ctx.Context(), id)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusNoContent).Send(nil)
}
