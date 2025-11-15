package merge_pr_test

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/pull_request"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/team"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/user"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/specification/status"
	"github.com/loloneme/potential-waffle/internal/infrastructure/utils/test_env"
	rpc_errors "github.com/loloneme/potential-waffle/internal/rpc/errors"
	"github.com/loloneme/potential-waffle/internal/usecase/merge_pr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testEnv struct {
	ctx      context.Context
	db       *sqlx.DB
	userRepo *user.Repository
	prRepo   *pull_request.Repository
	teamRepo *team.Repository
	service  *merge_pr.Service
}

func setupTest(t *testing.T) *testEnv {
	t.Helper()

	ctx := context.Background()

	db, err := test_env.NewTestDatabaseConnection(ctx)
	require.NoError(t, err)

	userRepo := user.NewRepository(db)
	prRepo := pull_request.NewRepository(db)
	teamRepo := team.NewRepository(db)
	service := merge_pr.New(prRepo)

	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "TRUNCATE users, teams, pull_requests, reviewers RESTART IDENTITY CASCADE")
		_ = db.Close()
	})

	return &testEnv{
		ctx:      ctx,
		db:       db,
		userRepo: userRepo,
		prRepo:   prRepo,
		teamRepo: teamRepo,
		service:  service,
	}
}

func TestService_MergePullRequest_Successful(t *testing.T) {
	env := setupTest(t)

	// Подготовка данных
	teamName := "test-team"
	err := env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := env.teamRepo.CreateTeam(ctx, tx, models.Team{TeamName: teamName})
		return err
	})
	require.NoError(t, err)

	authorID := "author-1"
	reviewer1ID := "reviewer-1"
	reviewer2ID := "reviewer-2"

	err = env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		users := []models.User{
			{ID: authorID, Username: "author", IsActive: true, TeamName: teamName},
			{ID: reviewer1ID, Username: "reviewer1", IsActive: true, TeamName: teamName},
			{ID: reviewer2ID, Username: "reviewer2", IsActive: true, TeamName: teamName},
		}
		_, err := env.userRepo.UpsertUsers(ctx, tx, users)
		return err
	})
	require.NoError(t, err)

	// Создаем PR
	prID := "pr-1"
	err = env.prRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		pr := &models.PullRequest{
			ID:       prID,
			Name:     "Test PR",
			AuthorID: authorID,
			Status:   &models.Status{Name: "OPEN"},
		}

		foundStatus, err := env.prRepo.FindStatus(ctx, status.NewGetStatusByNameSpecification("OPEN"))
		require.NoError(t, err)
		pr.StatusID = foundStatus.ID

		_, err = env.prRepo.InsertPullRequest(ctx, tx, pr)
		if err != nil {
			return err
		}

		reviewers := []string{reviewer1ID, reviewer2ID}
		return env.prRepo.InsertReviewers(ctx, tx, prID, reviewers)
	})
	require.NoError(t, err)

	// Мержим PR
	mergedPR, err := env.service.MergePullRequest(env.ctx, prID, "MERGED")
	require.NoError(t, err)

	// Проверяем результат
	assert.Equal(t, prID, mergedPR.ID)
	assert.Equal(t, "Test PR", mergedPR.Name)
	assert.Equal(t, authorID, mergedPR.AuthorID)
	assert.NotNil(t, mergedPR.Status)
	assert.Equal(t, "MERGED", mergedPR.Status.Name)

	// Проверяем, что статус действительно изменен в БД
	prFromDB, err := env.prRepo.GetPRByID(env.ctx, prID)
	require.NoError(t, err)
	assert.Equal(t, "MERGED", prFromDB.Status.Name)
}

func TestService_MergePullRequest_NotFound(t *testing.T) {
	env := setupTest(t)

	nonExistentPRID := "non-existent-pr"

	_, err := env.service.MergePullRequest(env.ctx, nonExistentPRID, "MERGED")
	require.Error(t, err)
	var notFoundErr *rpc_errors.NotFoundError
	assert.ErrorAs(t, err, &notFoundErr)
}

func TestService_MergePullRequest_StatusNotFound(t *testing.T) {
	env := setupTest(t)

	// Подготовка данных
	teamName := "test-team-status-not-found"
	err := env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := env.teamRepo.CreateTeam(ctx, tx, models.Team{TeamName: teamName})
		return err
	})
	require.NoError(t, err)

	authorID := "author-status-not-found"
	err = env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		users := []models.User{
			{ID: authorID, Username: "author-status-not-found", IsActive: true, TeamName: teamName},
		}
		_, err := env.userRepo.UpsertUsers(ctx, tx, users)
		return err
	})
	require.NoError(t, err)

	// Создаем PR
	prID := "pr-status-not-found"
	err = env.prRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		pr := &models.PullRequest{
			ID:       prID,
			Name:     "PR for status not found test",
			AuthorID: authorID,
			Status:   &models.Status{Name: "OPEN"},
		}

		foundStatus, err := env.prRepo.FindStatus(ctx, status.NewGetStatusByNameSpecification("OPEN"))
		require.NoError(t, err)
		pr.StatusID = foundStatus.ID

		_, err = env.prRepo.InsertPullRequest(ctx, tx, pr)
		return err
	})
	require.NoError(t, err)

	// Проверяем, что PR существует
	exists, err := env.prRepo.PullRequestExists(env.ctx, prID)
	require.NoError(t, err)
	require.True(t, exists, "PR should exist")

	// Пытаемся изменить статус на несуществующий статус
	nonExistentStatus := "NONEXISTENT_STATUS"
	_, err = env.service.MergePullRequest(env.ctx, prID, nonExistentStatus)
	require.Error(t, err)

	// Проверяем, что ошибка связана с тем, что статус не найден
	assert.ErrorIs(t, err, pull_request.ErrStatusNotFound, "Error should be ErrStatusNotFound")
	assert.Contains(t, err.Error(), "find pull request status", "Error message should mention status lookup")
}
