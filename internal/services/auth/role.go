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

var roleHierarchy = map[models.AccessLevel][]models.RoleName{
	models.AccessLevelWorker: {
		models.RoleOwner,
		models.RoleAdmin,
		models.RoleManager,
	},
	models.AccessLevelManager: {
		models.RoleOwner,
		models.RoleAdmin,
		models.RoleManager,
	},
	models.AccessLevelAdmin: {
		models.RoleOwner,
		models.RoleAdmin,
	},
	models.AccessLevelOwner: {
		models.RoleOwner,
	},
}

func (s *AuthService) SetUserRole(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, roleID models.RoleID) error {
	ctx, span := s.tracer.Start(ctx, "SetUserRole",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("user_id", userID.String()),
			attribute.Int("role.id", int(roleID)),
		),
	)
	defer span.End()

	if roleID < models.RoleOwnerID || roleID > models.RoleWorkerID {
		span.RecordError(services.ErrValidationError)
		span.SetStatus(codes.Error, "invalid role ID")
		return services.ErrValidationError
	}

	err := s.queries.AssignRoleToUser(ctx, sqlc.AssignRoleToUserParams{
		OrgID:  database.PgUUID(orgID),
		UserID: database.PgUUID(userID),
		RoleID: int32(roleID),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to set user role")
		return fmt.Errorf("failed to set user role: %w", err)
	}

	span.SetStatus(codes.Ok, "user role set")
	return nil
}

func (s *AuthService) RemoveUserRole(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "RemoveUserRole",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	err := s.queries.UnassignRoleFromUser(ctx, sqlc.UnassignRoleFromUserParams{
		OrgID:  database.PgUUID(orgID),
		UserID: database.PgUUID(userID),
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to unassign role from user")
		return fmt.Errorf("failed to unassign role from user: %w", err)
	}

	span.SetStatus(codes.Ok, "role removed from user")
	return nil
}

func (s *AuthService) CheckUserAccess(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, accessLevel models.AccessLevel) (bool, error) {
	ctx, span := s.tracer.Start(ctx, "CheckUserAccess",
		trace.WithAttributes(
			attribute.String("user.id", userID.String()),
			attribute.String("org.id", orgID.String()),
			attribute.String("access_level", string(accessLevel)),
		),
	)
	defer span.End()

	role, err := s.GetUserRole(ctx, userID, orgID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get user role")
		return false, fmt.Errorf("failed to get user role: %w", err)
	}

	allowedRoles, exists := roleHierarchy[accessLevel]
	if !exists {
		span.RecordError(fmt.Errorf("invalid access level: %s", accessLevel))
		span.SetStatus(codes.Error, "invalid access level")
		return false, fmt.Errorf("invalid access level: %s", accessLevel)
	}

	for _, allowedRole := range allowedRoles {
		if role.Name == allowedRole {
			return true, nil
		}
	}

	return false, nil
}

func (s *AuthService) GetUserRole(ctx context.Context, userID uuid.UUID, orgID uuid.UUID) (*models.Role, error) {
	ctx, span := s.tracer.Start(ctx, "GetUserRole",
		trace.WithAttributes(
			attribute.String("user.id", userID.String()),
			attribute.String("org.id", orgID.String()),
		),
	)
	defer span.End()

	role, err := s.queries.GetUserRoleInOrg(ctx, sqlc.GetUserRoleInOrgParams{
		OrgID:  database.PgUUID(orgID),
		UserID: database.PgUUID(userID),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get user role")
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	span.SetAttributes(
		attribute.Int("role.id", int(role.AppRole.ID)),
		attribute.String("role.name", role.AppRole.Name),
	)
	span.SetStatus(codes.Ok, "user role retrieved")
	return toRoleModel(role.AppRole), nil
}

func (s *AuthService) GetAvailableRoles(ctx context.Context) ([]*models.Role, error) {
	ctx, span := s.tracer.Start(ctx, "GetAvailableRoles")
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

func (s *AuthService) GetUserRoleInOrg(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) (*models.Role, error) {
	ctx, span := s.tracer.Start(ctx, "GetUserRoleInOrg",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	result, err := s.queries.GetUserRoleInOrg(ctx, sqlc.GetUserRoleInOrgParams{
		OrgID:  database.PgUUID(orgID),
		UserID: database.PgUUID(userID),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get user role")
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	span.SetStatus(codes.Ok, "user role retrieved successfully")
	return toRoleModel(result.AppRole), nil
}

func (s *AuthService) GetRoleById(ctx context.Context, id int) (*models.Role, error) {
	ctx, span := s.tracer.Start(ctx, "GetRoleById",
		trace.WithAttributes(
			attribute.Int("role.id", id),
		),
	)
	defer span.End()

	role, err := s.queries.GetRoleById(ctx, int32(id))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get role")
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	span.SetStatus(codes.Ok, "role retrieved successfully")
	return toRoleModel(role), nil
}
