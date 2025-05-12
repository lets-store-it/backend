package tvboard

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/common"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/organization"
	"github.com/let-store-it/backend/internal/services/tvboard"
	"github.com/let-store-it/backend/internal/usecases"
)

type TvBoardUseCase struct {
	tvBoardService      *tvboard.TvBoardService
	organizationService *organization.OrganizationService
	authService         *auth.AuthService
}

type TvBoardUseCaseConfig struct {
	TvBoardService      *tvboard.TvBoardService
	OrganizationService *organization.OrganizationService
	AuthService         *auth.AuthService
}

func New(config TvBoardUseCaseConfig) *TvBoardUseCase {
	return &TvBoardUseCase{
		tvBoardService:      config.TvBoardService,
		organizationService: config.OrganizationService,
		authService:         config.AuthService,
	}
}

func (uc *TvBoardUseCase) CreateTvBoard(ctx context.Context, tvBoard *models.TvBoard) (*models.TvBoard, error) {
	valRes, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return nil, err
	}

	if !valRes.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	tvBoard.OrgID = valRes.OrgID

	res, err := uc.tvBoardService.CreateTvBoard(ctx, tvBoard)
	if err != nil {
		return nil, err
	}

	unit, err := uc.organizationService.GetUnitByID(ctx, res.OrgID, res.UnitID)
	if err != nil {
		return nil, err
	}

	res.Unit = unit

	return res, nil
}

func (uc *TvBoardUseCase) GetTvBoard(ctx context.Context, id uuid.UUID) (*models.TvBoard, error) {
	valRes, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return nil, err
	}

	if !valRes.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	res, err := uc.tvBoardService.GetTvBoard(ctx, valRes.OrgID, id)
	if err != nil {
		return nil, err
	}

	unit, err := uc.organizationService.GetUnitByID(ctx, valRes.OrgID, res.UnitID)
	if err != nil {
		return nil, err
	}

	res.Unit = unit

	return res, nil
}

func (uc *TvBoardUseCase) GetTvBoards(ctx context.Context) ([]*models.TvBoard, error) {
	valRes, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return nil, err
	}

	if !valRes.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	res, err := uc.tvBoardService.GetTvBoards(ctx, valRes.OrgID)
	if err != nil {
		return nil, err
	}

	for _, tvBoard := range res {
		unit, err := uc.organizationService.GetUnitByID(ctx, valRes.OrgID, tvBoard.UnitID)
		if err != nil {
			return nil, err
		}
		tvBoard.Unit = unit
	}
	return res, nil
}

func (uc *TvBoardUseCase) DeleteTvBoard(ctx context.Context, id uuid.UUID) error {
	valRes, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return err
	}

	if !valRes.IsAllowed {
		return usecases.ErrNotAuthorized
	}

	return uc.tvBoardService.DeleteTvBoard(ctx, valRes.OrgID, id)
}

func (uc *TvBoardUseCase) GetTvBoardByToken(ctx context.Context, token string) (*models.TvBoard, error) {
	res, err := uc.tvBoardService.GetTvBoardByToken(ctx, token)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return res, nil
}
