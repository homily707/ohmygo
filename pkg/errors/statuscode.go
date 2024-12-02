package errors

import (
	"fmt"
	"net/http"
)

var (
	NotFoundError = NewErrorWithCode(http.StatusNotFound, "not found")
)

type ErrorWithCode struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e ErrorWithCode) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}

func NewErrorWithCode(statusCode int, message string) *ErrorWithCode {
	return &ErrorWithCode{
		StatusCode: statusCode,
		Message:    message,
	}
}
