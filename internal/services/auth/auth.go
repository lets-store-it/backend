package auth

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/generated/database"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidSession      = errors.New("invalid session")
	ErrSessionNotFound     = errors.New("session not found")
	ErrInvalidRole         = errors.New("invalid role")
	ErrInvalidEmail        = errors.New("invalid email")
	ErrInvalidUserId       = errors.New("invalid user")
	ErrInvalidApiToken     = errors.New("invalid API token")
	ErrApiTokenNotFound    = errors.New("API token not found")
	ErrDuplicateUser       = errors.New("user already exists")
	ErrInvalidOrganization = errors.New("invalid organization")
	ErrValidationError     = errors.New("validation error")
)

type AuthService struct {
	queries *database.Queries
	pgxPool *pgxpool.Pool
}

func New(queries *database.Queries, pgxPool *pgxpool.Pool) *AuthService {
	return &AuthService{queries: queries, pgxPool: pgxPool}
}
