package api

import (
	"hotel/types"

	"github.com/gofiber/fiber/v2"
)

func HandleListUsers(c *fiber.Ctx) error {
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

func HandleGetUser(c *fiber.Ctx) error {
	user := types.User{
		ID:        "123",
		FirstName: "James",
		LastName:  "Mister",
	}
	return c.JSON(user)
}
