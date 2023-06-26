package api

import (
	"hotel/db"
	"hotel/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (self *UserHandler) HandleListUsers(ctx *fiber.Ctx) error {
	users, err := self.userStore.GetUsers(ctx.Context())
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

	user, err := self.userStore.GetUserByID(ctx.Context(), id)
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

	errs := params.Validate()
	if len(errs) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(errs)
	}

	user, err := types.NewUserFromCreateParams(params)
	if err != nil {
		return err
	}

	createdUser, err := self.userStore.CreateUser(ctx.Context(), user)
	if err != nil {
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

	errs := params.Validate()
	if len(errs) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(errs)
	}

	data, err := types.NewUserFromUpdateParams(params)
	if err != nil {
		return err
	}

	updatedUser, err := self.userStore.UpdateUserByID(ctx.Context(), id, data)
	if err != nil {
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

	err = self.userStore.DeleteUserByID(ctx.Context(), id)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusNoContent).Send(nil)
}
