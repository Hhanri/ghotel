package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/api/api_error"
	"github.com/hhanri/ghotel/api/api_util"
)

func AdminAuthentication(c *fiber.Ctx) error {
	user, err := api_util.GetAuth(c.Context())
	if err != nil {
		return api_error.FiberUnauthorizedErrorResponse(c)
	}
	if !user.IsAdmin {
		return api_error.FiberUnauthorizedErrorResponse(c)
	}
	return c.Next()
}
