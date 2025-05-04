package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/utils"
)

func (s *AuthService) CreateUserSession(ctx context.Context, userId uuid.UUID) (*models.Session, error) {
	if !utils.IsValidUUID(userId) {
		return nil, ErrInvalidUserId
	}

	session, err := s.queries.CreateUserSession(ctx, database.CreateUserSessionParams{
		UserID: utils.PgUUID(userId),
		Token:  uuid.New().String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user session: %w", err)
	}

	return &models.Session{
		ID:     session.ID.Bytes,
		UserID: session.UserID.Bytes,
		Secret: session.Token,
	}, nil
}

func (s *AuthService) GetUserBySessionSecret(ctx context.Context, sessionSecret string) (*models.User, error) {
	if sessionSecret == "" {
		return nil, ErrInvalidSession
	}

	user, err := s.queries.GetUserBySessionSecret(ctx, sessionSecret)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get user by session secret: %w", err)
	}

	return toUserModel(user), nil
}
