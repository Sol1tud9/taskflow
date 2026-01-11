package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/Sol1tud9/taskflow/internal/domain"
	"github.com/Sol1tud9/taskflow/internal/user/repository"
	"github.com/Sol1tud9/taskflow/pkg/logger"
	"go.uber.org/zap"
)

type EventPublisher interface {
	PublishUserCreated(ctx context.Context, event domain.UserCreatedEvent) error
	PublishUserUpdated(ctx context.Context, event domain.UserUpdatedEvent) error
}

type UserUseCase struct {
	userRepo  repository.UserRepository
	publisher EventPublisher
}

func NewUserUseCase(userRepo repository.UserRepository, publisher EventPublisher) *UserUseCase {
	return &UserUseCase{
		userRepo:  userRepo,
		publisher: publisher,
	}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, email, name string) (*domain.User, error) {
	now := time.Now()
	user := &domain.User{
		ID:        uuid.New().String(),
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	event := domain.UserCreatedEvent{
		UserID:    user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}
	if err := uc.publisher.PublishUserCreated(ctx, event); err != nil {
		logger.Error("failed to publish user.created event", zap.Error(err), zap.String("user_id", user.ID))
	}

	return user, nil
}

func (uc *UserUseCase) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, id, email, name string) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if email != "" {
		user.Email = email
	}
	if name != "" {
		user.Name = name
	}
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	event := domain.UserUpdatedEvent{
		UserID:    user.ID,
		Email:     user.Email,
		Name:      user.Name,
		UpdatedAt: user.UpdatedAt,
	}
	if err := uc.publisher.PublishUserUpdated(ctx, event); err != nil {
		logger.Error("failed to publish user.updated event", zap.Error(err), zap.String("user_id", user.ID))
	}

	return user, nil
}

