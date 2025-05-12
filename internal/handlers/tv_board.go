package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/common"
	"github.com/let-store-it/backend/internal/models"
)

func toTvBoard(tvBoard *models.TvBoard) api.TvBoard {
	return api.TvBoard{
		ID:    tvBoard.ID,
		Name:  tvBoard.Name,
		Unit:  convertUnitToDTO(tvBoard.Unit),
		Token: tvBoard.Token,
	}
}

func toTvBoards(tvBoards []*models.TvBoard) []api.TvBoard {
	res := make([]api.TvBoard, len(tvBoards))
	for i, tvBoard := range tvBoards {
		res[i] = toTvBoard(tvBoard)
	}
	return res
}

func (h *RestApiImplementation) CreateTvBoard(ctx context.Context, req *api.CreateTvBoardRequest) (api.CreateTvBoardRes, error) {
	res, err := h.tvBoardUseCase.CreateTvBoard(ctx, &models.TvBoard{
		UnitID: req.UnitId,
		Name:   req.Name,
	})
	if err != nil {
		return nil, err
	}
	return &api.CreateTvBoardResponse{
		Data: toTvBoard(res),
	}, nil
}

func (h *RestApiImplementation) DeleteTvBoard(ctx context.Context, params api.DeleteTvBoardParams) (api.DeleteTvBoardRes, error) {
	err := h.tvBoardUseCase.DeleteTvBoard(ctx, params.ID)
	if err != nil {
		return nil, err
	}
	return &api.DeleteTvBoardNoContent{}, nil
}

func (h *RestApiImplementation) GetTvBoards(ctx context.Context) (api.GetTvBoardsRes, error) {
	res, err := h.tvBoardUseCase.GetTvBoards(ctx)
	if err != nil {
		return nil, err
	}
	return &api.GetTvBoardsResponse{
		Data: toTvBoards(res),
	}, nil
}

func (h *RestApiImplementation) GetTvBoardsData(ctx context.Context, params api.GetTvBoardsDataParams) (api.GetTvBoardsDataRes, error) {
	tvBoard, err := h.tvBoardUseCase.GetTvBoardByToken(ctx, params.TvToken)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return nil, h.NewUnauthorizedErrorWithMessage(ctx, "Invalid TV board token")
		}
		return nil, fmt.Errorf("failed to get TV board: %w", err)
	}
	
	ctx = context.WithValue(ctx, models.OrganizationIDContextKey, tvBoard.OrgID)
	ctx = context.WithValue(ctx, models.IsSystemUserContextKey, true)
	ctx = context.WithValue(ctx, models.TvBoardIDContextKey, tvBoard.ID)

	tvBoardID, err := common.GetTvBoardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := h.tvBoardUseCase.GetTvBoard(ctx, tvBoardID)
	if err != nil {
		return nil, err
	}

	taskRes, err := h.taskUseCase.GetTasks(ctx)
	if err != nil {
		return nil, err
	}

	return &api.GetTvBoardDataResponse{
		Data: api.GetTvBoardDataResponseData{
			TvBoard: toTvBoard(res),
			Tasks:   tasksToDto(taskRes),
		},
	}, nil
}
