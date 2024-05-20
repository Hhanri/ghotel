package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/api/api_error"
	"github.com/hhanri/ghotel/api/api_util"
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
	var params db.HotelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return api_error.BadRequestErrorResponse
	}

	hotels, err := h.store.Hotel.List(c.Context(), params)
	if err != nil {
		return api_error.InternalErrorResponse
	}

	return c.JSON(
		api_util.NewResourceResponse(
			hotels,
			int(params.Page),
		),
	)

}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	rooms, err := h.store.Room.GetRoomsByID(c.Context(), id)
	if err != nil {
		return api_error.InternalErrorResponse
	}
	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotelByID(c *fiber.Ctx) error {
	id := c.Params("id")
	rooms, err := h.store.Hotel.GetByID(c.Context(), id)
	if err != nil {
		return api_error.InternalErrorResponse
	}
	return c.JSON(rooms)
}
