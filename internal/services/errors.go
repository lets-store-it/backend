package services

import (
	"errors"
	"fmt"

	"github.com/let-store-it/backend/internal/database"
)

var (
	ErrNotFoundError    = errors.New("not found error")
	ErrDuplicationError = errors.New("duplication error")
	ErrValidationError  = errors.New("validation error")
)

func MapDbErrorToService(err error) error {
	if err == nil {
		return fmt.Errorf("untranslated DB error: %w", err)
	}

	if database.IsNotFound(err) {
		return ErrNotFoundError
	}

	if database.IsUniqueViolation(err) {
		return ErrDuplicationError
	}

	return err
}
