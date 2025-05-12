package services

import (
	"fmt"

	"github.com/let-store-it/backend/internal/common"
	"github.com/let-store-it/backend/internal/database"
)

func MapDbErrorToService(err error) error {
	if err == nil {
		return fmt.Errorf("untranslated DB error: %w", err)
	}

	if database.IsNotFound(err) {
		return common.ErrNotFound
	}

	if database.IsUniqueViolation(err) {
		return common.ErrDuplicationError
	}

	return err
}
