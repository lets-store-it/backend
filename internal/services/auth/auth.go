package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidSession      = errors.New("invalid session")
	ErrSessionNotFound     = errors.New("session not found")
	ErrInvalidRole         = errors.New("invalid role")
	ErrInvalidApiToken     = errors.New("invalid API token")
	ErrApiTokenNotFound    = errors.New("API token not found")
	ErrDuplicateUser       = errors.New("user already exists")
	ErrInvalidOrganization = errors.New("invalid organization")
)

type AuthService struct {
	queries *database.Queries
	pgxPool *pgxpool.Pool
}

func New(queries *database.Queries, pgxPool *pgxpool.Pool) *AuthService {
	return &AuthService{queries: queries, pgxPool: pgxPool}
}

func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	slog.Debug("service:auth:GetUserByEmail", "email", email)
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

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
	}, nil
}

func (s *AuthService) CreateUserSession(ctx context.Context, userId uuid.UUID) (*models.Session, error) {
	slog.Debug("service:auth:CreateUserSession", "userId", userId)
	if userId == uuid.Nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	session, err := s.queries.CreateUserSession(ctx, database.CreateUserSessionParams{
		UserID: pgtype.UUID{Bytes: userId, Valid: true},
		Token:  uuid.New().String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user session: %w", err)
	}

	return &models.Session{
		ID:     session.ID.Bytes,
		UserID: session.UserID.Bytes,
		Secret: session.Token,
	}, nil
}

func (s *AuthService) GetUserBySessionSecret(ctx context.Context, sessionSecret string) (*models.User, error) {
	slog.Debug("service:auth:GetUserBySessionSecret", "sessionSecret", sessionSecret)
	if sessionSecret == "" {
		return nil, ErrInvalidSession
	}

	user, err := s.queries.GetUserBySessionSecret(ctx, sessionSecret)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get user by session secret: %w", err)
	}

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
	}, nil
}

func (s *AuthService) GetUserById(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	slog.Debug("service:auth:GetUserById", "userID", userID)
	if userID == uuid.Nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	user, err := s.queries.GetUserById(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

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
	}, nil
}

func (s *AuthService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	slog.Debug("service:auth:CreateUser", "user", user)
	if user == nil {
		return nil, fmt.Errorf("user cannot be nil")
	}
	if user.Email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	var middleName pgtype.Text
	if user.MiddleName != nil {
		middleName = pgtype.Text{String: *user.MiddleName, Valid: true}
	}
	var yandexID pgtype.Text
	if user.YandexID != nil {
		yandexID = pgtype.Text{String: *user.YandexID, Valid: true}
	}

	dbUser, err := s.queries.CreateUser(ctx, database.CreateUserParams{
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: middleName,
		YandexID:   yandexID,
	})
	if err != nil {
		// Check for unique constraint violation
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, ErrDuplicateUser
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &models.User{
		ID:         dbUser.ID.Bytes,
		Email:      dbUser.Email,
		FirstName:  dbUser.FirstName,
		LastName:   dbUser.LastName,
		MiddleName: user.MiddleName,
		YandexID:   user.YandexID,
	}, nil
}

func (s *AuthService) AssignRoleToUser(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, roleID int) error {
	slog.Debug("service:auth:AssignRoleToUser", "orgID", orgID, "userID", userID, "roleID", roleID)
	if orgID == uuid.Nil {
		return ErrInvalidOrganization
	}
	if userID == uuid.Nil {
		return fmt.Errorf("invalid user ID")
	}
	if roleID < 1 || roleID > 4 {
		return ErrInvalidRole
	}

	err := s.queries.AssignRoleToUser(ctx, database.AssignRoleToUserParams{
		OrgID:  pgtype.UUID{Bytes: orgID, Valid: true},
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		RoleID: int32(roleID),
	})
	if err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}

func (s *AuthService) GetUserRole(ctx context.Context, userID uuid.UUID, orgID uuid.UUID) (*models.Role, error) {
	slog.Debug("service:auth:GetUserRole", "userID", userID, "orgID", orgID)
	if userID == uuid.Nil {
		return nil, fmt.Errorf("invalid user ID")
	}
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}

	role, err := s.queries.GetUserRoleInOrg(ctx, database.GetUserRoleInOrgParams{
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		OrgID:  pgtype.UUID{Bytes: orgID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	return &models.Role{
		ID:          int(role.AppRole.ID),
		Name:        role.AppRole.Name,
		DisplayName: role.AppRole.DisplayName,
		Description: role.AppRole.Description,
	}, nil
}

func (s *AuthService) GetRoles(ctx context.Context) ([]*models.Role, error) {
	slog.Debug("service:auth:GetRoles")
	roles, err := s.queries.GetRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	rolesModels := make([]*models.Role, len(roles))
	for i, role := range roles {
		rolesModels[i] = &models.Role{
			ID:          int(role.ID),
			Name:        role.Name,
			DisplayName: role.DisplayName,
			Description: role.Description,
		}
	}
	return rolesModels, nil
}

func tokenToModel(token database.AppApiToken) *models.ApiToken {
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

func (s *AuthService) GetOrgIdByApiToken(ctx context.Context, token string) (uuid.UUID, error) {
	slog.Debug("service:auth:GetOrgIdByApiToken", "token", token)
	if token == "" {
		return uuid.Nil, ErrInvalidApiToken
	}

	orgID, err := s.queries.GetOrgIdByApiToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, ErrApiTokenNotFound
		}
		return uuid.Nil, fmt.Errorf("failed to get organization ID by API token: %w", err)
	}
	return orgID.Bytes, nil
}

func (s *AuthService) CreateApiToken(ctx context.Context, orgID uuid.UUID, name string) (*models.ApiToken, error) {
	slog.Debug("service:auth:CreateApiToken", "orgID", orgID, "name", name)
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if name == "" {
		return nil, fmt.Errorf("token name cannot be empty")
	}

	token, err := s.queries.CreateApiToken(ctx, database.CreateApiTokenParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		Token: uuid.New().String(),
		Name:  name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create API token: %w", err)
	}
	return tokenToModel(token), nil
}

func (s *AuthService) RevokeApiToken(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	slog.Debug("service:auth:RevokeApiToken", "orgID", orgID, "id", id)
	if orgID == uuid.Nil {
		return ErrInvalidOrganization
	}
	if id == uuid.Nil {
		return ErrInvalidApiToken
	}

	err := s.queries.RevokeApiToken(ctx, database.RevokeApiTokenParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrApiTokenNotFound
		}
		return fmt.Errorf("failed to revoke API token: %w", err)
	}
	return nil
}

func (s *AuthService) GetApiTokens(ctx context.Context, orgID uuid.UUID) ([]*models.ApiToken, error) {
	slog.Debug("service:auth:GetApiTokens", "orgID", orgID)
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}

	tokens, err := s.queries.GetApiTokens(ctx, pgtype.UUID{Bytes: orgID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get API tokens: %w", err)
	}
	modelsTokens := make([]*models.ApiToken, len(tokens))
	for i, token := range tokens {
		modelsTokens[i] = tokenToModel(token)
	}
	return modelsTokens, nil
}

func (s *AuthService) GetEmployees(ctx context.Context, orgID uuid.UUID) ([]*models.Employee, error) {
	slog.Debug("service:auth:GetEmployees", "orgID", orgID)
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}

	employees, err := s.queries.GetEmployees(ctx, pgtype.UUID{Bytes: orgID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}

	employeesModels := make([]*models.Employee, len(employees))
	for i, employee := range employees {
		var middleName *string
		if employee.AppUser.MiddleName.Valid {
			middleName = &employee.AppUser.MiddleName.String
		}
		employeesModels[i] = &models.Employee{
			UserID:     employee.AppUser.ID.Bytes,
			Email:      employee.AppUser.Email,
			FirstName:  employee.AppUser.FirstName,
			LastName:   employee.AppUser.LastName,
			MiddleName: middleName,
			RoleID:     int(employee.AppRole.ID),
			Role: &models.Role{
				ID:          int(employee.AppRole.ID),
				Name:        employee.AppRole.Name,
				DisplayName: employee.AppRole.DisplayName,
				Description: employee.AppRole.Description,
			},
		}
	}
	return employeesModels, nil
}

func (s *AuthService) DeleteEmployee(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error {
	slog.Debug("service:auth:DeleteEmployee", "orgID", orgID, "userID", userID)
	if orgID == uuid.Nil {
		return ErrInvalidOrganization
	}
	if userID == uuid.Nil {
		return fmt.Errorf("invalid user ID")
	}

	err := s.queries.UnassignRoleFromUser(ctx, database.UnassignRoleFromUserParams{
		OrgID:  pgtype.UUID{Bytes: orgID, Valid: true},
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to unassign role from user: %w", err)
	}
	return nil
}

func (s *AuthService) GetEmployee(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) (*models.Employee, error) {
	slog.Debug("service:auth:GetEmployee", "orgID", orgID, "userID", userID)
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if userID == uuid.Nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	employee, err := s.queries.GetEmployee(ctx, database.GetEmployeeParams{
		OrgID:  pgtype.UUID{Bytes: orgID, Valid: true},
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	var middleName *string
	if employee.AppUser.MiddleName.Valid {
		middleName = &employee.AppUser.MiddleName.String
	}

	return &models.Employee{
		UserID:     employee.AppUser.ID.Bytes,
		Email:      employee.AppUser.Email,
		FirstName:  employee.AppUser.FirstName,
		LastName:   employee.AppUser.LastName,
		MiddleName: middleName,
		RoleID:     int(employee.AppRole.ID),
		Role: &models.Role{
			ID:          int(employee.AppRole.ID),
			Name:        employee.AppRole.Name,
			DisplayName: employee.AppRole.DisplayName,
			Description: employee.AppRole.Description,
		},
	}, nil
}

func (s *AuthService) SetEmployeeRole(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, roleID int) error {
	slog.Debug("service:auth:SetEmployeeRole", "orgID", orgID, "userID", userID, "roleID", roleID)
	if orgID == uuid.Nil {
		return ErrInvalidOrganization
	}
	if userID == uuid.Nil {
		return fmt.Errorf("invalid user ID")
	}
	if roleID < 1 || roleID > 4 {
		return ErrInvalidRole
	}

	err := s.queries.AssignRoleToUser(ctx, database.AssignRoleToUserParams{
		OrgID:  pgtype.UUID{Bytes: orgID, Valid: true},
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		RoleID: int32(roleID),
	})
	if err != nil {
		return fmt.Errorf("failed to set employee role: %w", err)
	}
	return nil
}
