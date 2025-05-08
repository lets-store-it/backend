package auth

import (
	"time"

	database "github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/models"
)

func toUserModel(user database.AppUser) *models.User {
	var middleName *string
	if user.MiddleName.Valid {
		middleName = &user.MiddleName.String
	}

	return &models.User{
		ID:         user.ID.Bytes,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: middleName,
	}
}

func toRoleModel(role database.AppRole) *models.Role {
	return &models.Role{
		ID:          int(role.ID),
		Name:        models.RoleName(role.Name),
		DisplayName: role.DisplayName,
		Description: role.Description,
	}
}

func toTokenModel(token database.AppApiToken) *models.ApiToken {
	var revokedAt *time.Time
	if token.RevokedAt.Valid {
		revokedAt = &token.RevokedAt.Time
	}
	return &models.ApiToken{
		ID:        token.ID.Bytes,
		OrgID:     token.OrgID.Bytes,
		Name:      token.Name,
		Token:     token.Token,
		CreatedAt: token.CreatedAt.Time,
		RevokedAt: revokedAt,
	}
}

func toEmployeeModel(employee interface{}) *models.Employee {
	var middleName *string
	var appUser database.AppUser
	var appRole database.AppRole

	switch e := employee.(type) {
	case database.GetEmployeesRow:
		appUser = e.AppUser
		appRole = e.AppRole
	case database.GetEmployeeRow:
		appUser = e.AppUser
		appRole = e.AppRole
	}

	if appUser.MiddleName.Valid {
		middleName = &appUser.MiddleName.String
	}
	return &models.Employee{
		UserID:     appUser.ID.Bytes,
		Email:      appUser.Email,
		FirstName:  appUser.FirstName,
		LastName:   appUser.LastName,
		MiddleName: middleName,
		RoleID:     int(appRole.ID),
		Role:       toRoleModel(appRole),
	}
}
