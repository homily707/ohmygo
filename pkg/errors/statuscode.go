package errors

import (
	"fmt"
	"net/http"
)

var (
	NotFoundError = NewErrorWithCode(http.StatusNotFound, "not found")
)

type HttpError interface {
	StatusCode() int
}

type ErrorWithCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e ErrorWithCode) StatusCode() int {
	return e.Code
}

func (e ErrorWithCode) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func NewErrorWithCode(statusCode int, message string) *ErrorWithCode {
	return &ErrorWithCode{
		Code:    statusCode,
		Message: message,
	}
}
