package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/api"
	"github.com/hhanri/ghotel/types"
)

func AdminAuthentication(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return api.FiberBadRequestErrorResponse(c)
	}
	if !user.IsAdmin {
		return api.FiberUnauthorizedErrorResponse(c)
	}
	return c.Next()
}
