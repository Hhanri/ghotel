package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/db"
	"github.com/hhanri/ghotel/types"
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

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return fiberInternalErrorResponse(c)
	}

	user, err := h.store.User.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		return fiberNotFoundErrorResponse(c)
	}

	if ok := user.VerifyPassword(params.Password); !ok {
		return fiberInvalidCredentialsErrorResponse(c)
	}

	resp := AuthResponse{
		User:  user,
		Token: user.CreateJWT(),
	}
	return c.JSON(resp)
}
