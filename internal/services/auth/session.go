package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
	"github.com/let-store-it/backend/internal/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *AuthService) CreateSession(ctx context.Context, userId uuid.UUID) (*models.UserSession, error) {
	return telemetry.WithTrace(ctx, s.tracer, "CreateSession", func(ctx context.Context, span trace.Span) (*models.UserSession, error) {
		span.SetAttributes(
			attribute.String("user.id", userId.String()),
		)

		session, err := s.queries.CreateUserSession(ctx, sqlc.CreateUserSessionParams{
			UserID: database.PgUUID(userId),
			Token:  uuid.New().String(),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		return toSessionModel(session), nil
	})
}

func (s *AuthService) GetSessionBySecret(ctx context.Context, sessionSecret string) (*models.UserSession, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetSessionBySecret", func(ctx context.Context, span trace.Span) (*models.UserSession, error) {
		session, err := s.queries.GetSessionBySecret(ctx, sessionSecret)
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		return toSessionModel(session), nil
	})
}

func (s *AuthService) InvalidateSession(ctx context.Context, sessionID uuid.UUID) error {
	return telemetry.WithVoidTrace(ctx, s.tracer, "InvalidateSession", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("session.id", sessionID.String()),
		)

		err := s.queries.InvalidateSession(ctx, database.PgUUID(sessionID))
		if err != nil {
			return services.MapDbErrorToService(err)
		}
		return nil
	})
}
