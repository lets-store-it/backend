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
)

func (s *AuthService) GetOrgIdByApiToken(ctx context.Context, token string) (uuid.UUID, error) {
	if token == "" {
		return uuid.Nil, ErrInvalidApiToken
	}

	orgID, err := s.queries.GetOrgIdByApiToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, ErrApiTokenNotFound
		}
		return uuid.Nil, fmt.Errorf("failed to get organization ID by API token: %w", err)
	}
	return orgID.Bytes, nil
}

func (s *AuthService) CreateApiToken(ctx context.Context, orgID uuid.UUID, name string) (*models.ApiToken, error) {
	if !utils.IsValidUUID(orgID) {
		return nil, ErrInvalidOrganization
	}
	if name == "" {
		return nil, ErrValidationError
	}

	token, err := s.queries.CreateApiToken(ctx, database.CreateApiTokenParams{
		OrgID: utils.PgUUID(orgID),
		Token: uuid.New().String(),
		Name:  name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create API token: %w", err)
	}
	return toTokenModel(token), nil
}

func (s *AuthService) RevokeApiToken(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if !utils.IsValidUUID(orgID) {
		return ErrInvalidOrganization
	}
	if !utils.IsValidUUID(id) {
		return ErrInvalidApiToken
	}

	err := s.queries.RevokeApiToken(ctx, database.RevokeApiTokenParams{
		OrgID: utils.PgUUID(orgID),
		ID:    utils.PgUUID(id),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrApiTokenNotFound
		}
		return fmt.Errorf("failed to revoke API token: %w", err)
	}
	return nil
}

func (s *AuthService) GetApiTokens(ctx context.Context, orgID uuid.UUID) ([]*models.ApiToken, error) {
	if !utils.IsValidUUID(orgID) {
		return nil, ErrInvalidOrganization
	}

	tokens, err := s.queries.GetApiTokens(ctx, utils.PgUUID(orgID))
	if err != nil {
		return nil, fmt.Errorf("failed to get API tokens: %w", err)
	}

	modelsTokens := make([]*models.ApiToken, len(tokens))
	for i, token := range tokens {
		modelsTokens[i] = toTokenModel(token)
	}
	return modelsTokens, nil
}
