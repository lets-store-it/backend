package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/utils"
)

func (s *AuthService) AssignRoleToUser(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, roleID int) error {
	if !utils.IsValidUUID(orgID) {
		return ErrInvalidOrganization
	}
	if !utils.IsValidUUID(userID) {
		return ErrInvalidUserId
	}
	if roleID < 1 || roleID > 4 {
		return ErrInvalidRole
	}

	err := s.queries.AssignRoleToUser(ctx, database.AssignRoleToUserParams{
		OrgID:  utils.PgUUID(orgID),
		UserID: utils.PgUUID(userID),
		RoleID: int32(roleID),
	})
	if err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}

func (s *AuthService) GetUserRole(ctx context.Context, userID uuid.UUID, orgID uuid.UUID) (*models.Role, error) {
	if !utils.IsValidUUID(userID) {
		return nil, ErrInvalidUserId
	}
	if !utils.IsValidUUID(orgID) {
		return nil, ErrInvalidOrganization
	}

	role, err := s.queries.GetUserRoleInOrg(ctx, database.GetUserRoleInOrgParams{
		UserID: utils.PgUUID(userID),
		OrgID:  utils.PgUUID(orgID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	return toRoleModel(role.AppRole), nil
}

func (s *AuthService) GetRoles(ctx context.Context) ([]*models.Role, error) {
	roles, err := s.queries.GetRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}

	rolesModels := make([]*models.Role, len(roles))
	for i, role := range roles {
		rolesModels[i] = toRoleModel(role)
	}
	return rolesModels, nil
}
