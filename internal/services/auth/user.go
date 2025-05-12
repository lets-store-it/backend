package auth

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/common"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
	"github.com/let-store-it/backend/internal/telemetry"
	"github.com/let-store-it/backend/internal/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func (s *AuthService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return telemetry.WithTrace(ctx, s.tracer, "CreateUser", func(ctx context.Context, span trace.Span) (*models.User, error) {
		span.SetAttributes(
			attribute.String("user.email", user.Email),
			attribute.String("user.first_name", user.FirstName),
			attribute.String("user.last_name", user.LastName),
			attribute.String("user.middle_name", utils.SafeString(user.MiddleName)),
			attribute.String("user.yandex_id", utils.SafeString(user.YandexID)),
		)

		if err := s.validateUserData(user); err != nil {
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
			return nil, services.MapDbErrorToService(err)
		}

		return toUserModel(createdUser), nil
	})
}

func (s *AuthService) validateUserData(user *models.User) error {
	if user == nil {
		return fmt.Errorf("%w: user is nil", common.ErrValidationError)
	}
	if strings.TrimSpace(user.Email) == "" {
		return fmt.Errorf("%w: email is required", common.ErrValidationError)
	}
	if !emailRegex.MatchString(user.Email) {
		return fmt.Errorf("%w: invalid email format", common.ErrValidationError)
	}
	if strings.TrimSpace(user.FirstName) == "" {
		return fmt.Errorf("%w: first name is required", common.ErrValidationError)
	}
	if strings.TrimSpace(user.LastName) == "" {
		return fmt.Errorf("%w: last name is required", common.ErrValidationError)
	}
	return nil
}

func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetUserByEmail", func(ctx context.Context, span trace.Span) (*models.User, error) {
		span.SetAttributes(
			attribute.String("user.email", email),
		)

		user, err := s.queries.GetUserByEmail(ctx, email)
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		return toUserModel(user), nil
	})
}

func (s *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetUserByID", func(ctx context.Context, span trace.Span) (*models.User, error) {
		span.SetAttributes(
			attribute.String("user.id", userID.String()),
		)

		user, err := s.queries.GetUserById(ctx, database.PgUUID(userID))
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		return toUserModel(user), nil
	})
}
