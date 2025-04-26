package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-faster/jx"
	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/models"
)

// GetAuditLogs implements api.Handler.

func toObjectTypeDTO(objectType *models.ObjectType) api.AuditLogTargetObjectType {
	return api.AuditLogTargetObjectType{
		ID:    int(objectType.ID),
		Name:  objectType.Name,
		Group: objectType.Group,
	}
}

func convertJsonRawMessageToJxMap(rawMsg json.RawMessage) (map[string]jx.Raw, error) {
	// Handle nil or empty JSON
	if len(rawMsg) == 0 || string(rawMsg) == "null" {
		return make(map[string]jx.Raw), nil
	}

	// Create a decoder from the raw message
	d := jx.DecodeBytes(rawMsg)

	// Create the target map
	result := make(map[string]jx.Raw)

	// Decode the object
	if err := d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		// For each field in the object
		v, err := d.RawAppend(nil)
		if err != nil {
			return err
		}
		result[string(k)] = jx.Raw(v)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("decode object: %w", err)
	}

	return result, nil
}

func toAuditLog(objectChange *models.ObjectChange) (*api.AuditLog, error) {
	var prechangeState api.NilAuditLogPrechangeState
	var postchangeState api.NilAuditLogPostchangeState

	// Handle prechange state
	if len(objectChange.PrechangeState) > 0 && string(objectChange.PrechangeState) != "null" {
		preStateMap, err := convertJsonRawMessageToJxMap(objectChange.PrechangeState)
		if err != nil {
			return nil, fmt.Errorf("convert prechange state: %w", err)
		}
		prechangeState.Value = api.AuditLogPrechangeState(preStateMap)
		prechangeState.Null = false
	} else {
		prechangeState.Null = true
	}

	// Handle postchange state
	if len(objectChange.PostchangeState) > 0 && string(objectChange.PostchangeState) != "null" {
		postStateMap, err := convertJsonRawMessageToJxMap(objectChange.PostchangeState)
		if err != nil {
			return nil, fmt.Errorf("convert postchange state: %w", err)
		}
		postchangeState.Value = api.AuditLogPostchangeState(postStateMap)
		postchangeState.Null = false
	} else {
		postchangeState.Null = true
	}

	return &api.AuditLog{
		ID:               objectChange.ID,
		Employee:         toEmployeeDTO(objectChange.Employee),
		Action:           api.AuditLogAction(objectChange.Action),
		Time:             objectChange.Timestamp,
		TargetObjectType: toObjectTypeDTO(objectChange.ObjectType),
		PrechangeState:   prechangeState,
		PostchangeState:  postchangeState,
	}, nil
}

func (h *RestApiImplementation) GetAuditLogs(ctx context.Context, params api.GetAuditLogsParams) (*api.GetAuditLogsResponse, error) {
	if !params.ObjectTypeID.Set || !params.ObjectID.Set {
		return nil, h.NewError(ctx, errors.New("object type id and object id are required"))
	}

	auditLogs, err := h.auditUseCase.GetObjectChanges(ctx, models.ObjectTypeId(params.ObjectTypeID.Value), params.ObjectID.Value)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	auditLogsDTO := make([]api.AuditLog, len(auditLogs))
	for i, auditLog := range auditLogs {
		auditLogDTO, err := toAuditLog(auditLog)
		if err != nil {
			return nil, fmt.Errorf("convert audit log: %w", err)
		}
		auditLogsDTO[i] = *auditLogDTO
	}

	return &api.GetAuditLogsResponse{
		Data: auditLogsDTO,
	}, nil
}
