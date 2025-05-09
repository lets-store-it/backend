package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func (s *AuthService) GetOrgIdByApiToken(ctx context.Context, token string) (uuid.UUID, error) {
	ctx, span := s.tracer.Start(ctx, "GetOrgIdByApiToken")
	defer span.End()

	orgID, err := s.queries.GetOrgIdByApiToken(ctx, token)
	if err != nil {
		if database.IsNotFound(err) {
			span.RecordError(services.ErrNotFoundError)
			span.SetStatus(codes.Error, "API token not found")
			return uuid.Nil, services.ErrNotFoundError
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get organization ID by API token")
		return uuid.Nil, fmt.Errorf("failed to get organization ID by API token: %w", err)
	}

	span.SetAttributes(attribute.String("org.id", uuid.UUID(orgID.Bytes).String()))
	span.SetStatus(codes.Ok, "organization found")
	return orgID.Bytes, nil
}

func (s *AuthService) CreateApiToken(ctx context.Context, orgID uuid.UUID, name string) (*models.ApiToken, error) {
	ctx, span := s.tracer.Start(ctx, "CreateApiToken",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("token.name", name),
		),
	)
	defer span.End()

	if name == "" {
		span.RecordError(services.ErrValidationError)
		span.SetStatus(codes.Error, "token name is required")
		return nil, services.ErrValidationError
	}

	token, err := s.queries.CreateApiToken(ctx, sqlc.CreateApiTokenParams{
		OrgID: database.PgUUID(orgID),
		Token: uuid.New().String(),
		Name:  name,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create API token")
		return nil, fmt.Errorf("failed to create API token: %w", err)
	}

	span.SetAttributes(attribute.String("token.id", database.UuidFromPgx(token.ID).String()))
	span.SetStatus(codes.Ok, "API token created")
	return toTokenModel(token), nil
}

func (s *AuthService) RevokeApiToken(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "RevokeApiToken",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("token.id", id.String()),
		),
	)
	defer span.End()

	err := s.queries.RevokeApiToken(ctx, sqlc.RevokeApiTokenParams{
		OrgID: database.PgUUID(orgID),
		ID:    database.PgUUID(id),
	})
	if err != nil {
		if database.IsNotFound(err) {
			span.RecordError(services.ErrNotFoundError)
			span.SetStatus(codes.Error, "API token not found")
			return services.ErrNotFoundError
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to revoke API token")
		return fmt.Errorf("failed to revoke API token: %w", err)
	}

	span.SetStatus(codes.Ok, "API token revoked")
	return nil
}

func (s *AuthService) GetApiTokens(ctx context.Context, orgID uuid.UUID) ([]*models.ApiToken, error) {
	ctx, span := s.tracer.Start(ctx, "GetApiTokens",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
		),
	)
	defer span.End()

	tokens, err := s.queries.GetApiTokens(ctx, database.PgUUID(orgID))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get API tokens")
		return nil, fmt.Errorf("failed to get API tokens: %w", err)
	}

	modelsTokens := make([]*models.ApiToken, len(tokens))
	for i, token := range tokens {
		modelsTokens[i] = toTokenModel(token)
	}

	span.SetAttributes(attribute.Int("tokens_count", len(modelsTokens)))
	span.SetStatus(codes.Ok, "API tokens retrieved")
	return modelsTokens, nil
}
