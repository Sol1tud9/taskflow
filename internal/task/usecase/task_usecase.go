package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/Sol1tud9/taskflow/internal/domain"
	"github.com/Sol1tud9/taskflow/internal/task/repository"
	"github.com/Sol1tud9/taskflow/pkg/logger"
	"go.uber.org/zap"
)

type EventPublisher interface {
	PublishTaskCreated(ctx context.Context, event domain.TaskCreatedEvent) error
	PublishTaskUpdated(ctx context.Context, event domain.TaskUpdatedEvent) error
}

type TaskUseCase struct {
	taskRepo        repository.TaskRepository
	taskHistoryRepo repository.TaskHistoryRepository
	publisher       EventPublisher
}

func NewTaskUseCase(
	taskRepo repository.TaskRepository,
	taskHistoryRepo repository.TaskHistoryRepository,
	publisher EventPublisher,
) *TaskUseCase {
	return &TaskUseCase{
		taskRepo:        taskRepo,
		taskHistoryRepo: taskHistoryRepo,
		publisher:       publisher,
	}
}

func (uc *TaskUseCase) CreateTask(ctx context.Context, input CreateTaskInput) (*domain.Task, error) {
	now := time.Now()
	task := &domain.Task{
		ID:          uuid.New().String(),
		Title:       input.Title,
		Description: input.Description,
		Status:      domain.TaskStatusTodo,
		Priority:    domain.TaskPriority(input.Priority),
		AssigneeID:  input.AssigneeID,
		CreatorID:   input.CreatorID,
		TeamID:      input.TeamID,
		DueDate:     time.Unix(input.DueDate, 0),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	event := domain.TaskCreatedEvent{
		TaskID:     task.ID,
		Title:      task.Title,
		CreatorID:  task.CreatorID,
		AssigneeID: task.AssigneeID,
		TeamID:     task.TeamID,
		CreatedAt:  task.CreatedAt,
	}
	if err := uc.publisher.PublishTaskCreated(ctx, event); err != nil {
		logger.Error("failed to publish task.created event", zap.Error(err), zap.String("task_id", task.ID))
	}

	return task, nil
}

type CreateTaskInput struct {
	Title       string
	Description string
	Priority    string
	AssigneeID  string
	CreatorID   string
	TeamID      string
	DueDate     int64
}

func (uc *TaskUseCase) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	return uc.taskRepo.GetByID(ctx, id)
}

func (uc *TaskUseCase) ListTasks(ctx context.Context, filter repository.TaskFilter) ([]*domain.Task, int, error) {
	return uc.taskRepo.List(ctx, filter)
}

func (uc *TaskUseCase) UpdateTask(ctx context.Context, id string, input UpdateTaskInput) (*domain.Task, error) {
	task, err := uc.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	changes := make(map[string]struct {
		old string
		new string
	})

	if input.Title != "" && input.Title != task.Title {
		changes["title"] = struct {
			old string
			new string
		}{task.Title, input.Title}
		task.Title = input.Title
	}

	if input.Description != "" && input.Description != task.Description {
		changes["description"] = struct {
			old string
			new string
		}{task.Description, input.Description}
		task.Description = input.Description
	}

	if input.Status != "" && domain.TaskStatus(input.Status) != task.Status {
		changes["status"] = struct {
			old string
			new string
		}{string(task.Status), input.Status}
		task.Status = domain.TaskStatus(input.Status)
	}

	if input.Priority != "" && domain.TaskPriority(input.Priority) != task.Priority {
		changes["priority"] = struct {
			old string
			new string
		}{string(task.Priority), input.Priority}
		task.Priority = domain.TaskPriority(input.Priority)
	}

	if input.AssigneeID != "" && input.AssigneeID != task.AssigneeID {
		changes["assignee_id"] = struct {
			old string
			new string
		}{task.AssigneeID, input.AssigneeID}
		task.AssigneeID = input.AssigneeID
	}

	task.UpdatedAt = time.Now()

	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	for field, change := range changes {
		history := &domain.TaskHistory{
			ID:        uuid.New().String(),
			TaskID:    task.ID,
			UserID:    input.UserID,
			Field:     field,
			OldValue:  change.old,
			NewValue:  change.new,
			ChangedAt: task.UpdatedAt,
		}
		_ = uc.taskHistoryRepo.Create(ctx, history)

		event := domain.TaskUpdatedEvent{
			TaskID:    task.ID,
			UserID:    input.UserID,
			Field:     field,
			OldValue:  change.old,
			NewValue:  change.new,
			UpdatedAt: task.UpdatedAt,
		}
		if err := uc.publisher.PublishTaskUpdated(ctx, event); err != nil {
			logger.Error("failed to publish task.updated event", zap.Error(err), zap.String("task_id", task.ID))
		}
	}

	return task, nil
}

type UpdateTaskInput struct {
	Title       string
	Description string
	Status      string
	Priority    string
	AssigneeID  string
	DueDate     int64
	UserID      string
}

func (uc *TaskUseCase) DeleteTask(ctx context.Context, id string) error {
	return uc.taskRepo.Delete(ctx, id)
}

func (uc *TaskUseCase) GetTaskHistory(ctx context.Context, taskID string) ([]*domain.TaskHistory, error) {
	return uc.taskHistoryRepo.GetByTaskID(ctx, taskID)
}

