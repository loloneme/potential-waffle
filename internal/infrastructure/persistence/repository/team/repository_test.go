package team_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/team"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/team/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTeamRepository_CreateTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockteamRepository(ctrl)
	ctx := context.Background()
	var tx *sqlx.Tx

	t.Run("successful create", func(t *testing.T) {
		tm := models.Team{
			TeamName: "team-1",
		}
		expectedTeam := models.Team{
			TeamName: "team-1",
		}

		mockRepo.EXPECT().
			CreateTeam(gomock.Any(), gomock.Any(), tm).
			Return(expectedTeam, nil)

		result, err := mockRepo.CreateTeam(ctx, tx, tm)
		assert.NoError(t, err)
		assert.Equal(t, expectedTeam.TeamName, result.TeamName)
	})

	t.Run("team already exists", func(t *testing.T) {
		tm := models.Team{
			TeamName: "team-1",
		}

		mockRepo.EXPECT().
			CreateTeam(gomock.Any(), gomock.Any(), tm).
			Return(tm, team.ErrAlreadyExists)

		result, err := mockRepo.CreateTeam(ctx, tx, tm)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, team.ErrAlreadyExists))
		assert.Equal(t, tm.TeamName, result.TeamName)
	})
}

func TestTeamRepository_FindTeamByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockteamRepository(ctrl)
	ctx := context.Background()
	teamName := "team-1"

	t.Run("successful find", func(t *testing.T) {
		expectedTeam := models.Team{
			TeamName: "team-1",
		}

		mockRepo.EXPECT().
			FindTeamByID(gomock.Any(), teamName).
			Return(expectedTeam, nil)

		tm, err := mockRepo.FindTeamByID(ctx, teamName)
		assert.NoError(t, err)
		assert.Equal(t, expectedTeam.TeamName, tm.TeamName)
	})

	t.Run("team not found", func(t *testing.T) {
		mockRepo.EXPECT().
			FindTeamByID(gomock.Any(), teamName).
			Return(models.Team{}, team.ErrNotFound)

		tm, err := mockRepo.FindTeamByID(ctx, teamName)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, team.ErrNotFound))
		assert.Equal(t, models.Team{}, tm)
	})
}

func TestTeamRepository_WithTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockteamRepository(ctrl)
	ctx := context.Background()

	t.Run("successful transaction", func(t *testing.T) {
		mockRepo.EXPECT().
			WithTx(gomock.Any(), gomock.Any()).
			Return(nil)

		err := mockRepo.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("transaction with error", func(t *testing.T) {
		testErr := errors.New("test error")
		mockRepo.EXPECT().
			WithTx(gomock.Any(), gomock.Any()).
			Return(testErr)

		err := mockRepo.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
			return testErr
		})
		assert.Error(t, err)
		assert.Equal(t, testErr, err)
	})
}
