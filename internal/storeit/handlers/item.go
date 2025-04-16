package handlers

import (
	"context"

	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/storeit/models"
	"github.com/let-store-it/backend/internal/storeit/usecases"
)

func convertItemToDTO(item *models.Item) *api.ItemFull {
	variants := make([]api.ItemVariantEditable, 0)
	if item.Variants != nil {
		for _, v := range *item.Variants {
			var article api.OptNilString
			if v.Article != nil {
				article = api.NewOptNilString(*v.Article)
			}
			var ean13 api.OptNilInt
			if v.EAN13 != nil {
				ean13 = api.NewOptNilInt(int(*v.EAN13))
			}
			variants = append(variants, api.ItemVariantEditable{
				ID:      api.NewOptUUID(v.ID),
				Name:    v.Name,
				Article: article,
				Ean13:   ean13,
			})
		}
	}

	description := ""
	if item.Description != nil {
		description = *item.Description
	}

	return &api.ItemFull{
		ID:          api.NewOptUUID(item.ID),
		Name:        item.Name,
		Description: description,
		Variants:    variants,
	}
}

func convertDTOToItem(dto *api.CreateItemRequest) *models.Item {
	variants := make([]models.ItemVariant, 0, len(dto.Variants))
	for _, v := range dto.Variants {
		var article *string
		if v.Article.Set {
			article = &v.Article.Value
		}
		var ean13 *int64
		if v.Ean13.Set {
			val := int64(v.Ean13.Value)
			ean13 = &val
		}
		variants = append(variants, models.ItemVariant{
			ID:      v.ID.Value,
			Name:    v.Name,
			Article: article,
			EAN13:   ean13,
		})
	}

	return &models.Item{
		ID:          dto.ID.Value,
		Name:        dto.Name,
		Description: &dto.Description,
		Variants:    &variants,
	}
}

// CreateItem implements api.Handler.
func (h *RestApiImplementation) CreateItem(ctx context.Context, req *api.CreateItemRequest) (*api.CreateItemResponse, error) {
	item := convertDTOToItem(req)
	createdItem, err := h.itemUseCase.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return &api.CreateItemResponse{
		Data: *convertItemToDTO(createdItem),
	}, nil
}

// DeleteItem implements api.Handler.
func (h *RestApiImplementation) DeleteItem(ctx context.Context, params api.DeleteItemParams) error {
	orgID, err := usecases.GetOrganizationIDFromContext(ctx)
	if err != nil {
		return err
	}
	return h.itemUseCase.Delete(ctx, orgID, params.ID)
}

// GetItemById implements api.Handler.
func (h *RestApiImplementation) GetItemById(ctx context.Context, params api.GetItemByIdParams) (*api.GetItemByIdResponse, error) {
	item, err := h.itemUseCase.GetByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &api.GetItemByIdResponse{
		Data: *convertItemToDTO(item),
	}, nil
}

// GetItems implements api.Handler.
func (h *RestApiImplementation) GetItems(ctx context.Context) (*api.GetItemsResponse, error) {
	items, err := h.itemUseCase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	dtoItems := make([]api.Item, 0, len(items))
	for _, item := range items {
		description := ""
		if item.Description != nil {
			description = *item.Description
		}
		dtoItems = append(dtoItems, api.Item{
			ID:          api.NewOptUUID(item.ID),
			Name:        item.Name,
			Description: description,
		})
	}

	return &api.GetItemsResponse{
		Data: dtoItems,
	}, nil
}

// PatchItem implements api.Handler.
func (h *RestApiImplementation) PatchItem(ctx context.Context, req *api.PatchItemRequest, params api.PatchItemParams) (*api.PatchItemResponse, error) {
	updates := make(map[string]interface{})

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = &req.Description
	}

	// Handle variants updates
	if req.Variants != nil {
		variants := make([]interface{}, len(req.Variants))
		for i, v := range req.Variants {
			variant := make(map[string]interface{})
			if v.ID.Set {
				variant["id"] = v.ID.Value.String()
			}
			variant["name"] = v.Name
			if v.Article.Set {
				variant["article"] = &v.Article.Value
			}
			if v.Ean13.Set {
				variant["ean13"] = float64(v.Ean13.Value)
			}
			variants[i] = variant
		}
		updates["variants"] = variants
	}

	orgID, err := usecases.GetOrganizationIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	item, err := h.itemUseCase.Patch(ctx, orgID, params.ID, updates)
	if err != nil {
		return nil, err
	}

	return &api.PatchItemResponse{
		Data: *convertItemToDTO(item),
	}, nil
}

// UpdateItem implements api.Handler.
func (h *RestApiImplementation) UpdateItem(ctx context.Context, req *api.UpdateItemRequest, params api.UpdateItemParams) (*api.UpdateItemResponse, error) {
	variants := make([]models.ItemVariant, 0, len(req.Variants))
	for _, v := range req.Variants {
		var article *string
		if v.Article.Set {
			article = &v.Article.Value
		}
		var ean13 *int64
		if v.Ean13.Set {
			val := int64(v.Ean13.Value)
			ean13 = &val
		}
		variants = append(variants, models.ItemVariant{
			ID:      v.ID.Value,
			Name:    v.Name,
			Article: article,
			EAN13:   ean13,
		})
	}

	item := &models.Item{
		ID:          params.ID,
		Name:        req.Name,
		Description: &req.Description,
		Variants:    &variants,
	}

	orgID, err := usecases.GetOrganizationIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	updatedItem, err := h.itemUseCase.Update(ctx, orgID, item)
	if err != nil {
		return nil, err
	}

	return &api.UpdateItemResponse{
		Data: *convertItemToDTO(updatedItem),
	}, nil
}
