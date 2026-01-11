package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/Sol1tud9/taskflow/internal/domain"
	repoMocks "github.com/Sol1tud9/taskflow/internal/task/repository/mocks"
	usecaseMocks "github.com/Sol1tud9/taskflow/internal/task/usecase/mocks"
)

type TaskUseCaseSuite struct {
	suite.Suite
	ctx              context.Context
	taskRepo         *repoMocks.TaskRepository
	taskHistoryRepo  *repoMocks.TaskHistoryRepository
	publisher        *usecaseMocks.EventPublisher
	taskUseCase      *TaskUseCase
}

func (s *TaskUseCaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.taskRepo = repoMocks.NewTaskRepository(s.T())
	s.taskHistoryRepo = repoMocks.NewTaskHistoryRepository(s.T())
	s.publisher = usecaseMocks.NewEventPublisher(s.T())
	s.taskUseCase = NewTaskUseCase(s.taskRepo, s.taskHistoryRepo, s.publisher)
}

func (s *TaskUseCaseSuite) TestCreateTask_Success() {
	input := CreateTaskInput{
		Title:       "Test Task",
		Description: "Test Description",
		Priority:    "high",
		AssigneeID:  uuid.New().String(),
		CreatorID:   uuid.New().String(),
		TeamID:      uuid.New().String(),
		DueDate:     time.Now().Unix(),
	}

	s.taskRepo.On("Create", s.ctx, mock.MatchedBy(func(t *domain.Task) bool {
		return t.Title == input.Title && t.Description == input.Description
	})).Return(nil).Run(func(args mock.Arguments) {
		task := args.Get(1).(*domain.Task)
		task.ID = uuid.New().String()
		task.CreatedAt = time.Now()
		task.UpdatedAt = task.CreatedAt
	})

	s.publisher.On("PublishTaskCreated", s.ctx, mock.MatchedBy(func(e domain.TaskCreatedEvent) bool {
		return e.Title == input.Title
	})).Return(nil)

	result, err := s.taskUseCase.CreateTask(s.ctx, input)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), input.Title, result.Title)
	assert.Equal(s.T(), domain.TaskStatusTodo, result.Status)
}

func (s *TaskUseCaseSuite) TestCreateTask_RepositoryError() {
	input := CreateTaskInput{
		Title: "Test Task",
	}
	repoErr := errors.New("repository error")

	s.taskRepo.On("Create", s.ctx, mock.Anything).Return(repoErr)

	result, err := s.taskUseCase.CreateTask(s.ctx, input)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), repoErr, err)
}

func (s *TaskUseCaseSuite) TestGetTask_Success() {
	taskID := uuid.New().String()
	expectedTask := &domain.Task{
		ID:          taskID,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      domain.TaskStatusTodo,
		Priority:    domain.TaskPriorityHigh,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.taskRepo.On("GetByID", s.ctx, taskID).Return(expectedTask, nil)

	result, err := s.taskUseCase.GetTask(s.ctx, taskID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedTask.ID, result.ID)
	assert.Equal(s.T(), expectedTask.Title, result.Title)
}

func (s *TaskUseCaseSuite) TestUpdateTask_Success() {
	taskID := uuid.New().String()
	userID := uuid.New().String()

	existingTask := &domain.Task{
		ID:          taskID,
		Title:       "Old Title",
		Description: "Old Description",
		Status:      domain.TaskStatusTodo,
		Priority:    domain.TaskPriorityLow,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	input := UpdateTaskInput{
		Title:  "New Title",
		Status: "in_progress",
		UserID: userID,
	}

	s.taskRepo.On("GetByID", s.ctx, taskID).Return(existingTask, nil)
	s.taskRepo.On("Update", s.ctx, mock.MatchedBy(func(t *domain.Task) bool {
		return t.ID == taskID && t.Title == input.Title
	})).Return(nil)
	s.taskHistoryRepo.On("Create", s.ctx, mock.Anything).Return(nil)
	s.publisher.On("PublishTaskUpdated", s.ctx, mock.Anything).Return(nil)

	result, err := s.taskUseCase.UpdateTask(s.ctx, taskID, input)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), input.Title, result.Title)
	assert.Equal(s.T(), domain.TaskStatusInProgress, result.Status)
}

func (s *TaskUseCaseSuite) TestDeleteTask_Success() {
	taskID := uuid.New().String()

	s.taskRepo.On("Delete", s.ctx, taskID).Return(nil)

	err := s.taskUseCase.DeleteTask(s.ctx, taskID)

	assert.NoError(s.T(), err)
}

func TestTaskUseCaseSuite(t *testing.T) {
	suite.Run(t, new(TaskUseCaseSuite))
}

