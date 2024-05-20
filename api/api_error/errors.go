package api_error

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

var InternalErrorResponse = ErrorResponse{
	Error:      "Internal Server Error",
	StatusCode: http.StatusInternalServerError,
}

var NotFoundErrorResponse = ErrorResponse{
	Error:      "Not Found",
	StatusCode: http.StatusNotFound,
}

var InvalidCredentialsErrorResponse = ErrorResponse{
	Error:      "Invalid Credentials",
	StatusCode: http.StatusUnauthorized,
}

var ExpiredTokenErrorResponse = ErrorResponse{
	Error:      "Expired Token",
	StatusCode: http.StatusUnauthorized,
}

var UnauthorizedErrorResponse = ErrorResponse{
	Error:      "Unauthorized",
	StatusCode: http.StatusUnauthorized,
}

var BadRequestErrorResponse = ErrorResponse{
	Error:      "Bad Request",
	StatusCode: http.StatusBadRequest,
}

func FiberErrorResponse(c *fiber.Ctx, err ErrorResponse) error {
	return c.Status(err.StatusCode).JSON(err)
}

func FiberInternalErrorResponse(c *fiber.Ctx) error {
	return FiberErrorResponse(c, InternalErrorResponse)
}

func FiberNotFoundErrorResponse(c *fiber.Ctx) error {
	return FiberErrorResponse(c, NotFoundErrorResponse)
}

func FiberUnauthorizedErrorResponse(c *fiber.Ctx) error {
	return FiberErrorResponse(c, UnauthorizedErrorResponse)
}

func FiberInvalidCredentialsErrorResponse(c *fiber.Ctx) error {
	return FiberErrorResponse(c, InvalidCredentialsErrorResponse)
}

func FiberExpiredTokenErrorResponse(c *fiber.Ctx) error {
	return FiberErrorResponse(c, ExpiredTokenErrorResponse)
}

func FiberBadRequestErrorResponse(c *fiber.Ctx) error {
	return FiberErrorResponse(c, BadRequestErrorResponse)
}
