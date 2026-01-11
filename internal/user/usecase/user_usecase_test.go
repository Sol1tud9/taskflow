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
	repoMocks "github.com/Sol1tud9/taskflow/internal/user/repository/mocks"
	usecaseMocks "github.com/Sol1tud9/taskflow/internal/user/usecase/mocks"
)

type UserUseCaseSuite struct {
	suite.Suite
	ctx            context.Context
	userRepo       *repoMocks.UserRepository
	publisher      *usecaseMocks.EventPublisher
	userUseCase    *UserUseCase
}

func (s *UserUseCaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.userRepo = repoMocks.NewUserRepository(s.T())
	s.publisher = usecaseMocks.NewEventPublisher(s.T())
	s.userUseCase = NewUserUseCase(s.userRepo, s.publisher)
}

func (s *UserUseCaseSuite) TestCreateUser_Success() {
	email := "test@example.com"
	name := "Test User"

	user := &domain.User{
		ID:        uuid.New().String(),
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.userRepo.On("Create", s.ctx, mock.MatchedBy(func(u *domain.User) bool {
		return u.Email == email && u.Name == name
	})).Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*domain.User)
		user.ID = u.ID
		user.CreatedAt = u.CreatedAt
		user.UpdatedAt = u.UpdatedAt
	})

	s.publisher.On("PublishUserCreated", s.ctx, mock.MatchedBy(func(e domain.UserCreatedEvent) bool {
		return e.Email == email && e.Name == name
	})).Return(nil)

	result, err := s.userUseCase.CreateUser(s.ctx, email, name)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), email, result.Email)
	assert.Equal(s.T(), name, result.Name)
	assert.NotEmpty(s.T(), result.ID)
}

func (s *UserUseCaseSuite) TestCreateUser_RepositoryError() {
	email := "test@example.com"
	name := "Test User"
	repoErr := errors.New("repository error")

	s.userRepo.On("Create", s.ctx, mock.Anything).Return(repoErr)

	result, err := s.userUseCase.CreateUser(s.ctx, email, name)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), repoErr, err)
}

func (s *UserUseCaseSuite) TestGetUser_Success() {
	userID := uuid.New().String()
	expectedUser := &domain.User{
		ID:        userID,
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.userRepo.On("GetByID", s.ctx, userID).Return(expectedUser, nil)

	result, err := s.userUseCase.GetUser(s.ctx, userID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedUser.ID, result.ID)
	assert.Equal(s.T(), expectedUser.Email, result.Email)
}

func (s *UserUseCaseSuite) TestGetUser_NotFound() {
	userID := uuid.New().String()
	repoErr := errors.New("user not found")

	s.userRepo.On("GetByID", s.ctx, userID).Return(nil, repoErr)

	result, err := s.userUseCase.GetUser(s.ctx, userID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), repoErr, err)
}

func (s *UserUseCaseSuite) TestUpdateUser_Success() {
	userID := uuid.New().String()
	newEmail := "new@example.com"
	newName := "New Name"

	existingUser := &domain.User{
		ID:        userID,
		Email:     "old@example.com",
		Name:      "Old Name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.userRepo.On("GetByID", s.ctx, userID).Return(existingUser, nil)
	s.userRepo.On("Update", s.ctx, mock.MatchedBy(func(u *domain.User) bool {
		return u.ID == userID && u.Email == newEmail && u.Name == newName
	})).Return(nil)
	s.publisher.On("PublishUserUpdated", s.ctx, mock.MatchedBy(func(e domain.UserUpdatedEvent) bool {
		return e.UserID == userID && e.Email == newEmail && e.Name == newName
	})).Return(nil)

	result, err := s.userUseCase.UpdateUser(s.ctx, userID, newEmail, newName)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), newEmail, result.Email)
	assert.Equal(s.T(), newName, result.Name)
}

func (s *UserUseCaseSuite) TestUpdateUser_PartialUpdate() {
	userID := uuid.New().String()
	newEmail := "new@example.com"

	existingUser := &domain.User{
		ID:        userID,
		Email:     "old@example.com",
		Name:      "Old Name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.userRepo.On("GetByID", s.ctx, userID).Return(existingUser, nil)
	s.userRepo.On("Update", s.ctx, mock.MatchedBy(func(u *domain.User) bool {
		return u.ID == userID && u.Email == newEmail && u.Name == existingUser.Name
	})).Return(nil)
	s.publisher.On("PublishUserUpdated", s.ctx, mock.Anything).Return(nil)

	result, err := s.userUseCase.UpdateUser(s.ctx, userID, newEmail, "")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), newEmail, result.Email)
	assert.Equal(s.T(), existingUser.Name, result.Name)
}

func (s *UserUseCaseSuite) TestUpdateUser_NotFound() {
	userID := uuid.New().String()
	repoErr := errors.New("user not found")

	s.userRepo.On("GetByID", s.ctx, userID).Return(nil, repoErr)

	result, err := s.userUseCase.UpdateUser(s.ctx, userID, "new@example.com", "New Name")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), repoErr, err)
}

func TestUserUseCaseSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseSuite))
}

