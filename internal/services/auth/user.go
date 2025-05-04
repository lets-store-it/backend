package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/utils"
)

func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}

	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return toUserModel(user), nil
}

func (s *AuthService) GetUserById(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	if !utils.IsValidUUID(userID) {
		return nil, ErrInvalidUserId
	}

	user, err := s.queries.GetUserById(ctx, utils.PgUUID(userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return toUserModel(user), nil
}

func (s *AuthService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if user == nil {
		return nil, ErrInvalidUserId
	}
	if user.Email == "" {
		return nil, ErrInvalidEmail
	}

	dbUser, err := s.queries.CreateUser(ctx, database.CreateUserParams{
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: utils.PgTextPtr(user.MiddleName),
		YandexID:   utils.PgTextPtr(user.YandexID),
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, ErrDuplicateUser
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return toUserModel(dbUser), nil
}
