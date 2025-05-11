package tasks

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/item"
	"github.com/let-store-it/backend/internal/services/organization"
	"github.com/let-store-it/backend/internal/services/storage"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type TaskService struct {
	queries        *sqlc.Queries
	pgxpool        *pgxpool.Pool
	tracer         trace.Tracer
	auth           *auth.AuthService
	org            *organization.OrganizationService
	storageService *storage.StorageService
	item           *item.ItemService
}

type TaskServiceConfig struct {
	Queries        *sqlc.Queries
	PGXPool        *pgxpool.Pool
	Auth           *auth.AuthService
	Org            *organization.OrganizationService
	StorageService *storage.StorageService
	ItemService    *item.ItemService
}

func New(cfg TaskServiceConfig) *TaskService {
	return &TaskService{
		queries:        cfg.Queries,
		pgxpool:        cfg.PGXPool,
		auth:           cfg.Auth,
		org:            cfg.Org,
		tracer:         otel.GetTracerProvider().Tracer("tasks-service"),
		storageService: cfg.StorageService,
		item:           cfg.ItemService,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, orgID uuid.UUID, task *models.Task) (*models.Task, error) {
	ctx, span := s.tracer.Start(ctx, "CreateTask",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
		),
	)
	defer span.End()

	tx, err := s.pgxpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to begin transaction")
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := s.queries.WithTx(tx)

	createdTask, err := qtx.CreateTask(ctx, sqlc.CreateTaskParams{
		OrgID:            database.PgUUID(orgID),
		UnitID:           database.PgUUID(task.UnitID),
		Type:             string(task.Type),
		Name:             task.Name,
		Description:      database.PgTextPtr(task.Description),
		AssignedToUserID: database.PgUUIDPtr(task.AssignedToUserID),
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create task")
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	for _, item := range task.Items {
		sourceCell, err := s.item.GetItemInstanceFull(ctx, orgID, item.InstanceID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to get instance cell")
			return nil, fmt.Errorf("failed to get instance cell: %w", err)
		}

		taskItemDB, err := qtx.CreateTaskItem(ctx, sqlc.CreateTaskItemParams{
			OrgID:             database.PgUUID(orgID),
			TaskID:            createdTask.ID,
			ItemInstanceID:    database.PgUUID(item.InstanceID),
			DestinationCellID: database.PgUUIDPtr(item.TargetCellID),
			SourceCellID:      database.PgUUIDPtr(sourceCell.CellID),
		})
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to create task item")
			return nil, fmt.Errorf("failed to create task item: %w", err)
		}
		taskItem := toTaskItem(taskItemDB)

		taskItem.SourceCell = sourceCell.Cell

		cell, err := s.storageService.GetCellFull(ctx, orgID, *item.TargetCellID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to get cell")
			return nil, fmt.Errorf("failed to get cell: %w", err)
		}
		taskItem.TargetCell = cell

		instance, err := s.item.GetItemInstanceFull(ctx, orgID, item.InstanceID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to get instance")
			return nil, fmt.Errorf("failed to get instance: %w", err)
		}
		taskItem.Instance = instance
	}

	err = tx.Commit(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to commit transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return toTask(createdTask), nil
}

func (s *TaskService) GetTaskById(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.Task, error) {
	ctx, span := s.tracer.Start(ctx, "GetTaskById",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("task.id", id.String()),
		),
	)
	defer span.End()

	task, err := s.queries.GetTaskById(ctx, sqlc.GetTaskByIdParams{
		OrgID: database.PgUUID(orgID),
		ID:    database.PgUUID(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	res, err := toTask(task), nil
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	items, err := s.queries.GetTaskItems(ctx, sqlc.GetTaskItemsParams{
		OrgID:  database.PgUUID(orgID),
		TaskID: database.PgUUID(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get task items: %w", err)
	}

	taskItems := make([]*models.TaskItem, len(items))
	for i, item := range items {
		taskItems[i] = toTaskItem(item)
		taskItems[i].SourceCell, err = s.storageService.GetCellFull(ctx, orgID, *taskItems[i].SourceCellID)
		if err != nil {
			return nil, fmt.Errorf("failed to get source cell: %w", err)
		}
		taskItems[i].TargetCell, err = s.storageService.GetCellFull(ctx, orgID, *taskItems[i].TargetCellID)
		if err != nil {
			return nil, fmt.Errorf("failed to get target cell: %w", err)
		}
		taskItems[i].Instance, err = s.item.GetItemInstanceFull(ctx, orgID, taskItems[i].InstanceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get instance: %w", err)
		}
	}
	res.Items = taskItems

	if res.AssignedToUserID != nil {
		assignedTo, err := s.auth.GetEmployee(ctx, orgID, *res.AssignedToUserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get assigned to: %w", err)
		}
		res.AssignedTo = assignedTo
	}

	unit, err := s.org.GetUnitByID(ctx, orgID, res.UnitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unit: %w", err)
	}
	res.Unit = unit

	return res, nil
}

func (s *TaskService) GetTasks(ctx context.Context, orgID uuid.UUID) ([]*models.Task, error) {
	ctx, span := s.tracer.Start(ctx, "GetTasks",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
		),
	)
	defer span.End()

	tasks, err := s.queries.GetTasks(ctx, database.PgUUID(orgID))
	if err != nil {
		if database.IsNotFound(err) {
			return nil, services.ErrNotFoundError
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get tasks")
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	models := make([]*models.Task, len(tasks))
	for i, task := range tasks {
		models[i] = toTask(task)
		if models[i].AssignedToUserID != nil {
			empl, err := s.auth.GetEmployee(ctx, orgID, *models[i].AssignedToUserID)
			if err != nil {
				return nil, err
			}
			models[i].AssignedTo = empl
		}
		unit, err := s.org.GetUnitByID(ctx, orgID, models[i].UnitID)
		if err != nil {
			return nil, err
		}
		models[i].Unit = unit
	}

	return models, nil
}

func (s *TaskService) PickInstance(ctx context.Context, orgID uuid.UUID, taskID uuid.UUID, instanceID uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "PickInstance",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("task.id", taskID.String()),
			attribute.String("instance.id", instanceID.String()),
		),
	)
	defer span.End()

	err := s.item.SetItemInstanceStatus(ctx, &models.ItemInstance{
		OrgID:                 orgID,
		ID:                    instanceID,
		Status:                models.ItemInstanceStatusReserved,
		AffectedByOperationID: &taskID,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to set item instance status")
		return fmt.Errorf("failed to set item instance status: %w", err)
	}

	// set cell
	err = s.item.SetInstanceCell(ctx, orgID, instanceID, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to set item instance cell")
		return fmt.Errorf("failed to set item instance cell: %w", err)
	}

	// set task item status to picked
	err = s.queries.SetTaskItemStatus(ctx, sqlc.SetTaskItemStatusParams{
		OrgID:          database.PgUUID(orgID),
		ItemInstanceID: database.PgUUID(instanceID),
		Status:         string(models.TaskItemStatusPicked),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to set task item status")
		return fmt.Errorf("failed to set task item status: %w", err)
	}
	return nil
}
