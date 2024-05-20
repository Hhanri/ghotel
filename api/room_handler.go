package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/api/api_error"
	"github.com/hhanri/ghotel/api/api_util"
	"github.com/hhanri/ghotel/db"
	"github.com/hhanri/ghotel/types"
	"go.mongodb.org/mongo-driver/bson"
)

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

type BookRoomParams struct {
	FromDate  time.Time `json:"fromDate"`
	UntilDate time.Time `json:"untilDate"`
	NumPeople int       `json:"numPeople"`
}

func (h *RoomHandler) HandleGetAllRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), struct{}{})
	if err != nil {
		return api_error.FiberInternalErrorResponse(c)
	}
	return c.JSON(rooms)
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || p.UntilDate.Before(now) {
		return fmt.Errorf("Can not book a room from before now")
	}
	if p.UntilDate.Before(p.FromDate) {
		return fmt.Errorf("Until must be after From")
	}
	return nil
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	roomId := c.Params("id")
	user, err := api_util.GetAuth(c.Context())
	if err != nil {
		return api_error.FiberUnauthorizedErrorResponse(c)
	}

	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return api_error.FiberBadRequestErrorResponse(c)
	}

	if err := params.validate(); err != nil {
		return api_error.FiberErrorResponse(
			c,
			api_error.ErrorResponse{
				Error:      err.Error(),
				StatusCode: http.StatusBadRequest,
			},
		)
	}
	exists, err := h.roomExists(c.Context(), roomId)
	if err != nil || !exists {
		return api_error.FiberBadRequestErrorResponse(c)
	}

	ok, err := h.isRoomAvailable(c.Context(), roomId, params)
	if err != nil {
		return api_error.FiberInternalErrorResponse(c)
	}
	if !ok {
		return api_error.FiberErrorResponse(
			c,
			api_error.ErrorResponse{
				Error:      "Room already booked",
				StatusCode: http.StatusNotAcceptable,
			},
		)
	}

	booking := types.Booking{
		UserID:    user.ID,
		RoomdID:   roomId,
		FromDate:  params.FromDate,
		UntilDate: params.UntilDate,
		NumPeople: params.NumPeople,
	}

	inserted, err := h.store.Booking.Insert(c.Context(), &booking)
	if err != nil {
		return api_error.FiberBadRequestErrorResponse(c)
	}

	return c.JSON(inserted)
}

func (h *RoomHandler) roomExists(ctx context.Context, roomId string) (bool, error) {
	room, err := h.store.Room.GetRoomByID(ctx, roomId)
	if err != nil {
		return false, err
	}
	return room != nil, nil
}

func (h *RoomHandler) isRoomAvailable(
	ctx context.Context,
	roomId string,
	params BookRoomParams,
) (bool, error) {
	filter := bson.M{
		"roomId": roomId,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"untilDate": bson.M{
			"$lte": params.UntilDate,
		},
		"canceled": false,
	}

	bookings, err := h.store.Booking.GetAll(
		ctx,
		filter,
	)

	if err != nil {
		return false, err
	}

	return len(bookings) < 1, nil
}
