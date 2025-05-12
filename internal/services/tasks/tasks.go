package tasks

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/common"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/employee"
	"github.com/let-store-it/backend/internal/services/item"
	"github.com/let-store-it/backend/internal/services/organization"
	"github.com/let-store-it/backend/internal/services/storage"
	"github.com/let-store-it/backend/internal/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
	employee       *employee.EmployeeService
}

type TaskServiceConfig struct {
	Queries        *sqlc.Queries
	PGXPool        *pgxpool.Pool
	Auth           *auth.AuthService
	Org            *organization.OrganizationService
	StorageService *storage.StorageService
	ItemService    *item.ItemService
	EmployeeService *employee.EmployeeService
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
		employee:       cfg.EmployeeService,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, orgID uuid.UUID, task *models.Task) (*models.Task, error) {
	return telemetry.WithTrace(ctx, s.tracer, "CreateTask", func(ctx context.Context, span trace.Span) (*models.Task, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("unit.id", task.UnitID.String()),
			attribute.String("task.type", string(task.Type)),
			attribute.String("task.name", task.Name),
		)

		return database.WithTransaction(ctx, s.pgxpool, s.tracer, func(ctx context.Context, tx pgx.Tx) (*models.Task, error) {
			qtx := s.queries.WithTx(tx)

			createdTask, err := qtx.CreateTask(ctx, sqlc.CreateTaskParams{
				OrgID:            database.PgUUID(orgID),
				UnitID:           database.PgUUID(task.UnitID),
				Type:             sqlc.TaskType(task.Type),
				Name:             task.Name,
				Description:      database.PgTextPtr(task.Description),
				AssignedToUserID: database.PgUUIDPtr(task.AssignedToUserID),
			})
			if err != nil {
				return nil, services.MapDbErrorToService(err)
			}

			resultTask := toTask(createdTask)
			resultTask.Items = make([]*models.TaskItem, 0, len(task.Items))

			for _, item := range task.Items {
				sourceCell, err := s.item.GetItemInstanceFull(ctx, orgID, item.InstanceID)
				if err != nil {
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
					return nil, services.MapDbErrorToService(err)
				}
				taskItem := toTaskItem(taskItemDB)

				taskItem.SourceCell = sourceCell.Cell

				if item.TargetCellID != nil {
					cell, err := s.storageService.GetCellFull(ctx, orgID, *item.TargetCellID)
					if err != nil {
						return nil, fmt.Errorf("failed to get cell: %w", err)
					}
					taskItem.TargetCell = cell
				}

				instance, err := s.item.GetItemInstanceFull(ctx, orgID, item.InstanceID)
				if err != nil {
					return nil, fmt.Errorf("failed to get instance: %w", err)
				}
				taskItem.Instance = instance

				resultTask.Items = append(resultTask.Items, taskItem)
			}

			// Get the unit
			unit, err := s.org.GetUnitByID(ctx, orgID, resultTask.UnitID)
			if err != nil {
				return nil, fmt.Errorf("failed to get unit: %w", err)
			}
			resultTask.Unit = unit

			// Get assigned to if set
			if resultTask.AssignedToUserID != nil {
				assignedTo, err := s.employee.GetEmployee(ctx, orgID, *resultTask.AssignedToUserID)
				if err != nil {
					return nil, fmt.Errorf("failed to get assigned to: %w", err)
				}
				resultTask.AssignedTo = assignedTo
			}

			return resultTask, nil
		})
	})
}

func (s *TaskService) GetTaskById(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.Task, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetTaskById", func(ctx context.Context, span trace.Span) (*models.Task, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("task.id", id.String()),
		)

		task, err := s.queries.GetTaskById(ctx, sqlc.GetTaskByIdParams{
			OrgID: database.PgUUID(orgID),
			ID:    database.PgUUID(id),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		res := toTask(task)

		items, err := s.queries.GetTaskItems(ctx, sqlc.GetTaskItemsParams{
			OrgID:  database.PgUUID(orgID),
			TaskID: database.PgUUID(id),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		taskItems := make([]*models.TaskItem, len(items))
		for i, item := range items {
			taskItems[i] = toTaskItem(item)

			if taskItems[i].SourceCellID != nil {
				sourceCell, err := s.storageService.GetCellFull(ctx, orgID, *taskItems[i].SourceCellID)
				if err != nil {
					return nil, fmt.Errorf("failed to get source cell: %w", err)
				}
				taskItems[i].SourceCell = sourceCell
			}

			if taskItems[i].TargetCellID != nil {
				targetCell, err := s.storageService.GetCellFull(ctx, orgID, *taskItems[i].TargetCellID)
				if err != nil {
					return nil, fmt.Errorf("failed to get target cell: %w", err)
				}
				taskItems[i].TargetCell = targetCell
			}

			instance, err := s.item.GetItemInstanceFull(ctx, orgID, taskItems[i].InstanceID)
			if err != nil {
				return nil, fmt.Errorf("failed to get instance: %w", err)
			}
			taskItems[i].Instance = instance
		}
		res.Items = taskItems

		if res.AssignedToUserID != nil {
			assignedTo, err := s.employee.GetEmployee(ctx, orgID, *res.AssignedToUserID)
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
	})
}

func (s *TaskService) GetTasks(ctx context.Context, orgID uuid.UUID) ([]*models.Task, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetTasks", func(ctx context.Context, span trace.Span) ([]*models.Task, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
		)

		tasks, err := s.queries.GetTasks(ctx, database.PgUUID(orgID))
		if err != nil {
			if database.IsNotFound(err) {
				return nil, common.ErrNotFound
			}
			return nil, services.MapDbErrorToService(err)
		}

		models := make([]*models.Task, len(tasks))
		for i, task := range tasks {
			models[i] = toTask(task)
			if models[i].AssignedToUserID != nil {
				empl, err := s.employee.GetEmployee(ctx, orgID, *models[i].AssignedToUserID)
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
	})
}

func (s *TaskService) PickInstance(ctx context.Context, orgID uuid.UUID, taskID uuid.UUID, instanceID uuid.UUID) error {
	return telemetry.WithVoidTrace(ctx, s.tracer, "PickInstance", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("task.id", taskID.String()),
			attribute.String("instance.id", instanceID.String()),
		)

		err := s.item.SetItemInstanceStatus(ctx, &models.ItemInstance{
			OrgID:                 orgID,
			ID:                    instanceID,
			Status:                models.ItemInstanceStatusReserved,
			AffectedByOperationID: &taskID,
		})
		if err != nil {
			return fmt.Errorf("failed to set item instance status: %w", err)
		}

		err = s.item.SetInstanceCell(ctx, orgID, instanceID, nil)
		if err != nil {
			return fmt.Errorf("failed to set item instance cell: %w", err)
		}

		err = s.queries.SetTaskItemStatus(ctx, sqlc.SetTaskItemStatusParams{
			OrgID:          database.PgUUID(orgID),
			ItemInstanceID: database.PgUUID(instanceID),
			Status:         sqlc.TaskItemStatus(models.TaskItemStatusPicked),
		})
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		return nil
	})
}
