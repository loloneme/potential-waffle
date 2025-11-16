package create_team_test

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/team"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/user"
	"github.com/loloneme/potential-waffle/internal/infrastructure/utils/test_env"
	"github.com/loloneme/potential-waffle/internal/usecase/create_team"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testEnv struct {
	ctx      context.Context
	db       *sqlx.DB
	userRepo *user.Repository
	teamRepo *team.Repository
	service  *create_team.Service
}

func setupTest(t *testing.T) *testEnv {
	t.Helper()

	ctx := context.Background()

	db, err := test_env.NewTestDatabaseConnection(ctx)
	require.NoError(t, err)

	userRepo := user.NewRepository(db)
	teamRepo := team.NewRepository(db)
	service := create_team.New(teamRepo, userRepo)

	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "TRUNCATE users, teams, pull_requests, reviewers RESTART IDENTITY CASCADE")
		_ = db.Close()
	})

	return &testEnv{
		ctx:      ctx,
		db:       db,
		userRepo: userRepo,
		teamRepo: teamRepo,
		service:  service,
	}
}

func TestService_CreateTeam_Successful(t *testing.T) {
	env := setupTest(t)

	teamName := "test-team"
	members := []models.User{
		{ID: "user-1", Username: "user1", IsActive: true, TeamName: teamName},
		{ID: "user-2", Username: "user2", IsActive: true, TeamName: teamName},
		{ID: "user-3", Username: "user3", IsActive: false, TeamName: teamName},
	}

	team := &models.Team{
		TeamName: teamName,
		Members:  members,
	}

	createdTeam, err := env.service.CreateTeam(env.ctx, team)
	require.NoError(t, err)

	assert.Equal(t, teamName, createdTeam.TeamName)
	assert.Len(t, createdTeam.Members, 3)
	assert.Equal(t, "user-1", createdTeam.Members[0].ID)
	assert.Equal(t, "user1", createdTeam.Members[0].Username)
	assert.True(t, createdTeam.Members[0].IsActive)
	assert.Equal(t, "user-2", createdTeam.Members[1].ID)
	assert.Equal(t, "user-3", createdTeam.Members[2].ID)
	assert.False(t, createdTeam.Members[2].IsActive)

	exists, err := env.teamRepo.Exists(env.ctx, teamName)
	require.NoError(t, err)
	assert.True(t, exists)

	for _, member := range members {
		teamNameFromDB, err := env.userRepo.GetUserTeamName(env.ctx, member.ID)
		require.NoError(t, err)
		assert.Equal(t, teamName, teamNameFromDB)
	}
}

func TestService_CreateTeam_EmptyMembers(t *testing.T) {
	env := setupTest(t)

	teamName := "empty-team"
	team := &models.Team{
		TeamName: teamName,
		Members:  []models.User{},
	}

	createdTeam, err := env.service.CreateTeam(env.ctx, team)
	require.NoError(t, err)

	assert.Equal(t, teamName, createdTeam.TeamName)
	assert.Len(t, createdTeam.Members, 0)

	exists, err := env.teamRepo.Exists(env.ctx, teamName)
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestService_CreateTeam_TeamAlreadyExists(t *testing.T) {
	env := setupTest(t)

	teamName := "existing-team"
	members := []models.User{
		{ID: "user-1", Username: "user1", IsActive: true, TeamName: teamName},
	}

	team := &models.Team{
		TeamName: teamName,
		Members:  members,
	}

	_, err := env.service.CreateTeam(env.ctx, team)
	require.NoError(t, err)

	newMembers := []models.User{
		{ID: "user-2", Username: "user2", IsActive: true, TeamName: teamName},
		{ID: "user-3", Username: "user3", IsActive: true, TeamName: teamName},
	}

	team2 := &models.Team{
		TeamName: teamName,
		Members:  newMembers,
	}

	createdTeam, err := env.service.CreateTeam(env.ctx, team2)
	require.NoError(t, err)

	assert.Equal(t, teamName, createdTeam.TeamName)
	assert.Len(t, createdTeam.Members, 2)
}

func TestService_CreateTeam_MultipleTeams(t *testing.T) {
	env := setupTest(t)

	team1Name := "team-1"
	team1 := &models.Team{
		TeamName: team1Name,
		Members: []models.User{
			{ID: "user-1", Username: "user1", IsActive: true, TeamName: team1Name},
		},
	}

	createdTeam1, err := env.service.CreateTeam(env.ctx, team1)
	require.NoError(t, err)
	assert.Equal(t, team1Name, createdTeam1.TeamName)

	team2Name := "team-2"
	team2 := &models.Team{
		TeamName: team2Name,
		Members: []models.User{
			{ID: "user-2", Username: "user2", IsActive: true, TeamName: team2Name},
		},
	}

	createdTeam2, err := env.service.CreateTeam(env.ctx, team2)
	require.NoError(t, err)
	assert.Equal(t, team2Name, createdTeam2.TeamName)

	exists1, err := env.teamRepo.Exists(env.ctx, team1Name)
	require.NoError(t, err)
	assert.True(t, exists1)

	exists2, err := env.teamRepo.Exists(env.ctx, team2Name)
	require.NoError(t, err)
	assert.True(t, exists2)
}
