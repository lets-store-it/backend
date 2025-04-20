package handlers

import (
	"context"
	"errors"

	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/usecases"
)

func convertItemToFullDTO(item *models.Item, itemInstances *[]models.ItemInstance) *api.ItemFull {
	var description api.NilString
	if item.Description != nil {
		description.SetTo(*item.Description)
	}

	var variants []api.ItemVariant
	if item.Variants != nil {
		for _, variant := range *item.Variants {
			var article api.NilString
			if variant.Article != nil {
				article.SetTo(*variant.Article)
			} else {
				article.SetToNull()
			}
			var ean13 api.NilInt
			if variant.EAN13 != nil {
				ean13.SetTo(*variant.EAN13)
			} else {
				ean13.SetToNull()
			}
			variants = append(variants, api.ItemVariant{
				ID:      variant.ID,
				Name:    variant.Name,
				Article: article,
				Ean13:   ean13,
			})
		}
	}

	var instances []api.Instance
	if itemInstances != nil {
		for _, instance := range *itemInstances {
			var cellPath []api.CellForInstanceCellPathItem
			if instance.Cell != nil {
				for _, pathItem := range instance.Cell.Path {
					cellPath = append(cellPath, api.CellForInstanceCellPathItem{
						ID:         pathItem.ID,
						ObjectType: api.CellForInstanceCellPathItemObjectType(pathItem.ObjectType),
						Alias:      pathItem.Alias,
					})
				}
			}

			var article api.NilString
			if instance.Variant.Article != nil {
				article.SetTo(*instance.Variant.Article)
			} else {
				article.SetToNull()
			}

			var ean13 api.NilInt
			if instance.Variant.EAN13 != nil {
				ean13.SetTo(*instance.Variant.EAN13)
			} else {
				ean13.SetToNull()
			}

			instances = append(instances, api.ItemFullInstancesItem{
				ID:     instance.ID,
				Status: api.ItemFullInstancesItemStatus(instance.Status),
				Variant: api.ItemVariant{
					ID:      instance.VariantID,
					Name:    instance.Variant.Name,
					Article: article,
					Ean13:   ean13,
				},
				Cell: api.CellForInstance{
					ID:   instance.Cell.ID,
					Path: cellPath,
				},
			})
		}
	}
	return &api.ItemFull{
		ID:          item.ID,
		Name:        item.Name,
		Description: description,
		Variants:    variants,
		Instances:   instances,
	}
}

// CreateItem implements api.Handler.
func (h *RestApiImplementation) CreateItem(ctx context.Context, req *api.CreateItemRequest) (*api.CreateItemResponse, error) {
	var description *string
	if val, ok := req.Description.Get(); ok {
		description = &val
	}

	var variants []models.ItemVariant

	for _, variant := range req.Variants {
		var ean13 *int
		if val, ok := variant.Ean13.Get(); ok {
			ean13 = &val
		}
		var article *string
		if val, ok := variant.Article.Get(); ok {
			article = &val
		}

		variant := models.ItemVariant{
			Name:    variant.Name,
			Article: article,
			EAN13:   ean13,
		}
		variants = append(variants, variant)
	}

	item := &models.Item{
		Name:        req.Name,
		Description: description,
		Variants:    &variants,
	}

	createdItem, err := h.itemUseCase.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return &api.CreateItemResponse{
		Data: *convertItemToFullDTO(createdItem),
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
		Data: *convertItemToFullDTO(item),
	}, nil
}

// // GetItems implements api.Handler.
func (h *RestApiImplementation) GetItems(ctx context.Context) (*api.GetItemsResponse, error) {
	items, err := h.itemUseCase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	dtoItems := make([]api.ItemForList, 0, len(items))
	for _, item := range items {
		if item.Variants == nil {
			return nil, errors.New("variants are nil")
		}

		var variants []api.ItemVariant
		for _, variant := range *item.Variants {
			var article api.NilString
			if variant.Article != nil {
				article.SetTo(*variant.Article)
			} else {
				article.SetToNull()
			}
			var ean13 api.NilInt
			if variant.EAN13 != nil {
				ean13.SetTo(*variant.EAN13)
			} else {
				ean13.SetToNull()
			}

			variants = append(variants, api.ItemVariant{
				ID:      variant.ID,
				Name:    variant.Name,
				Article: article,
				Ean13:   ean13,
			})
		}

		var description api.NilString
		if item.Description != nil {
			description.SetTo(*item.Description)
		}
		dtoItems = append(dtoItems, api.ItemForList{
			ID:          item.ID,
			Name:        item.Name,
			Description: description,
			Variants:    variants,
		})
	}

	return &api.GetItemsResponse{
		Data: dtoItems,
	}, nil
}

// PatchItem implements api.Handler.
func (h *RestApiImplementation) PatchItem(ctx context.Context, req *api.PatchItemRequest, params api.PatchItemParams) (*api.PatchItemResponse, error) {
	updates := make(map[string]interface{})

	if val, ok := req.Name.Get(); ok {
		updates["name"] = val
	}
	if val, ok := req.Description.Get(); ok {
		updates["description"] = &val
	}

	// Handle variants updates
	if req.Variants != nil {
		variants := make([]interface{}, len(req.Variants))
		for i, v := range req.Variants {
			variant := make(map[string]interface{})
			variant["id"] = v.ID
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
		Data: *convertItemToFullDTO(item),
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
		var ean13 *int
		if v.Ean13.Set {
			val := v.Ean13.Value
			ean13 = &val
		}
		variants = append(variants, models.ItemVariant{
			ID:      v.ID,
			Name:    v.Name,
			Article: article,
			EAN13:   ean13,
		})
	}

	var description *string
	if val, ok := req.Description.Get(); ok {
		description = &val
	}
	item := &models.Item{
		ID:          params.ID,
		Name:        req.Name,
		Description: description,
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
		Data: *convertItemToFullDTO(updatedItem),
	}, nil
}
