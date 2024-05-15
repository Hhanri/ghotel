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

var unauthorizedErrorResponse = ErrorResponse{
	Error:      "Unauthorized",
	StatusCode: http.StatusUnauthorized,
}

func fiberErrorResponse(c *fiber.Ctx, err ErrorResponse) error {
	return c.Status(err.StatusCode).JSON(err)
}

func fiberInternalErrorResponse(c *fiber.Ctx) error {
	return fiberErrorResponse(c, internalErrorResponse)
}

func fiberNotFoundErrorResponse(c *fiber.Ctx) error {
	return fiberErrorResponse(c, notFoundErrorResponse)
}

func fiberUnauthorizedErrorResponse(c *fiber.Ctx) error {
	return fiberErrorResponse(c, unauthorizedErrorResponse)
}

func fiberInvalidCredentialsErrorResponse(c *fiber.Ctx) error {
	return fiberErrorResponse(c, invalidCredentialsErrorResponse)
}
