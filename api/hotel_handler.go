package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/db"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {

	hotels, err := h.store.Hotel.List(c.Context(), struct{}{})
	if err != nil {
		return fiberInternalErrorResponse(c)
	}

	return c.JSON(hotels)

}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	rooms, err := h.store.Room.GetRoomsByID(c.Context(), id)
	if err != nil {
		return fiberInternalErrorResponse(c)
	}
	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotelByID(c *fiber.Ctx) error {
	id := c.Params("id")
	rooms, err := h.store.Hotel.GetByID(c.Context(), id)
	if err != nil {
		return fiberInternalErrorResponse(c)
	}
	return c.JSON(rooms)
}
