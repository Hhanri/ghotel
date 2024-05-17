package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/db"
	"github.com/hhanri/ghotel/types"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetAll(c.Context(), struct{}{})
	if err != nil {
		return FiberInternalErrorResponse(c)
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetByID(c.Context(), id)
	if err != nil {
		return FiberNotFoundErrorResponse(c)
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return FiberInternalErrorResponse(c)
	}
	if user.ID != booking.UserID {
		return FiberUnauthorizedErrorResponse(c)
	}

	return c.JSON(booking)
}
