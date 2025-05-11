package auth

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
	"github.com/let-store-it/backend/internal/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func (s *AuthService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "CreateUser",
		trace.WithAttributes(
			attribute.String("user.email", user.Email),
			attribute.String("user.first_name", user.FirstName),
			attribute.String("user.last_name", user.LastName),
			attribute.String("user.middle_name", utils.SafeString(user.MiddleName)),
			attribute.String("user.yandex_id", utils.SafeString(user.YandexID)),
		),
	)
	defer span.End()

	if err := s.validateUserData(user); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	createdUser, err := s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: database.PgTextPtr(user.MiddleName),
		YandexID:   database.PgTextPtr(user.YandexID),
	})
	if err != nil {
		if database.IsUniqueViolation(err) {
			span.RecordError(services.ErrDuplicationError)
			span.SetStatus(codes.Error, "user already exists")
			return nil, services.ErrDuplicationError
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	span.SetStatus(codes.Ok, "user created successfully")
	return toUserModel(createdUser), nil
}

func (s *AuthService) validateUserData(user *models.User) error {
	if user == nil {
		return fmt.Errorf("%w: user is nil", services.ErrValidationError)
	}
	if strings.TrimSpace(user.Email) == "" {
		return fmt.Errorf("%w: email is required", services.ErrValidationError)
	}
	if !emailRegex.MatchString(user.Email) {
		return fmt.Errorf("%w: invalid email format", services.ErrValidationError)
	}
	if strings.TrimSpace(user.FirstName) == "" {
		return fmt.Errorf("%w: first name is required", services.ErrValidationError)
	}
	if strings.TrimSpace(user.LastName) == "" {
		return fmt.Errorf("%w: last name is required", services.ErrValidationError)
	}
	return nil
}

func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "GetUserByEmail",
		trace.WithAttributes(
			attribute.String("user.email", email),
		),
	)
	defer span.End()

	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if database.IsNotFound(err) {
			span.RecordError(services.ErrNotFoundError)
			span.SetStatus(codes.Error, "user not found")
			return nil, services.ErrNotFoundError
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
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	user, err := s.queries.GetUserById(ctx, database.PgUUID(userID))
	if err != nil {
		if database.IsNotFound(err) {
			span.RecordError(services.ErrNotFoundError)
			span.SetStatus(codes.Error, "user not found")
			return nil, services.ErrNotFoundError
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	span.SetStatus(codes.Ok, "user retrieved successfully")
	return toUserModel(user), nil
}
