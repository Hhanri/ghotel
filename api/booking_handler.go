package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/api/api_error"
	"github.com/hhanri/ghotel/api/api_util"
	"github.com/hhanri/ghotel/db"
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
		return api_error.InternalErrorResponse
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetByID(c.Context(), id)
	if err != nil {
		return api_error.NotFoundErrorResponse
	}

	user, err := api_util.GetAuth(c.Context())
	if err != nil {
		return api_error.UnauthorizedErrorResponse
	}
	if user.ID != booking.UserID {
		return api_error.UnauthorizedErrorResponse
	}

	return c.JSON(booking)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	booking, err := h.store.Booking.GetByID(c.Context(), id)
	if err != nil {
		return api_error.InternalErrorResponse
	}
	user, err := api_util.GetAuth(c.Context())
	if err != nil {
		return api_error.UnauthorizedErrorResponse
	}
	if booking.UserID != user.ID && !user.IsAdmin {
		return api_error.UnauthorizedErrorResponse
	}

	err = h.store.Booking.Cancel(c.Context(), id)
	if err != nil {
		fmt.Println(err)
		return api_error.InternalErrorResponse
	}

	return c.JSON(
		map[string]string{
			"data": id,
		},
	)
}
