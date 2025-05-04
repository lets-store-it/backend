package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/utils"
)

func (s *AuthService) GetEmployees(ctx context.Context, orgID uuid.UUID) ([]*models.Employee, error) {
	if !utils.IsValidUUID(orgID) {
		return nil, ErrInvalidOrganization
	}

	employees, err := s.queries.GetEmployees(ctx, utils.PgUUID(orgID))
	if err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}

	employeesModels := make([]*models.Employee, len(employees))
	for i, employee := range employees {
		employeesModels[i] = toEmployeeModel(employee)
	}
	return employeesModels, nil
}

func (s *AuthService) DeleteEmployee(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error {
	if !utils.IsValidUUID(orgID) {
		return ErrInvalidOrganization
	}
	if !utils.IsValidUUID(userID) {
		return ErrInvalidUserId
	}

	err := s.queries.UnassignRoleFromUser(ctx, database.UnassignRoleFromUserParams{
		OrgID:  utils.PgUUID(orgID),
		UserID: utils.PgUUID(userID),
	})
	if err != nil {
		return fmt.Errorf("failed to unassign role from user: %w", err)
	}
	return nil
}

func (s *AuthService) GetEmployee(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) (*models.Employee, error) {
	if !utils.IsValidUUID(orgID) {
		return nil, ErrInvalidOrganization
	}
	if !utils.IsValidUUID(userID) {
		return nil, ErrInvalidUserId
	}

	employee, err := s.queries.GetEmployee(ctx, database.GetEmployeeParams{
		OrgID:  utils.PgUUID(orgID),
		UserID: utils.PgUUID(userID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	return toEmployeeModel(employee), nil
}

func (s *AuthService) SetEmployeeRole(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, roleID int) error {
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
		return fmt.Errorf("failed to set employee role: %w", err)
	}
	return nil
}

func (s *AuthService) GetEmployeeWithRole(ctx context.Context, orgID, userID uuid.UUID) (*models.Employee, error) {
	employee, err := s.queries.GetEmployeeByUserId(ctx, database.GetEmployeeByUserIdParams{
		OrgID:  utils.PgUUID(orgID),
		UserID: utils.PgUUID(userID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	role, err := s.queries.GetRoleById(ctx, employee.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	var middleName *string
	if employee.MiddleName.Valid {
		middleName = &employee.MiddleName.String
	}

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
