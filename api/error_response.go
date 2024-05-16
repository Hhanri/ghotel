package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

var internalErrorResponse = ErrorResponse{
	Error:      "Internal Server Error",
	StatusCode: http.StatusInternalServerError,
}

var notFoundErrorResponse = ErrorResponse{
	Error:      "Not Found",
	StatusCode: http.StatusNotFound,
}

var invalidCredentialsErrorResponse = ErrorResponse{
	Error:      "Invalid Credentials",
	StatusCode: http.StatusUnauthorized,
}

var expiredTokenErrorResponse = ErrorResponse{
	Error:      "Expired Token",
	StatusCode: http.StatusUnauthorized,
}

var unauthorizedErrorResponse = ErrorResponse{
	Error:      "Unauthorized",
	StatusCode: http.StatusUnauthorized,
}

var badRequestErrorResponse = ErrorResponse{
	Error:      "Bad Request",
	StatusCode: http.StatusBadRequest,
}

func FiberErrorResponse(c *fiber.Ctx, err ErrorResponse) error {
	return c.Status(err.StatusCode).JSON(err)
}

func FiberInternalErrorResponse(c *fiber.Ctx) error {
	return FiberErrorResponse(c, internalErrorResponse)
}

func FiberNotFoundErrorResponse(c *fiber.Ctx) error {
	return FiberErrorResponse(c, notFoundErrorResponse)
}

func FiberUnauthorizedErrorResponse(c *fiber.Ctx) error {
	return FiberErrorResponse(c, unauthorizedErrorResponse)
}

func FiberInvalidCredentialsErrorResponse(c *fiber.Ctx) error {
	return FiberErrorResponse(c, invalidCredentialsErrorResponse)
}

func FiberExpiredTokenErrorResponse(c *fiber.Ctx) error {
	return FiberErrorResponse(c, expiredTokenErrorResponse)
}

func FiberBadRequestErrorResponse(c *fiber.Ctx) error {
	return FiberErrorResponse(c, badRequestErrorResponse)
}
