package errs

import (
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ERROR_BAD_REQUEST           = &AppError{Code: http.StatusBadRequest, Message: "Bad request"}
	ERROR_INTERNAL_SERVER_ERROR = &AppError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	ERROR_NOT_FOUND             = &AppError{Code: http.StatusNotFound, Message: "Not found"}
	ERROR_UNAUTHORIZED          = &AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	ERROR_USER_EXIST            = &AppError{Code: http.StatusConflict, Message: "User already exists"}
	ERROR_INVALID_CREDENTIALS   = &AppError{Code: http.StatusBadRequest, Message: "Invalid credentials"}
)

type ErrorResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *AppError   `json:"error,omitempty"`
}

func NewError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func NewErrorResponse(err error) ErrorResponse {
	appErr, ok := err.(*AppError)
	if !ok {
		appErr = ERROR_INTERNAL_SERVER_ERROR
	}
	return ErrorResponse{
		Success: false,
		Error:   appErr,
	}
}
