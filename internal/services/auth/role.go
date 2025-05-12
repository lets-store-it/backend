package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/common"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
	"github.com/let-store-it/backend/internal/telemetry"
	"go.opentelemetry.io/otel/attribute"
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
	return telemetry.WithVoidTrace(ctx, s.tracer, "SetUserRole", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("user_id", userID.String()),
			attribute.Int("role.id", int(roleID)),
		)

		if roleID < models.RoleOwnerID || roleID > models.RoleWorkerID {
			return common.ErrValidationError
		}

		err := s.queries.AssignRoleToUser(ctx, sqlc.AssignRoleToUserParams{
			OrgID:  database.PgUUID(orgID),
			UserID: database.PgUUID(userID),
			RoleID: int32(roleID),
		})
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		return nil
	})
}

func (s *AuthService) RemoveUserRole(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error {
	return telemetry.WithVoidTrace(ctx, s.tracer, "RemoveUserRole", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("user.id", userID.String()),
		)

		err := s.queries.UnassignRoleFromUser(ctx, sqlc.UnassignRoleFromUserParams{
			OrgID:  database.PgUUID(orgID),
			UserID: database.PgUUID(userID),
		})
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		return nil
	})
}

func (s *AuthService) CheckUserAccess(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, accessLevel models.AccessLevel) (bool, error) {
	return telemetry.WithTrace(ctx, s.tracer, "CheckUserAccess", func(ctx context.Context, span trace.Span) (bool, error) {
		span.SetAttributes(
			attribute.String("user.id", userID.String()),
			attribute.String("org.id", orgID.String()),
			attribute.String("access_level", string(accessLevel)),
		)

		role, err := s.GetUserRole(ctx, userID, orgID)
		if err != nil {
			return false, services.MapDbErrorToService(err)
		}

		allowedRoles, exists := roleHierarchy[accessLevel]
		if !exists {
			return false, fmt.Errorf("invalid access level: %s", accessLevel)
		}

		for _, allowedRole := range allowedRoles {
			if role.Name == allowedRole {
				return true, nil
			}
		}

		return false, nil
	})
}

func (s *AuthService) GetUserRole(ctx context.Context, userID uuid.UUID, orgID uuid.UUID) (*models.Role, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetUserRole", func(ctx context.Context, span trace.Span) (*models.Role, error) {
		span.SetAttributes(
			attribute.String("user.id", userID.String()),
			attribute.String("org.id", orgID.String()),
		)

		role, err := s.queries.GetUserRoleInOrg(ctx, sqlc.GetUserRoleInOrgParams{
			OrgID:  database.PgUUID(orgID),
			UserID: database.PgUUID(userID),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		span.SetAttributes(
			attribute.Int("role.id", int(role.AppRole.ID)),
			attribute.String("role.name", role.AppRole.Name),
		)
		return toRoleModel(role.AppRole), nil
	})
}

func (s *AuthService) GetAvailableRoles(ctx context.Context) ([]*models.Role, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetAvailableRoles", func(ctx context.Context, span trace.Span) ([]*models.Role, error) {
		roles, err := s.queries.GetRoles(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get roles: %w", err)
		}

		rolesModels := make([]*models.Role, len(roles))
		for i, role := range roles {
			rolesModels[i] = toRoleModel(role)
		}

		span.SetAttributes(attribute.Int("roles_count", len(rolesModels)))
		return rolesModels, nil
	})
}

func (s *AuthService) GetUserRoleInOrg(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) (*models.Role, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetUserRoleInOrg", func(ctx context.Context, span trace.Span) (*models.Role, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("user.id", userID.String()),
		)

		result, err := s.queries.GetUserRoleInOrg(ctx, sqlc.GetUserRoleInOrgParams{
			OrgID:  database.PgUUID(orgID),
			UserID: database.PgUUID(userID),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		return toRoleModel(result.AppRole), nil
	})
}

func (s *AuthService) GetRoleById(ctx context.Context, id int) (*models.Role, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetRoleById", func(ctx context.Context, span trace.Span) (*models.Role, error) {
		span.SetAttributes(
			attribute.Int("role.id", id),
		)

		role, err := s.queries.GetRoleById(ctx, int32(id))
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		return toRoleModel(role), nil
	})
}
