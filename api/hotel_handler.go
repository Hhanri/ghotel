package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/db"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore:  rs,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {

	hotels, err := h.hotelStore.List(c.Context(), struct{}{})
	if err != nil {
		return err
	}

	return c.JSON(hotels)

}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	rooms, err := h.roomStore.GetRoomsByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}
