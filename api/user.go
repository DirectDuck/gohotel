package api

import (
	"context"
	"hotel/db"
	"hotel/types"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (self *UserHandler) HandleListUsers(c *fiber.Ctx) error {
	user1 := types.User{
		ID:        "123",
		FirstName: "James",
		LastName:  "Mister",
	}
	user2 := types.User{
		ID:        "124",
		FirstName: "Albert",
		LastName:  "Second",
	}
	return c.JSON([]types.User{user1, user2})
}

func (self *UserHandler) HandleGetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := self.userStore.GetUserByID(context.Background(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.JSON(user)
}
