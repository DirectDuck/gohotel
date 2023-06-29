package api

import (
	"hotel/db"
	"hotel/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (self *HotelHandler) HandleListHotels(ctx *fiber.Ctx) error {
	hotels, err := self.store.Hotels.Get(ctx.Context())
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

	hotel, err := self.store.Hotels.GetByID(ctx.Context(), id)
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

	errs := params.Validate()
	if len(errs) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(errs)
	}

	hotel, err := types.NewHotelFromCreateParams(params)
	if err != nil {
		return err
	}

	createdHotel, err := self.store.Hotels.Create(ctx.Context(), hotel)
	if err != nil {
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

	errs := params.Validate()
	if len(errs) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(errs)
	}

	data, err := types.NewHotelFromUpdateParams(params)
	if err != nil {
		return err
	}

	updatedHotel, err := self.store.Hotels.UpdateByID(ctx.Context(), id, data)
	if err != nil {
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

	err = self.store.Hotels.DeleteByID(ctx.Context(), id)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusNoContent).Send(nil)
}
