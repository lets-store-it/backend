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
			attribute.String("org_id", orgID.String()),
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

	span.SetAttributes(attribute.Int("employees_count", len(employeesModels)))
	span.SetStatus(codes.Ok, "employees retrieved")
	return employeesModels, nil
}

func (s *AuthService) DeleteEmployee(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "DeleteEmployee",
		trace.WithAttributes(
			attribute.String("org_id", orgID.String()),
			attribute.String("user_id", userID.String()),
		),
	)
	defer span.End()

	if userID == uuid.Nil {
		span.RecordError(ErrInvalidUserId)
		span.SetStatus(codes.Error, "invalid user ID")
		return ErrInvalidUserId
	}

	err := s.queries.UnassignRoleFromUser(ctx, sqlc.UnassignRoleFromUserParams{
		OrgID:  database.PgUUID(orgID),
		UserID: database.PgUUID(userID),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to unassign role from user")
		return fmt.Errorf("failed to unassign role from user: %w", err)
	}

	span.SetStatus(codes.Ok, "employee deleted")
	return nil
}

func (s *AuthService) GetEmployee(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) (*models.Employee, error) {
	ctx, span := s.tracer.Start(ctx, "GetEmployee",
		trace.WithAttributes(
			attribute.String("org_id", orgID.String()),
			attribute.String("user_id", userID.String()),
		),
	)
	defer span.End()

	if userID == uuid.Nil {
		span.RecordError(ErrInvalidUserId)
		span.SetStatus(codes.Error, "invalid user ID")
		return nil, ErrInvalidUserId
	}

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

func (s *AuthService) SetEmployeeRole(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, roleID int) error {
	ctx, span := s.tracer.Start(ctx, "SetEmployeeRole",
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
		span.SetStatus(codes.Error, "failed to set employee role")
		return fmt.Errorf("failed to set employee role: %w", err)
	}

	span.SetStatus(codes.Ok, "employee role set")
	return nil
}

func (s *AuthService) GetEmployeeWithRole(ctx context.Context, orgID, userID uuid.UUID) (*models.Employee, error) {
	ctx, span := s.tracer.Start(ctx, "GetEmployeeWithRole",
		trace.WithAttributes(
			attribute.String("org_id", orgID.String()),
			attribute.String("user_id", userID.String()),
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
