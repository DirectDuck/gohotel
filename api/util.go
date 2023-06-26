package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func HandleAPIError(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	fiberErr := ctx.Status(code).JSON(map[string]interface{}{
		"error": err.Error(),
	})

	if fiberErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return nil
}
