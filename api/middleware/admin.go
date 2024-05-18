package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/api"
)

func AdminAuthentication(c *fiber.Ctx) error {
	user, err := api.GetAuth(c.Context())
	if err != nil {
		return api.FiberUnauthorizedErrorResponse(c)
	}
	if !user.IsAdmin {
		return api.FiberUnauthorizedErrorResponse(c)
	}
	return c.Next()
}
