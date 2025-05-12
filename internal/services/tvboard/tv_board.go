package tvboard

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

type TvBoardService struct {
	queries *sqlc.Queries
	pgxpool *pgxpool.Pool
	tracer  trace.Tracer
}

type TvBoardServiceConfig struct {
	Queries *sqlc.Queries
	PGXPool *pgxpool.Pool
}

func New(cfg TvBoardServiceConfig) *TvBoardService {
	return &TvBoardService{
		queries: cfg.Queries,
		pgxpool: cfg.PGXPool,
		tracer:  otel.GetTracerProvider().Tracer("tvboard-service"),
	}
}

func toTvBoard(tvBoard sqlc.TvBoard) *models.TvBoard {
	return &models.TvBoard{
		ID:     database.UUIDFromPgx(tvBoard.ID),
		OrgID:  database.UUIDFromPgx(tvBoard.OrgID),
		UnitID: database.UUIDFromPgx(tvBoard.UnitID),
		Name:   tvBoard.Name,
		Token:  tvBoard.Token,
	}
}

func toTvBoards(tvBoards []sqlc.TvBoard) []*models.TvBoard {
	res := make([]*models.TvBoard, len(tvBoards))
	for i, tvBoard := range tvBoards {
		res[i] = toTvBoard(tvBoard)
	}
	return res
}

func (s *TvBoardService) CreateTvBoard(ctx context.Context, tvBoard *models.TvBoard) (*models.TvBoard, error) {
	return telemetry.WithTrace(ctx, s.tracer, "CreateTvBoard", func(ctx context.Context, span trace.Span) (*models.TvBoard, error) {
		span.SetAttributes(
			attribute.String("org.id", tvBoard.OrgID.String()),
			attribute.String("unit.id", tvBoard.UnitID.String()),
			attribute.String("tv_board.name", tvBoard.Name),
		)

		res, err := s.queries.CreateTvBoard(ctx, sqlc.CreateTvBoardParams{
			OrgID:  database.PgUUID(tvBoard.OrgID),
			UnitID: database.PgUUID(tvBoard.UnitID),
			Name:   tvBoard.Name,
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}
		return toTvBoard(res), nil
	})
}

func (s *TvBoardService) GetTvBoard(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.TvBoard, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetTvBoard", func(ctx context.Context, span trace.Span) (*models.TvBoard, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("tv_board.id", id.String()),
		)

		res, err := s.queries.GetTvBoardById(ctx, sqlc.GetTvBoardByIdParams{
			OrgID: database.PgUUID(orgID),
			ID:    database.PgUUID(id),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}
		return toTvBoard(res), nil
	})
}

func (s *TvBoardService) GetTvBoards(ctx context.Context, orgID uuid.UUID) ([]*models.TvBoard, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetTvBoards", func(ctx context.Context, span trace.Span) ([]*models.TvBoard, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
		)

		res, err := s.queries.GetTvBoards(ctx, database.PgUUID(orgID))
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}
		return toTvBoards(res), nil
	})
}

func (s *TvBoardService) DeleteTvBoard(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	return telemetry.WithVoidTrace(ctx, s.tracer, "DeleteTvBoard", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("tv_board.id", id.String()),
		)

		err := s.queries.DeleteTvBoard(ctx, sqlc.DeleteTvBoardParams{
			OrgID: database.PgUUID(orgID),
			ID:    database.PgUUID(id),
		})
		if err != nil {
			return services.MapDbErrorToService(err)
		}
		return nil
	})
}

func (s *TvBoardService) GetTvBoardByToken(ctx context.Context, token string) (*models.TvBoard, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetTvBoardByToken", func(ctx context.Context, span trace.Span) (*models.TvBoard, error) {
		span.SetAttributes(
			attribute.String("tv_board.token", token),
		)

		res, err := s.queries.GetTvBoardByToken(ctx, token)
		if err != nil {
			if database.IsNotFound(err) {
				return nil, services.ErrNotFoundError
			}
			return nil, services.MapDbErrorToService(err)
		}
		return toTvBoard(res), nil
	})
}
