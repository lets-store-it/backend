package usecases

import (
	"errors"
)

var (
	ErrForbidden     = errors.New("forbidden")
	ErrNotAuthorized = errors.New("not authorized")
)
