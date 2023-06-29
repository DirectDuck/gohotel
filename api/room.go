package api

import (
	"hotel/db"
	"hotel/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	roomStore db.RoomStore
}

func NewRoomHandler(roomStore db.RoomStore) *RoomHandler {
	return &RoomHandler{
		roomStore: roomStore,
	}
}

func (self *RoomHandler) HandleListRooms(ctx *fiber.Ctx) error {
	var query db.ListRoomsQueryParams
	err := ctx.QueryParser(&query)
	if err != nil {
		return err
	}
	rooms, err := self.roomStore.Get(ctx.Context(), &query)
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

	room, err := self.roomStore.GetByID(ctx.Context(), id)
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

	errs := params.Validate()
	if len(errs) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(errs)
	}

	room, err := types.NewRoomFromCreateParams(params)
	if err != nil {
		return err
	}

	createdRoom, err := self.roomStore.Create(ctx.Context(), room)
	if err != nil {
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

	errs := params.Validate()
	if len(errs) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(errs)
	}

	data, err := types.NewRoomFromUpdateParams(params)
	if err != nil {
		return err
	}

	updatedRoom, err := self.roomStore.UpdateByID(ctx.Context(), id, data)
	if err != nil {
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

	err = self.roomStore.DeleteByID(ctx.Context(), id)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusNoContent).Send(nil)
}
