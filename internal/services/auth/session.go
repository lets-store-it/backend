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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func (s *AuthService) CreateUserSession(ctx context.Context, userId uuid.UUID) (*models.Session, error) {
	ctx, span := s.tracer.Start(ctx, "CreateUserSession",
		trace.WithAttributes(
			attribute.String("user_id", userId.String()),
		),
	)
	defer span.End()

	session, err := s.queries.CreateUserSession(ctx, database.CreateUserSessionParams{
		UserID: utils.PgUUID(userId),
		Token:  uuid.New().String(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create user session")
		return nil, fmt.Errorf("failed to create user session: %w", err)
	}

	span.SetStatus(codes.Ok, "session created")
	return &models.Session{
		ID:     session.ID.Bytes,
		UserID: session.UserID.Bytes,
		Secret: session.Token,
	}, nil
}

func (s *AuthService) GetUserBySessionSecret(ctx context.Context, sessionSecret string) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "GetUserBySessionSecret")
	defer span.End()

	user, err := s.queries.GetUserBySessionSecret(ctx, sessionSecret)
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

	span.SetStatus(codes.Ok, "user found")
	return toUserModel(user), nil
}
