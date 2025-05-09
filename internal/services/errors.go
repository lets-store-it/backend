package services

import "errors"

var (
	ErrNotFoundError    = errors.New("not found error")
	ErrDuplicationError = errors.New("duplication error")
	ErrValidationError  = errors.New("validation error")
	ErrNotAuthorized    = errors.New("not authorized")
)
