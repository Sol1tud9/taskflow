package repository

import (
	"context"

	"github.com/Sol1tud9/taskflow/internal/domain"
)

type ActivityRepository interface {
	Create(ctx context.Context, activity *domain.Activity) error
	GetByUserID(ctx context.Context, userID string, filter ActivityFilter) ([]*domain.Activity, int, error)
	GetByEntity(ctx context.Context, entityType, entityID string, filter ActivityFilter) ([]*domain.Activity, int, error)
	GetAll(ctx context.Context, filter ActivityFilter) ([]*domain.Activity, int, error)
}

type ActivityFilter struct {
	FromTimestamp int64
	ToTimestamp   int64
	Limit         int
	Offset        int
}

