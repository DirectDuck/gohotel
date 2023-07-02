package api

import (
	"hotel/db"
	"hotel/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (self *BookingHandler) HandleListBookings(ctx *fiber.Ctx) error {
	rooms, err := self.store.Bookings.Get(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(rooms)
}

func (self *BookingHandler) HandleGetBooking(ctx *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		return err
	}

	room, err := self.store.Bookings.GetUnfoldedByID(ctx.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if room == nil {
		return fiber.NewError(fiber.StatusNotFound, EntityNotFoundMessage)
	}

	return ctx.JSON(room)
}

func (self *BookingHandler) HandleCreateBooking(ctx *fiber.Ctx) error {
	var params types.CreateBookingParams
	err := ctx.BodyParser(&params)
	if err != nil {
		return err
	}

	room, err := types.NewBookingFromCreateParams(params)
	if err != nil {
		return err
	}

	createdBooking, err := self.store.Bookings.Create(ctx.Context(), room)
	if err != nil {
		validationError, ok := err.(db.ValidationError)
		if ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationError.Fields)
		}
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(createdBooking)
}

func (self *BookingHandler) HandleUpdateBooking(ctx *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		return err
	}

	var params types.UpdateBookingParams
	err = ctx.BodyParser(&params)
	if err != nil {
		return err
	}

	data, err := types.NewBookingFromUpdateParams(params)
	if err != nil {
		return err
	}

	updatedBooking, err := self.store.Bookings.UpdateByID(ctx.Context(), id, data)
	if err != nil {
		validationError, ok := err.(db.ValidationError)
		if ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationError.Fields)
		}
		return err
	}
	if updatedBooking == nil {
		return fiber.NewError(fiber.StatusNotFound, EntityNotFoundMessage)
	}

	return ctx.JSON(updatedBooking)
}

func (self *BookingHandler) HandleDeleteBooking(ctx *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		return err
	}

	err = self.store.Bookings.DeleteByID(ctx.Context(), id)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusNoContent).Send(nil)
}
