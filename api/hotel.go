package api

import (
	"hotel/controllers"
	"hotel/db"
	"hotel/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	controller *controllers.HotelController
}

func NewHotelHandler(controller *controllers.HotelController) *HotelHandler {
	return &HotelHandler{
		controller: controller,
	}
}

func (self *HotelHandler) HandleListHotels(ctx *fiber.Ctx) error {
	hotels, err := self.controller.Get(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(hotels)
}

func (self *HotelHandler) HandleGetHotel(ctx *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		return err
	}

	hotel, err := self.controller.GetWithRoomsByID(ctx.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if hotel == nil {
		return fiber.NewError(fiber.StatusNotFound, EntityNotFoundMessage)
	}

	return ctx.JSON(hotel)
}

func (self *HotelHandler) HandleCreateHotel(ctx *fiber.Ctx) error {
	var params types.CreateHotelParams
	err := ctx.BodyParser(&params)
	if err != nil {
		return err
	}

	hotel, err := types.NewHotelFromCreateParams(params)
	if err != nil {
		return err
	}

	createdHotel, err := self.controller.Create(ctx.Context(), hotel)
	if err != nil {
		validationError, ok := err.(db.ValidationError)
		if ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationError.Fields)
		}
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(createdHotel)
}

func (self *HotelHandler) HandleUpdateHotel(ctx *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		return err
	}

	var params types.UpdateHotelParams
	err = ctx.BodyParser(&params)
	if err != nil {
		return err
	}

	data, err := types.NewHotelFromUpdateParams(params)
	if err != nil {
		return err
	}

	updatedHotel, err := self.controller.UpdateByID(ctx.Context(), id, data)
	if err != nil {
		validationError, ok := err.(db.ValidationError)
		if ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationError.Fields)
		}
		return err
	}
	if updatedHotel == nil {
		return fiber.NewError(fiber.StatusNotFound, EntityNotFoundMessage)
	}

	return ctx.JSON(updatedHotel)
}

func (self *HotelHandler) HandleDeleteHotel(ctx *fiber.Ctx) error {
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
