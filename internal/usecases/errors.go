package usecases

import (
	"errors"
)

type ErrDetailedValidationError struct {
	Message string
}

func (e *ErrDetailedValidationError) Error() string {
	return e.Message
}

func ErrDetailedValidationErrorWithMessage(message string) *ErrDetailedValidationError {
	return &ErrDetailedValidationError{
		Message: message,
	}
}

var (
	ErrNotAuthorized   = errors.New("not authorized")
	ErrValidationError = &ErrDetailedValidationError{Message: "validation error"}

	ErrUserIDMissing         = errors.New("user ID missing")
	ErrOrganizationIDMissing = errors.New("organization ID missing")
	ErrTvBoardIDMissing      = errors.New("TV board ID missing")

	ErrNotFound = errors.New("not found")
)
