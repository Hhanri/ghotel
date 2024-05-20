package api_error

import (
	"net/http"
)

type ErrorResponse struct {
	Message    string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

func (err ErrorResponse) Error() string {
	return err.Message
}

func NewErrorResponse(statusCode int, message string) ErrorResponse {
	return ErrorResponse{
		Message:    message,
		StatusCode: statusCode,
	}
}

var InternalErrorResponse = NewErrorResponse(
	http.StatusInternalServerError,
	"Internal Server Error",
)

var NotFoundErrorResponse = NewErrorResponse(
	http.StatusNotFound,
	"Not Found",
)

var InvalidCredentialsErrorResponse = NewErrorResponse(
	http.StatusUnauthorized,
	"Invalid Credentials",
)

var ExpiredTokenErrorResponse = NewErrorResponse(
	http.StatusUnauthorized,
	"Expired Token",
)

var UnauthorizedErrorResponse = NewErrorResponse(
	http.StatusUnauthorized,
	"Unauthorized",
)

var BadRequestErrorResponse = NewErrorResponse(
	http.StatusBadRequest,
	"Bad Request",
)
