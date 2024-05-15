package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/db"
	"github.com/hhanri/ghotel/types"
)

type UserHandler struct {
	store *db.Store
}

func NewUserHandler(store *db.Store) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.store.User.GetUserByID(c.Context(), id)
	if err != nil {
		return FiberInternalErrorResponse(c)
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.store.User.GetUsers(c.Context())
	if err != nil {
		return FiberInternalErrorResponse(c)
	}
	return c.JSON(users)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return FiberBadRequestErrorResponse(c)
	}

	if err := params.Validate(); len(err) > 0 {
		return c.JSON(err)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return FiberInternalErrorResponse(c)
	}

	newUser, err := h.store.User.InsertUser(c.Context(), user)
	if err != nil {
		return FiberInternalErrorResponse(c)
	}

	return c.JSON(newUser)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	err := h.store.User.DeleteUser(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(
		map[string]string{
			"data": userID,
		},
	)
}

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	var params types.UpdateUserParams
	if err := c.BodyParser(&params); err != nil {
		return FiberBadRequestErrorResponse(c)
	}

	err := h.store.User.UpdateUser(c.Context(), userID, params)
	if err != nil {
		return FiberInternalErrorResponse(c)
	}

	return c.JSON(
		map[string]string{
			"data": userID,
		},
	)

}
