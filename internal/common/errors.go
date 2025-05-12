package common

import "errors"

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
	ErrValidationError  = &ErrDetailedValidationError{Message: "validation error"}
	ErrNotFound         = errors.New("not found")
	ErrDuplicationError = errors.New("duplication error")
)

var (
	ErrUserIDMissing         = errors.New("user ID missing")
	ErrOrganizationIDMissing = errors.New("organization ID missing")
	ErrTvBoardIDMissing      = errors.New("TV board ID missing")
)
