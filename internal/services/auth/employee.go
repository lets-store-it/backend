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

func (s *AuthService) GetEmployees(ctx context.Context, orgID uuid.UUID) ([]*models.Employee, error) {
	ctx, span := s.tracer.Start(ctx, "GetEmployees",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
		),
	)
	defer span.End()

	employees, err := s.queries.GetEmployees(ctx, database.PgUUID(orgID))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get employees")
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}

	employeesModels := make([]*models.Employee, len(employees))
	for i, employee := range employees {
		employeesModels[i] = toEmployeeModel(employee)
	}

	span.SetAttributes(attribute.Int("response.count", len(employeesModels)))
	span.SetStatus(codes.Ok, "employees retrieved")
	return employeesModels, nil
}

func (s *AuthService) GetEmployee(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) (*models.Employee, error) {
	ctx, span := s.tracer.Start(ctx, "GetEmployee",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	employee, err := s.queries.GetEmployee(ctx, sqlc.GetEmployeeParams{
		OrgID:  database.PgUUID(orgID),
		UserID: database.PgUUID(userID),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get employee")
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	span.SetStatus(codes.Ok, "employee found")
	return toEmployeeModel(employee), nil
}

func (s *AuthService) GetUserAsEmployeeInOrg(ctx context.Context, orgID, userID uuid.UUID) (*models.Employee, error) {
	ctx, span := s.tracer.Start(ctx, "GetUserAsEmployeeInOrg",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	employee, err := s.queries.GetEmployeeByUserId(ctx, sqlc.GetEmployeeByUserIdParams{
		OrgID:  database.PgUUID(orgID),
		UserID: database.PgUUID(userID),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get employee")
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	role, err := s.queries.GetRoleById(ctx, employee.RoleID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get role")
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	var middleName *string
	if employee.MiddleName.Valid {
		middleName = &employee.MiddleName.String
	}

	span.SetAttributes(
		attribute.Int("role_id", int(employee.RoleID)),
		attribute.String("role_name", role.Name),
	)
	span.SetStatus(codes.Ok, "employee with role found")

	return &models.Employee{
		UserID:     employee.UserID.Bytes,
		Email:      employee.Email,
		FirstName:  employee.FirstName,
		LastName:   employee.LastName,
		MiddleName: middleName,
		RoleID:     int(employee.RoleID),
		Role:       toRoleModel(role),
	}, nil
}
