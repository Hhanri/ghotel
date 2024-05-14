package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/db"
)

type AuthHandler struct {
	store *db.Store
}

func NewAuthHandler(store *db.Store) *AuthHandler {
	return &AuthHandler{
		store: store,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := h.store.User.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		return err
	}

	if ok := user.VerifyPassword(params.Password); !ok {
		return fmt.Errorf("invalid credentials")
	}

	fmt.Println("authenticated")

	return c.JSON(user)
}
