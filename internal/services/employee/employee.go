package employee

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
	"github.com/let-store-it/backend/internal/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type EmployeeService struct {
	queries *sqlc.Queries
	pgxPool *pgxpool.Pool
	tracer  trace.Tracer
}

type EmployeeServiceConfig struct {
	Queries *sqlc.Queries
	PGXPool *pgxpool.Pool
}

func New(cfg EmployeeServiceConfig) *EmployeeService {
	if cfg.Queries == nil || cfg.PGXPool == nil {
		panic("EmployeeServiceConfig is invalid")
	}

	return &EmployeeService{
		queries: cfg.Queries,
		pgxPool: cfg.PGXPool,
		tracer:  otel.GetTracerProvider().Tracer("employee-service"),
	}
}

func (s *EmployeeService) GetEmployees(ctx context.Context, orgID uuid.UUID) ([]*models.Employee, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetEmployees", func(ctx context.Context, span trace.Span) ([]*models.Employee, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
		)

		employees, err := s.queries.GetEmployees(ctx, database.PgUUID(orgID))
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		employeesModels := make([]*models.Employee, len(employees))
		for i, employee := range employees {
			employeesModels[i] = &models.Employee{
				UserID:     database.UUIDFromPgx(employee.AppUser.ID),
				Email:      employee.AppUser.Email,
				FirstName:  employee.AppUser.FirstName,
				LastName:   employee.AppUser.LastName,
				MiddleName: database.PgTextPtrFromPgx(employee.AppUser.MiddleName),
				RoleID:     int(employee.AppRole.ID),
				Role: &models.Role{
					ID:          int(employee.AppRole.ID),
					Name:        models.RoleName(employee.AppRole.Name),
					Description: employee.AppRole.Description,
					DisplayName: employee.AppRole.DisplayName,
				},
			}
		}

		span.SetAttributes(attribute.Int("response.count", len(employeesModels)))
		return employeesModels, nil
	})
}

func (s *EmployeeService) GetEmployee(ctx context.Context, orgID, userID uuid.UUID) (*models.Employee, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetUserAsEmployeeInOrg", func(ctx context.Context, span trace.Span) (*models.Employee, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("user.id", userID.String()),
		)

		employee, err := s.queries.GetEmployeeByUserId(ctx, sqlc.GetEmployeeByUserIdParams{
			OrgID:  database.PgUUID(orgID),
			UserID: database.PgUUID(userID),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		role, err := s.queries.GetRoleById(ctx, employee.RoleID)
		if err != nil {
			return nil, err
		}

		var middleName *string
		if employee.MiddleName.Valid {
			middleName = &employee.MiddleName.String
		}

		span.SetAttributes(
			attribute.Int("role_id", int(employee.RoleID)),
			attribute.String("role_name", role.Name),
		)

		return &models.Employee{
			UserID:     employee.UserID.Bytes,
			Email:      employee.Email,
			FirstName:  employee.FirstName,
			LastName:   employee.LastName,
			MiddleName: middleName,
			RoleID:     int(employee.RoleID),
			Role: &models.Role{
				ID:          int(role.ID),
				Name:        models.RoleName(role.Name),
				Description: role.Description,
				DisplayName: role.DisplayName,
			},
		}, nil
	})
}
