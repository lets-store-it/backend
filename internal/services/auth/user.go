package auth

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func (s *AuthService) validateUserData(user *models.User) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}
	if strings.TrimSpace(user.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if !emailRegex.MatchString(user.Email) {
		return fmt.Errorf("invalid email format")
	}
	if strings.TrimSpace(user.FirstName) == "" {
		return fmt.Errorf("first name is required")
	}
	if strings.TrimSpace(user.LastName) == "" {
		return fmt.Errorf("last name is required")
	}
	return nil
}

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

func (s *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "GetUserByID",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
		),
	)
	defer span.End()

	if userID == uuid.Nil {
		span.RecordError(ErrInvalidUserId)
		span.SetStatus(codes.Error, "invalid user ID")
		return nil, ErrInvalidUserId
	}

	user, err := s.queries.GetUserById(ctx, database.PgUUID(userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			span.RecordError(ErrUserNotFound)
			span.SetStatus(codes.Error, "user not found")
			return nil, ErrUserNotFound
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	span.SetStatus(codes.Ok, "user retrieved successfully")
	return toUserModel(user), nil
}

func (s *AuthService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "CreateUser",
		trace.WithAttributes(
			attribute.String("email", user.Email),
			attribute.String("first_name", user.FirstName),
			attribute.String("last_name", user.LastName),
		),
	)
	defer span.End()

	if err := s.validateUserData(user); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	var middleName, yandexID string
	if user.MiddleName != nil {
		middleName = *user.MiddleName
	}
	if user.YandexID != nil {
		yandexID = *user.YandexID
	}

	createdUser, err := s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: database.PgText(middleName),
		YandexID:   database.PgText(yandexID),
	})
	if err != nil {
		if database.IsUniqueViolation(err) {
			span.RecordError(ErrDuplicateUser)
			span.SetStatus(codes.Error, "user already exists")
			return nil, ErrDuplicateUser
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	span.SetStatus(codes.Ok, "user created successfully")
	return toUserModel(createdUser), nil
}
