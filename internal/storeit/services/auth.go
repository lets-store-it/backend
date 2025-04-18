package services

import (
	"context"
	"errors"
	"log/slog"

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
	slog.Info("service:GetUserByEmail", "email", email)
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
	slog.Info("service:CreateUserSession", "userId", userId)
	slog.Info("repository:CreateUserSession", "userId", userId)
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
	slog.Info("service:GetUserBySessionSecret", "sessionSecret", sessionSecret)
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
	slog.Info("service:GetCurrentUser", "userID", userID)
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
	slog.Info("service:CreateUser", "user", user)
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
	slog.Info("service:AssignRoleToUser", "orgID", orgID, "userID", userID, "role", role)
	slog.Info("service:AssignRoleToUser", "queries", s.queries)
	return s.queries.AssignRoleToUser(ctx, database.AssignRoleToUserParams{
		OrgID:  pgtype.UUID{Bytes: orgID, Valid: true},
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		RoleID: int32(role),
	})
}

func (s *AuthService) GetUserRoles(ctx context.Context, userID uuid.UUID, orgID uuid.UUID) (map[Role]struct{}, error) {
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
