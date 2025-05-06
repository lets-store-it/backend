package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func (s *AuthService) AssignRoleToUser(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, roleID int) error {
	ctx, span := s.tracer.Start(ctx, "AssignRoleToUser",
		trace.WithAttributes(
			attribute.String("org_id", orgID.String()),
			attribute.String("user_id", userID.String()),
			attribute.Int("role_id", roleID),
		),
	)
	defer span.End()

	if roleID < 1 || roleID > 4 {
		span.RecordError(ErrInvalidRole)
		span.SetStatus(codes.Error, "invalid role ID")
		return ErrInvalidRole
	}

	err := s.queries.AssignRoleToUser(ctx, database.AssignRoleToUserParams{
		OrgID:  utils.PgUUID(orgID),
		UserID: utils.PgUUID(userID),
		RoleID: int32(roleID),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to assign role to user")
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	span.SetStatus(codes.Ok, "role assigned to user")
	return nil
}

func (s *AuthService) GetUserRole(ctx context.Context, userID uuid.UUID, orgID uuid.UUID) (*models.Role, error) {
	ctx, span := s.tracer.Start(ctx, "GetUserRole",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("org_id", orgID.String()),
		),
	)
	defer span.End()

	role, err := s.queries.GetUserRoleInOrg(ctx, database.GetUserRoleInOrgParams{
		UserID: utils.PgUUID(userID),
		OrgID:  utils.PgUUID(orgID),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get user role")
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	span.SetAttributes(
		attribute.Int("role_id", int(role.AppRole.ID)),
		attribute.String("role_name", role.AppRole.Name),
	)
	span.SetStatus(codes.Ok, "user role retrieved")
	return toRoleModel(role.AppRole), nil
}

func (s *AuthService) GetRoles(ctx context.Context) ([]*models.Role, error) {
	ctx, span := s.tracer.Start(ctx, "GetRoles")
	defer span.End()

	roles, err := s.queries.GetRoles(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get roles")
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}

	rolesModels := make([]*models.Role, len(roles))
	for i, role := range roles {
		rolesModels[i] = toRoleModel(role)
	}

	span.SetAttributes(attribute.Int("roles_count", len(rolesModels)))
	span.SetStatus(codes.Ok, "roles retrieved")
	return rolesModels, nil
}
