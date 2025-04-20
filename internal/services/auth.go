package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/storeit/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type AuthService struct {
	queries *database.Queries
	pgxPool *pgxpool.Pool
}

func NewAuthService(queries *database.Queries, pgxPool *pgxpool.Pool) *AuthService {
	return &AuthService{queries: queries, pgxPool: pgxPool}
}

func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	slog.Debug("service:auth:GetUserByEmail", "email", email)
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
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

	session, err := s.queries.CreateUserSession(ctx, database.CreateUserSessionParams{
		UserID: pgtype.UUID{Bytes: userId, Valid: true},
		Token:  uuid.New().String(),
	})
	if err != nil {
		return nil, err
	}

	return &models.Session{
		ID:     session.ID.Bytes,
		UserID: session.UserID.Bytes,
		Secret: session.Token,
	}, nil
}

func (s *AuthService) GetUserBySessionSecret(ctx context.Context, sessionSecret string) (*models.User, error) {
	slog.Debug("service:auth:GetUserBySessionSecret", "sessionSecret", sessionSecret)

	user, err := s.queries.GetUserBySessionSecret(ctx, sessionSecret)
	if err != nil {
		return nil, err
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

	user, err := s.queries.GetUserById(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
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
		return nil, err
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

type Role int

const (
	RoleOwner   Role = 1
	RoleAdmin   Role = 2
	RoleManager Role = 3
	RoleWorker  Role = 4
)

func (s *AuthService) AssignRoleToUser(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, role Role) error {
	slog.Debug("service:auth:AssignRoleToUser", "orgID", orgID, "userID", userID, "role", role)
	return s.queries.AssignRoleToUser(ctx, database.AssignRoleToUserParams{
		OrgID:  pgtype.UUID{Bytes: orgID, Valid: true},
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		RoleID: int32(role),
	})
}

func (s *AuthService) GetUserRoles(ctx context.Context, userID uuid.UUID, orgID uuid.UUID) (map[Role]struct{}, error) {
	slog.Debug("service:auth:GetUserRoles", "userID", userID, "orgID", orgID)
	dbRoles, err := s.queries.GetUserRolesInOrg(ctx, database.GetUserRolesInOrgParams{
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		OrgID:  pgtype.UUID{Bytes: orgID, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	roles := make(map[Role]struct{}, len(dbRoles))
	for _, role := range dbRoles {
		roles[Role(role.RoleID)] = struct{}{}
	}

	return roles, nil
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
	orgID, err := s.queries.GetOrgIdByApiToken(ctx, token)
	if err != nil {
		return uuid.Nil, err
	}
	return orgID.Bytes, nil
}

func (s *AuthService) CreateApiToken(ctx context.Context, orgID uuid.UUID, name string) (*models.ApiToken, error) {
	slog.Debug("service:auth:CreateApiToken", "orgID", orgID, "name", name)
	token, err := s.queries.CreateApiToken(ctx, database.CreateApiTokenParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		Token: uuid.New().String(),
	})
	if err != nil {
		return nil, err
	}
	return tokenToModel(token), nil
}

func (s *AuthService) RevokeApiToken(ctx context.Context, orgID uuid.UUID, token string) error {
	slog.Debug("service:auth:RevokeApiToken", "orgID", orgID, "token", token)
	return s.queries.RevokeApiToken(ctx, database.RevokeApiTokenParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		Token: token,
	})
}

func (s *AuthService) GetApiTokens(ctx context.Context, orgID uuid.UUID) ([]*models.ApiToken, error) {
	slog.Debug("service:auth:GetApiTokens", "orgID", orgID)
	tokens, err := s.queries.GetApiTokens(ctx, pgtype.UUID{Bytes: orgID, Valid: true})
	if err != nil {
		return nil, err
	}
	modelsTokens := make([]*models.ApiToken, len(tokens))
	for i, token := range tokens {
		modelsTokens[i] = tokenToModel(token)
	}
	return modelsTokens, nil
}
