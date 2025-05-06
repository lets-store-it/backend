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

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "GetUserByEmail")
	defer span.End()

	if email == "" {
		span.RecordError(ErrInvalidEmail)
		span.SetStatus(codes.Error, "invalid email")
		return nil, ErrInvalidEmail
	}

	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			span.RecordError(ErrUserNotFound)
			span.SetStatus(codes.Error, "user not found")
			return nil, ErrUserNotFound
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get user by email")
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	span.SetStatus(codes.Ok, "user found")
	return toUserModel(user), nil
}

func (s *AuthService) GetUserById(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "GetUserById",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
		),
	)
	defer span.End()

	if !utils.IsValidUUID(userID) {
		span.RecordError(ErrInvalidUserId)
		span.SetStatus(codes.Error, "invalid user ID")
		return nil, ErrInvalidUserId
	}

	user, err := s.queries.GetUserById(ctx, utils.PgUUID(userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			span.RecordError(ErrUserNotFound)
			span.SetStatus(codes.Error, "user not found")
			return nil, ErrUserNotFound
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get user by ID")
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	span.SetStatus(codes.Ok, "user found")
	return toUserModel(user), nil
}

func (s *AuthService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "CreateUser")
	defer span.End()

	if user == nil {
		span.RecordError(ErrInvalidUserId)
		span.SetStatus(codes.Error, "invalid user")
		return nil, ErrInvalidUserId
	}
	if user.Email == "" {
		span.RecordError(ErrInvalidEmail)
		span.SetStatus(codes.Error, "invalid email")
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
			span.RecordError(ErrDuplicateUser)
			span.SetStatus(codes.Error, "user already exists")
			return nil, ErrDuplicateUser
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	span.SetStatus(codes.Ok, "user created")
	return toUserModel(dbUser), nil
}
