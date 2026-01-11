package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/Sol1tud9/taskflow/internal/activity/repository"
	repoMocks "github.com/Sol1tud9/taskflow/internal/activity/repository/mocks"
	"github.com/Sol1tud9/taskflow/internal/domain"
)

type ActivityUseCaseSuite struct {
	suite.Suite
	ctx           context.Context
	activityRepo  *repoMocks.ActivityRepository
	activityUseCase *ActivityUseCase
}

func (s *ActivityUseCaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.activityRepo = repoMocks.NewActivityRepository(s.T())
	s.activityUseCase = NewActivityUseCase(s.activityRepo)
}

func (s *ActivityUseCaseSuite) TestRecordUserCreated_Success() {
	event := domain.UserCreatedEvent{
		UserID:    uuid.New().String(),
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
	}

	s.activityRepo.On("Create", s.ctx, mock.MatchedBy(func(a *domain.Activity) bool {
		return a.UserID == event.UserID && a.EntityType == domain.EntityTypeUser && a.Action == domain.ActionTypeCreated
	})).Return(nil)

	err := s.activityUseCase.RecordUserCreated(s.ctx, event)

	assert.NoError(s.T(), err)
}

func (s *ActivityUseCaseSuite) TestRecordTaskCreated_Success() {
	event := domain.TaskCreatedEvent{
		TaskID:     uuid.New().String(),
		Title:      "Test Task",
		CreatorID:  uuid.New().String(),
		AssigneeID: uuid.New().String(),
		TeamID:     uuid.New().String(),
		CreatedAt:  time.Now(),
	}

	s.activityRepo.On("Create", s.ctx, mock.MatchedBy(func(a *domain.Activity) bool {
		return a.EntityType == domain.EntityTypeTask && a.Action == domain.ActionTypeCreated
	})).Return(nil)

	err := s.activityUseCase.RecordTaskCreated(s.ctx, event)

	assert.NoError(s.T(), err)
}

func (s *ActivityUseCaseSuite) TestGetUserActivities_Success() {
	userID := uuid.New().String()
	expectedActivities := []*domain.Activity{
		{
			ID:         uuid.New().String(),
			UserID:     userID,
			EntityType: domain.EntityTypeUser,
			Action:     domain.ActionTypeCreated,
			CreatedAt:  time.Now(),
		},
	}

	filter := repository.ActivityFilter{
		Limit:  10,
		Offset: 0,
	}

	s.activityRepo.On("GetByUserID", s.ctx, userID, filter).Return(expectedActivities, 1, nil)

	result, total, err := s.activityUseCase.GetUserActivities(s.ctx, userID, 0, 0, 10, 0)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), 1, total)
	assert.Len(s.T(), result, 1)
}

func (s *ActivityUseCaseSuite) TestGetActivities_Success() {
	entityType := string(domain.EntityTypeTask)
	entityID := uuid.New().String()

	expectedActivities := []*domain.Activity{
		{
			ID:         uuid.New().String(),
			EntityType: domain.EntityTypeTask,
			EntityID:   entityID,
			Action:     domain.ActionTypeCreated,
			CreatedAt:  time.Now(),
		},
	}

	filter := repository.ActivityFilter{
		Limit:  10,
		Offset: 0,
	}

	s.activityRepo.On("GetByEntity", s.ctx, entityType, entityID, filter).Return(expectedActivities, 1, nil)

	result, total, err := s.activityUseCase.GetActivities(s.ctx, entityType, entityID, 0, 0, 10, 0)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), 1, total)
}

func TestActivityUseCaseSuite(t *testing.T) {
	suite.Run(t, new(ActivityUseCaseSuite))
}

