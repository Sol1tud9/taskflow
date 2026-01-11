package repository

import (
	"context"

	"github.com/Sol1tud9/taskflow/internal/domain"
)

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, id string) (*domain.Task, error)
	List(ctx context.Context, filter TaskFilter) ([]*domain.Task, int, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id string) error
}

type TaskFilter struct {
	TeamID     string
	AssigneeID string
	Status     string
	Limit      int
	Offset     int
}

type TaskHistoryRepository interface {
	Create(ctx context.Context, history *domain.TaskHistory) error
	GetByTaskID(ctx context.Context, taskID string) ([]*domain.TaskHistory, error)
}

