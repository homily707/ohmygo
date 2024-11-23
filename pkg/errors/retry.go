package errors

import "errors"

var ErrRetriable = errors.New("retryable error wrapper")

type RetriableError struct {
	Err error
}

func NewRetriableError(err error) RetriableError {
	return RetriableError{Err: err}
}

func (e RetriableError) Error() string {
	return e.Err.Error()
}

func (e RetriableError) Is(target error) bool {
	return target == ErrRetriable
}
