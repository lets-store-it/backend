package auth

import (
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
)

func toUserModel(user sqlc.AppUser) *models.User {
	return &models.User{
		ID:         user.ID.Bytes,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: database.PgTextPtrFromPgx(user.MiddleName),
		YandexID:   database.PgTextPtrFromPgx(user.YandexID),
	}
}

func toRoleModel(role sqlc.AppRole) *models.Role {
	return &models.Role{
		ID:          int(role.ID),
		Name:        models.RoleName(role.Name),
		DisplayName: role.DisplayName,
		Description: role.Description,
	}
}

func toApiTokenModel(token sqlc.AppApiToken) *models.ApiToken {
	return &models.ApiToken{
		ID:        token.ID.Bytes,
		OrgID:     token.OrgID.Bytes,
		Name:      token.Name,
		Token:     token.Token,
		CreatedAt: token.CreatedAt.Time,
		RevokedAt: database.PgTimePtrFromPgx(token.RevokedAt),
	}
}

func toEmployeeModel(employee interface{}) *models.Employee {
	var appUser sqlc.AppUser
	var appRole sqlc.AppRole

	switch e := employee.(type) {
	case sqlc.GetEmployeesRow:
		appUser = e.AppUser
		appRole = e.AppRole
	case sqlc.GetEmployeeRow:
		appUser = e.AppUser
		appRole = e.AppRole
	}

	return &models.Employee{
		UserID:     appUser.ID.Bytes,
		Email:      appUser.Email,
		FirstName:  appUser.FirstName,
		LastName:   appUser.LastName,
		MiddleName: database.PgTextPtrFromPgx(appUser.MiddleName),
		RoleID:     int(appRole.ID),
		Role:       toRoleModel(appRole),
	}
}

func toUserSessionModel(session sqlc.AppUserSession) *models.UserSession {
	return &models.UserSession{
		ID:     session.ID.Bytes,
		UserID: session.UserID.Bytes,
		Token:  session.Token,

		CreatedAt: session.CreatedAt.Time,
		ExpiresAt: database.PgTimePtrFromPgx(session.ExpiresAt),
		RevokedAt: database.PgTimePtrFromPgx(session.RevokedAt),
	}
}
