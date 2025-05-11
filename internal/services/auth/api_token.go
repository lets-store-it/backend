package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
	"github.com/let-store-it/backend/internal/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *AuthService) GetOrgIdByApiToken(ctx context.Context, token string) (uuid.UUID, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetOrgIdByApiToken", func(ctx context.Context, span trace.Span) (uuid.UUID, error) {
		orgID, err := s.queries.GetOrgIdByApiToken(ctx, token)
		if err != nil {
			return uuid.Nil, services.MapDbErrorToService(err)
		}

		span.SetAttributes(attribute.String("org.id", uuid.UUID(orgID.Bytes).String()))
		return orgID.Bytes, nil
	})
}

func (s *AuthService) CreateApiToken(ctx context.Context, orgID uuid.UUID, name string) (*models.ApiToken, error) {
	return telemetry.WithTrace(ctx, s.tracer, "CreateApiToken", func(ctx context.Context, span trace.Span) (*models.ApiToken, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("token.name", name),
		)

		if name == "" {
			return nil, services.ErrValidationError
		}

		token, err := s.queries.CreateApiToken(ctx, sqlc.CreateApiTokenParams{
			OrgID: database.PgUUID(orgID),
			Token: uuid.New().String(),
			Name:  name,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create API token: %w", err)
		}

		span.SetAttributes(attribute.String("token.id", database.UUIDFromPgx(token.ID).String()))
		return toApiTokenModel(token), nil
	})
}

func (s *AuthService) GetApiTokens(ctx context.Context, orgID uuid.UUID) ([]*models.ApiToken, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetApiTokens", func(ctx context.Context, span trace.Span) ([]*models.ApiToken, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
		)

		tokens, err := s.queries.GetApiTokens(ctx, database.PgUUID(orgID))
		if err != nil {
			return nil, fmt.Errorf("failed to get API tokens: %w", err)
		}

		modelsTokens := make([]*models.ApiToken, len(tokens))
		for i, token := range tokens {
			modelsTokens[i] = toApiTokenModel(token)
		}

		span.SetAttributes(attribute.Int("tokens_count", len(modelsTokens)))
		return modelsTokens, nil
	})
}

func (s *AuthService) RevokeApiToken(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	return telemetry.WithVoidTrace(ctx, s.tracer, "RevokeApiToken", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("token.id", id.String()),
		)

		err := s.queries.RevokeApiToken(ctx, sqlc.RevokeApiTokenParams{
			OrgID: database.PgUUID(orgID),
			ID:    database.PgUUID(id),
		})
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		return nil
	})
}
