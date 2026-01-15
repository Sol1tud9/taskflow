package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
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

type ActivityUseCase struct {
	activityRepo ActivityRepository
}

func NewActivityUseCase(activityRepo ActivityRepository) *ActivityUseCase {
	return &ActivityUseCase{
		activityRepo: activityRepo,
	}
}

func (uc *ActivityUseCase) RecordUserCreated(ctx context.Context, event domain.UserCreatedEvent) error {
	metadata, _ := json.Marshal(event)
	activity := &domain.Activity{
		ID:         uuid.New().String(),
		UserID:     event.UserID,
		EntityType: domain.EntityTypeUser,
		EntityID:   event.UserID,
		Action:     domain.ActionTypeCreated,
		Metadata:   string(metadata),
		CreatedAt:  event.CreatedAt,
	}
	return uc.activityRepo.Create(ctx, activity)
}

func (uc *ActivityUseCase) RecordUserUpdated(ctx context.Context, event domain.UserUpdatedEvent) error {
	metadata, _ := json.Marshal(event)
	activity := &domain.Activity{
		ID:         uuid.New().String(),
		UserID:     event.UserID,
		EntityType: domain.EntityTypeUser,
		EntityID:   event.UserID,
		Action:     domain.ActionTypeUpdated,
		Metadata:   string(metadata),
		CreatedAt:  event.UpdatedAt,
	}
	return uc.activityRepo.Create(ctx, activity)
}

func (uc *ActivityUseCase) RecordTaskCreated(ctx context.Context, event domain.TaskCreatedEvent) error {
	metadata, _ := json.Marshal(event)
	activity := &domain.Activity{
		ID:         uuid.New().String(),
		UserID:     event.CreatorID,
		EntityType: domain.EntityTypeTask,
		EntityID:   event.TaskID,
		Action:     domain.ActionTypeCreated,
		Metadata:   string(metadata),
		CreatedAt:  event.CreatedAt,
	}
	return uc.activityRepo.Create(ctx, activity)
}

func (uc *ActivityUseCase) RecordTaskUpdated(ctx context.Context, event domain.TaskUpdatedEvent) error {
	metadata, _ := json.Marshal(event)
	activity := &domain.Activity{
		ID:         uuid.New().String(),
		UserID:     event.UserID,
		EntityType: domain.EntityTypeTask,
		EntityID:   event.TaskID,
		Action:     domain.ActionTypeUpdated,
		Metadata:   string(metadata),
		CreatedAt:  event.UpdatedAt,
	}
	return uc.activityRepo.Create(ctx, activity)
}

func (uc *ActivityUseCase) GetUserActivities(ctx context.Context, userID string, from, to int64, limit, offset int) ([]*domain.Activity, int, error) {
	filter := ActivityFilter{
		FromTimestamp: from,
		ToTimestamp:   to,
		Limit:         limit,
		Offset:        offset,
	}
	return uc.activityRepo.GetByUserID(ctx, userID, filter)
}

func (uc *ActivityUseCase) GetActivities(ctx context.Context, entityType, entityID string, from, to int64, limit, offset int) ([]*domain.Activity, int, error) {
	filter := ActivityFilter{
		FromTimestamp: from,
		ToTimestamp:   to,
		Limit:         limit,
		Offset:        offset,
	}

	if entityType != "" && entityID != "" {
		return uc.activityRepo.GetByEntity(ctx, entityType, entityID, filter)
	}

	return uc.activityRepo.GetAll(ctx, filter)
}

func (uc *ActivityUseCase) RecordActivity(ctx context.Context, userID string, entityType domain.EntityType, entityID string, action domain.ActionType, metadata string) error {
	activity := &domain.Activity{
		ID:         uuid.New().String(),
		UserID:     userID,
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		Metadata:   metadata,
		CreatedAt:  time.Now(),
	}
	return uc.activityRepo.Create(ctx, activity)
}

