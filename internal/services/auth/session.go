package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func toSessionModel(session sqlc.AppUserSession) *models.UserSession {
	return &models.UserSession{
		ID:     session.ID.Bytes,
		UserID: session.UserID.Bytes,
		Token:  session.Token,

		CreatedAt: session.CreatedAt.Time,
		ExpiresAt: database.PgTimePtrFromPgx(session.ExpiresAt),
		RevokedAt: database.PgTimePtrFromPgx(session.RevokedAt),
	}
}

func (s *AuthService) CreateSession(ctx context.Context, userId uuid.UUID) (*models.UserSession, error) {
	ctx, span := s.tracer.Start(ctx, "CreateSession",
		trace.WithAttributes(
			attribute.String("user_id", userId.String()),
		),
	)
	defer span.End()

	session, err := s.queries.CreateUserSession(ctx, sqlc.CreateUserSessionParams{
		UserID: database.PgUUID(userId),
		Token:  uuid.New().String(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create session")
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	span.SetStatus(codes.Ok, "session created")
	return toSessionModel(session), nil
}

func (s *AuthService) GetSessionBySecret(ctx context.Context, sessionSecret string) (*models.UserSession, error) {
	ctx, span := s.tracer.Start(ctx, "GetSessionBySecret")
	defer span.End()

	session, err := s.queries.GetSessionBySecret(ctx, sessionSecret)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			span.RecordError(ErrSessionNotFound)
			span.SetStatus(codes.Error, "session not found")
			return nil, ErrSessionNotFound
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get user by session")
		return nil, fmt.Errorf("failed to get user by session secret: %w", err)
	}

	span.SetStatus(codes.Ok, "session found")
	return toSessionModel(session), nil
}

func (s *AuthService) InvalidateSession(ctx context.Context, sessionID uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "InvalidateSession")
	defer span.End()

	return s.queries.InvalidateSession(ctx, database.PgUUID(sessionID))
}
