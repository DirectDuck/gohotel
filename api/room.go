package api

import (
	"hotel/controllers"
	"hotel/db"
	"hotel/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	controller *controllers.RoomController
}

func NewRoomHandler(controller *controllers.RoomController) *RoomHandler {
	return &RoomHandler{
		controller: controller,
	}
}

func (self *RoomHandler) HandleListRooms(ctx *fiber.Ctx) error {
	var query controllers.RoomGetQueryParams
	err := ctx.QueryParser(&query)
	if err != nil {
		return err
	}
	rooms, err := self.controller.Get(ctx.Context(), &query)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(rooms)
}

func (self *RoomHandler) HandleGetRoom(ctx *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		return err
	}

	room, err := self.controller.GetUnfoldedByID(ctx.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if room == nil {
		return fiber.NewError(fiber.StatusNotFound, EntityNotFoundMessage)
	}

	return ctx.JSON(room)
}

func (self *RoomHandler) HandleCreateRoom(ctx *fiber.Ctx) error {
	var params types.CreateRoomParams
	err := ctx.BodyParser(&params)
	if err != nil {
		return err
	}

	room, err := types.NewRoomFromCreateParams(params)
	if err != nil {
		return err
	}

	createdRoom, err := self.controller.Create(ctx.Context(), room)
	if err != nil {
		validationError, ok := err.(db.ValidationError)
		if ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationError.Fields)
		}
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(createdRoom)
}

func (self *RoomHandler) HandleUpdateRoom(ctx *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		return err
	}

	var params types.UpdateRoomParams
	err = ctx.BodyParser(&params)
	if err != nil {
		return err
	}

	data, err := types.NewRoomFromUpdateParams(params)
	if err != nil {
		return err
	}

	updatedRoom, err := self.controller.UpdateByID(ctx.Context(), id, data)
	if err != nil {
		validationError, ok := err.(db.ValidationError)
		if ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationError.Fields)
		}
		return err
	}
	if updatedRoom == nil {
		return fiber.NewError(fiber.StatusNotFound, EntityNotFoundMessage)
	}

	return ctx.JSON(updatedRoom)
}

func (self *RoomHandler) HandleDeleteRoom(ctx *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		return err
	}

	err = self.controller.DeleteByID(ctx.Context(), id)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusNoContent).Send(nil)
}
