package usecases

import (
	"errors"
)

var (
	ErrNotAuthorized         = errors.New("not authorized")
	ErrOrganizationIDMissing = errors.New("organization ID missing")
	ErrValidationError       = errors.New("validation error")
)
