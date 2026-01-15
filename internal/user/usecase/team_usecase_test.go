package usecase_test

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
	userUsecase "github.com/Sol1tud9/taskflow/internal/user/usecase"
	usecaseMocks "github.com/Sol1tud9/taskflow/internal/user/usecase/mocks"
)

type TeamUseCaseSuite struct {
	suite.Suite
	ctx              context.Context
	teamRepo         *repoMocks.TeamRepository
	teamMemberRepo   *repoMocks.TeamMemberRepository
	publisher        *usecaseMocks.TeamEventPublisher
	teamUseCase      *userUsecase.TeamUseCase
}

func (s *TeamUseCaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.teamRepo = repoMocks.NewTeamRepository(s.T())
	s.teamMemberRepo = repoMocks.NewTeamMemberRepository(s.T())
	s.publisher = usecaseMocks.NewTeamEventPublisher(s.T())
	s.teamUseCase = userUsecase.NewTeamUseCase(s.teamRepo, s.teamMemberRepo, s.publisher)
}

func (s *TeamUseCaseSuite) TestCreateTeam_Success() {
	name := "Test Team"
	ownerID := uuid.New().String()

	s.teamRepo.On("Create", s.ctx, mock.MatchedBy(func(t *domain.Team) bool {
		return t.Name == name && t.OwnerID == ownerID
	})).Return(nil).Run(func(args mock.Arguments) {
		team := args.Get(1).(*domain.Team)
		team.ID = uuid.New().String()
		team.CreatedAt = time.Now()
		team.UpdatedAt = team.CreatedAt
	})

	s.teamMemberRepo.On("Add", s.ctx, mock.MatchedBy(func(m *domain.TeamMember) bool {
		return m.UserID == ownerID && m.Role == "owner"
	})).Return(nil)

	s.publisher.On("PublishTeamUpdated", s.ctx, mock.MatchedBy(func(e domain.TeamUpdatedEvent) bool {
		return e.Name == name && e.OwnerID == ownerID
	})).Return(nil)

	result, err := s.teamUseCase.CreateTeam(s.ctx, name, ownerID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), name, result.Name)
	assert.Equal(s.T(), ownerID, result.OwnerID)
	assert.NotEmpty(s.T(), result.ID)
}

func (s *TeamUseCaseSuite) TestCreateTeam_RepositoryError() {
	name := "Test Team"
	ownerID := uuid.New().String()
	repoErr := errors.New("repository error")

	s.teamRepo.On("Create", s.ctx, mock.Anything).Return(repoErr)

	result, err := s.teamUseCase.CreateTeam(s.ctx, name, ownerID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), repoErr, err)
}

func (s *TeamUseCaseSuite) TestGetTeam_Success() {
	teamID := uuid.New().String()
	expectedTeam := &domain.Team{
		ID:        teamID,
		Name:      "Test Team",
		OwnerID:   uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.teamRepo.On("GetByID", s.ctx, teamID).Return(expectedTeam, nil)

	result, err := s.teamUseCase.GetTeam(s.ctx, teamID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedTeam.ID, result.ID)
	assert.Equal(s.T(), expectedTeam.Name, result.Name)
}

func (s *TeamUseCaseSuite) TestAddTeamMember_Success() {
	teamID := uuid.New().String()
	userID := uuid.New().String()
	role := "member"

	s.teamMemberRepo.On("Add", s.ctx, mock.MatchedBy(func(m *domain.TeamMember) bool {
		return m.TeamID == teamID && m.UserID == userID && m.Role == role
	})).Return(nil).Run(func(args mock.Arguments) {
		member := args.Get(1).(*domain.TeamMember)
		member.ID = uuid.New().String()
		member.JoinedAt = time.Now()
	})

	result, err := s.teamUseCase.AddTeamMember(s.ctx, teamID, userID, role)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), teamID, result.TeamID)
	assert.Equal(s.T(), userID, result.UserID)
	assert.Equal(s.T(), role, result.Role)
}

func (s *TeamUseCaseSuite) TestGetTeamMembers_Success() {
	teamID := uuid.New().String()
	expectedMembers := []*domain.TeamMember{
		{
			ID:       uuid.New().String(),
			TeamID:   teamID,
			UserID:   uuid.New().String(),
			Role:     "owner",
			JoinedAt: time.Now(),
		},
		{
			ID:       uuid.New().String(),
			TeamID:   teamID,
			UserID:   uuid.New().String(),
			Role:     "member",
			JoinedAt: time.Now(),
		},
	}

	s.teamMemberRepo.On("GetByTeamID", s.ctx, teamID).Return(expectedMembers, nil)

	result, err := s.teamUseCase.GetTeamMembers(s.ctx, teamID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 2)
	assert.Equal(s.T(), expectedMembers[0].ID, result[0].ID)
}

func TestTeamUseCaseSuite(t *testing.T) {
	suite.Run(t, new(TeamUseCaseSuite))
}

