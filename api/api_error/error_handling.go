package api_error

import (
	"github.com/gofiber/fiber/v2"
)

func FiberErrorHandler(c *fiber.Ctx, err error) error {
	errResp, ok := err.(ErrorResponse)
	if ok {
		return fiberErrorResponse(c, errResp)
	}
	return fiberErrorResponse(c, InternalErrorResponse)
}

func fiberErrorResponse(c *fiber.Ctx, err ErrorResponse) error {
	return c.Status(err.StatusCode).JSON(err)
}
