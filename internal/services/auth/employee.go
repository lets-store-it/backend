package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
	"github.com/let-store-it/backend/internal/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *AuthService) GetEmployees(ctx context.Context, orgID uuid.UUID) ([]*models.Employee, error) {
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
			employeesModels[i] = toEmployeeModel(employee)
		}

		span.SetAttributes(attribute.Int("response.count", len(employeesModels)))
		return employeesModels, nil
	})
}

func (s *AuthService) GetEmployee(ctx context.Context, orgID, userID uuid.UUID) (*models.Employee, error) {
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
			Role:       toRoleModel(role),
		}, nil
	})
}
