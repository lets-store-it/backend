package services

import (
	"errors"

	"storeit/internal/database"
)

var (
	ErrNotFoundError    = errors.New("not found error")
	ErrDuplicationError = errors.New("duplication error")
	ErrValidationError  = errors.New("validation error")
)

// MapDbErrorToService maps database errors to corresponding service errors.
// If the error is not a known database error, it returns the original error.
func MapDbErrorToService(err error) error {
	if err == nil {
		return nil
	}

	if database.IsNotFound(err) {
		return ErrNotFoundError
	}

	if database.IsUniqueViolation(err) {
		return ErrDuplicationError
	}

	return err
}
