package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/models"
)

// Storage Groups

func storageGroupToDTO(group *models.StorageGroup) api.StorageGroup {
	var parentID api.NilUUID
	PtrToApiNil(group.ParentID, &parentID)

	return api.StorageGroup{
		ID:       group.ID,
		ParentId: parentID,
		Name:     group.Name,
		Alias:    api.StorageAlias(group.Alias),
		UnitId:   group.UnitID,
	}
}

func (h *RestApiImplementation) CreateStorageGroup(ctx context.Context, req *api.StorageGroupBase) (api.CreateStorageGroupRes, error) {
	group, err := h.storageGroupUseCase.Create(ctx, &models.StorageGroup{
		UnitID:   req.UnitId,
		ParentID: ApiValueToPtr(req.ParentId),
		Name:     req.Name,
		Alias:    string(req.Alias),
	})
	if err != nil {
		return nil, err
	}

	groupDTO := storageGroupToDTO(group)
	return &api.CreateStorageGroupResponse{
		Data: groupDTO,
	}, nil
}

func (h *RestApiImplementation) GetStorageGroupById(ctx context.Context, params api.GetStorageGroupByIdParams) (api.GetStorageGroupByIdRes, error) {
	group, err := h.storageGroupUseCase.GetByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &api.GetStorageGroupByIdResponse{
		Data: storageGroupToDTO(group),
	}, nil
}

func (h *RestApiImplementation) GetStorageGroups(ctx context.Context) (api.GetStorageGroupsRes, error) {
	groups, err := h.storageGroupUseCase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]api.StorageGroup, 0, len(groups))
	for _, group := range groups {
		items = append(items, storageGroupToDTO(group))
	}

	return &api.GetStorageGroupsResponse{
		Data: items,
	}, nil
}

func (h *RestApiImplementation) DeleteStorageGroup(ctx context.Context, params api.DeleteStorageGroupParams) (api.DeleteStorageGroupRes, error) {
	err := h.storageGroupUseCase.Delete(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &api.DefaultNoContent{}, nil
}

func (h *RestApiImplementation) UpdateStorageGroup(ctx context.Context, req *api.StorageGroupBase, params api.UpdateStorageGroupParams) (api.UpdateStorageGroupRes, error) {
	group := &models.StorageGroup{
		ID:       params.ID,
		ParentID: ApiValueToPtr(req.ParentId),
		Name:     req.Name,
		Alias:    string(req.Alias),
		UnitID:   req.UnitId,
	}

	updatedGroup, err := h.storageGroupUseCase.Update(ctx, group)
	if err != nil {
		return nil, err
	}

	return &api.UpdateStorageGroupResponse{
		Data: storageGroupToDTO(updatedGroup),
	}, nil
}

// Cells group
func toCellsGroupDTO(group *models.CellsGroup) api.CellGroup {
	var storageGroupID api.NilUUID
	PtrToApiNil(group.StorageGroupID, &storageGroupID)

	return api.CellGroup{
		ID:             group.ID,
		Name:           group.Name,
		Alias:          api.StorageAlias(group.Alias),
		StorageGroupId: storageGroupID,
		UnitId:         group.UnitID,
	}
}

func (h *RestApiImplementation) GetCellsGroups(ctx context.Context) (api.GetCellsGroupsRes, error) {
	cellsGroups, err := h.storageGroupUseCase.GetCellsGroups(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]api.CellGroup, 0, len(cellsGroups))
	for _, group := range cellsGroups {
		items = append(items, toCellsGroupDTO(group))
	}

	return &api.GetCellsGroupsResponse{
		Data: items,
	}, nil
}

func (h *RestApiImplementation) CreateCellsGroup(ctx context.Context, req *api.CreateCellsGroupRequest) (api.CreateCellsGroupRes, error) {
	var storageGroupID *uuid.UUID
	if val, ok := req.StorageGroupId.Get(); ok {
		storageGroupID = &val
	}

	cellGroup, err := h.storageGroupUseCase.CreateCellsGroup(ctx, &models.CellsGroup{
		UnitID:         req.UnitId,
		StorageGroupID: storageGroupID,
		Name:           req.Name,
		Alias:          string(req.Alias),
	})
	if err != nil {
		return nil, err
	}

	return &api.CreateCellsGroupResponse{
		Data: toCellsGroupDTO(cellGroup),
	}, nil
}

// GetCellsGroupById implements api.Handler.
func (h *RestApiImplementation) GetCellsGroupById(ctx context.Context, params api.GetCellsGroupByIdParams) (api.GetCellsGroupByIdRes, error) {
	cellGroup, err := h.storageGroupUseCase.GetCellsGroupByID(ctx, params.GroupId)
	if err != nil {
		return nil, err
	}

	return &api.GetCellsGroupByIdResponse{
		Data: toCellsGroupDTO(cellGroup),
	}, nil
}

func (h *RestApiImplementation) UpdateCellsGroup(ctx context.Context, req *api.UpdateCellsGroupRequest, params api.UpdateCellsGroupParams) (api.UpdateCellsGroupRes, error) {
	groupInDB, err := h.storageGroupUseCase.GetCellsGroupByID(ctx, params.GroupId)
	if err != nil {
		return nil, err
	}

	model := &models.CellsGroup{
		ID:             params.GroupId,
		OrgID:          groupInDB.OrgID,
		UnitID:         req.UnitId,
		StorageGroupID: groupInDB.StorageGroupID,
		Name:           req.Name,
		Alias:          string(req.Alias),
	}
	cellGroup, err := h.storageGroupUseCase.UpdateCellsGroup(ctx, model)
	if err != nil {
		return nil, err
	}

	return &api.UpdateCellsGroupResponse{
		Data: toCellsGroupDTO(cellGroup),
	}, nil
}

func (h *RestApiImplementation) DeleteCellsGroup(ctx context.Context, params api.DeleteCellsGroupParams) (api.DeleteCellsGroupRes, error) {
	err := h.storageGroupUseCase.DeleteCellsGroup(ctx, params.GroupId)
	if err != nil {
		return nil, err
	}

	return &api.DeleteCellsGroupNoContent{}, nil
}

func cellToDTO(cell *models.Cell) api.Cell {
	return api.Cell{
		ID:       cell.ID,
		Alias:    cell.Alias,
		Row:      cell.Row,
		Level:    cell.Level,
		Position: cell.Position,
	}
}

// Cells
func (h *RestApiImplementation) CreateCell(ctx context.Context, req *api.CreateCellRequest, params api.CreateCellParams) (api.CreateCellRes, error) {
	cell := &models.Cell{
		Alias:        req.Alias,
		Row:          req.Row,
		Level:        req.Level,
		Position:     req.Position,
		CellsGroupID: params.GroupId,
	}
	cell, err := h.storageGroupUseCase.CreateCell(ctx, cell)
	if err != nil {
		return nil, err
	}

	return &api.CreateCellResponse{
		Data: cellToDTO(cell),
	}, nil
}

func (h *RestApiImplementation) GetCellById(ctx context.Context, params api.GetCellByIdParams) (api.GetCellByIdRes, error) {
	cell, err := h.storageGroupUseCase.GetCellByID(ctx, params.CellId)
	if err != nil {
		return nil, err
	}

	return &api.GetCellByIdResponse{
		Data: cellToDTO(cell),
	}, nil
}

// GetCells implements api.Handler.
func (h *RestApiImplementation) GetCells(ctx context.Context, params api.GetCellsParams) (api.GetCellsRes, error) {
	cells, err := h.storageGroupUseCase.GetCells(ctx, params.GroupId)
	if err != nil {
		return nil, err
	}

	items := make([]api.Cell, 0, len(cells))
	for _, cell := range cells {
		items = append(items, cellToDTO(cell))
	}

	return &api.GetCellsResponse{
		Data: items,
	}, nil
}

func (h *RestApiImplementation) UpdateCell(ctx context.Context, req *api.UpdateCellRequest, params api.UpdateCellParams) (api.UpdateCellRes, error) {
	cell := &models.Cell{
		ID:       params.CellId,
		Alias:    req.Alias,
		Row:      req.Row,
		Level:    req.Level,
		Position: req.Position,
	}

	updatedCell, err := h.storageGroupUseCase.UpdateCell(ctx, params.GroupId, cell)
	if err != nil {
		return nil, err
	}

	return &api.UpdateCellResponse{
		Data: api.Cell{
			ID:       updatedCell.ID,
			Alias:    updatedCell.Alias,
			Row:      updatedCell.Row,
			Level:    updatedCell.Level,
			Position: updatedCell.Position,
		},
	}, nil
}

func (h *RestApiImplementation) DeleteCell(ctx context.Context, params api.DeleteCellParams) (api.DeleteCellRes, error) {
	err := h.storageGroupUseCase.DeleteCell(ctx, params.CellId)
	if err != nil {
		return nil, err
	}

	return &api.DeleteCellNoContent{}, nil
}
