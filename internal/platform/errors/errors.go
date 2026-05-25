package errors

import (
	"fmt"
	"net/http"
)

type Code string

const (
	CodeInvalidRequest  Code = "INVALID_REQUEST"
	CodeInvalidBase64   Code = "INVALID_BASE64"
	CodePayloadTooLarge Code = "PAYLOAD_TOO_LARGE"
	CodeInternal        Code = "INTERNAL_ERROR"
)

type AppError struct {
	Code    Code
	Message string
	Details map[string]interface{}
	Status  int
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func InvalidRequest(msg string) *AppError {
	return &AppError{
		Code:    CodeInvalidRequest,
		Message: msg,
		Status:  http.StatusBadRequest,
		Details: map[string]interface{}{},
	}
}

func InvalidBase64(msg string) *AppError {
	return &AppError{
		Code:    CodeInvalidBase64,
		Message: msg,
		Status:  http.StatusBadRequest,
		Details: map[string]interface{}{},
	}
}

func PayloadTooLarge(max int64) *AppError {
	return &AppError{
		Code:    CodePayloadTooLarge,
		Message: fmt.Sprintf("request body exceeds %d bytes", max),
		Status:  http.StatusRequestEntityTooLarge,
		Details: map[string]interface{}{"max_bytes": max},
	}
}

func Internal(msg string) *AppError {
	return &AppError{
		Code:    CodeInternal,
		Message: msg,
		Status:  http.StatusInternalServerError,
		Details: map[string]interface{}{},
	}
}

func StatusFor(err error) int {
	if ae, ok := err.(*AppError); ok {
		return ae.Status
	}
	return http.StatusInternalServerError
}
