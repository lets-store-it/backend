package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// AccessLevel represents the different levels of access in the system
type AccessLevel string

const (
	AccessLevelWorker  AccessLevel = "worker"
	AccessLevelManager AccessLevel = "manager"
	AccessLevelAdmin   AccessLevel = "admin"
	AccessLevelOwner   AccessLevel = "owner"
)

// roleHierarchy defines which roles have access to each access level
var roleHierarchy = map[AccessLevel][]models.RoleName{
	AccessLevelWorker: {
		models.RoleOwner,
		models.RoleAdmin,
		models.RoleManager,
	},
	AccessLevelManager: {
		models.RoleOwner,
		models.RoleAdmin,
		models.RoleManager,
	},
	AccessLevelAdmin: {
		models.RoleOwner,
		models.RoleAdmin,
	},
	AccessLevelOwner: {
		models.RoleOwner,
	},
}

// CheckUserAccess checks if a user has the required access level
func (s *AuthService) CheckUserAccess(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, accessLevel AccessLevel) (bool, error) {
	ctx, span := s.tracer.Start(ctx, "CheckUserAccess",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("org_id", orgID.String()),
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

	err := s.queries.AssignRoleToUser(ctx, sqlc.AssignRoleToUserParams{
		OrgID:  database.PgUUID(orgID),
		UserID: database.PgUUID(userID),
		RoleID: int32(roleID),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to assign/update role for user")
		return fmt.Errorf("failed to assign/update role for user: %w", err)
	}

	span.SetStatus(codes.Ok, "role assigned/updated for user")
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
		attribute.Int("role_id", int(role.AppRole.ID)),
		attribute.String("role_name", role.AppRole.Name),
	)
	span.SetStatus(codes.Ok, "user role retrieved")
	return toRoleModel(role.AppRole), nil
}

func (s *AuthService) GetAvaiableRoles(ctx context.Context) ([]*models.Role, error) {
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

func (s *AuthService) GetRoles(ctx context.Context) ([]*models.Role, error) {
	ctx, span := s.tracer.Start(ctx, "GetRoles")
	defer span.End()

	roles, err := s.queries.GetRoles(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get roles")
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}

	result := make([]*models.Role, len(roles))
	for i, role := range roles {
		result[i] = toRoleModel(role)
	}

	span.SetStatus(codes.Ok, "roles retrieved successfully")
	return result, nil
}

func (s *AuthService) GetUserRoleInOrg(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) (*models.Role, error) {
	ctx, span := s.tracer.Start(ctx, "GetUserRoleInOrg",
		trace.WithAttributes(
			attribute.String("org_id", orgID.String()),
			attribute.String("user_id", userID.String()),
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

func (s *AuthService) GetRoleByID(ctx context.Context, id int) (*models.Role, error) {
	ctx, span := s.tracer.Start(ctx, "GetRoleByID",
		trace.WithAttributes(
			attribute.Int("role_id", id),
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
