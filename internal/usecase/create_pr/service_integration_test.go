package create_pr_test

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/pull_request"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/team"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/user"
	"github.com/loloneme/potential-waffle/internal/infrastructure/utils/test_env"
	rpc_errors "github.com/loloneme/potential-waffle/internal/rpc/errors"
	"github.com/loloneme/potential-waffle/internal/usecase/create_pr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testEnv struct {
	ctx context.Context
	db  *sqlx.DB

	userRepo *user.Repository
	prRepo   *pull_request.Repository
	teamRepo *team.Repository
	service  *create_pr.Service
}

func setupTest(t *testing.T) *testEnv {
	t.Helper()

	ctx := context.Background()

	db, err := test_env.NewTestDatabaseConnection(ctx)
	require.NoError(t, err)

	userRepo := user.NewRepository(db)
	prRepo := pull_request.NewRepository(db)
	teamRepo := team.NewRepository(db)
	service := create_pr.New(userRepo, prRepo)

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

func TestService_CreatePR_Successful(t *testing.T) {
	env := setupTest(t)

	teamName := "test-team"
	err := env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := env.teamRepo.CreateTeam(ctx, tx, models.Team{TeamName: teamName})
		return err
	})
	require.NoError(t, err)

	authorID := "author-1"
	reviewer1ID := "reviewer-1"
	reviewer2ID := "reviewer-2"
	reviewer3ID := "reviewer-3"

	err = env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		users := []models.User{
			{ID: authorID, Username: "author", IsActive: true, TeamName: teamName},
			{ID: reviewer1ID, Username: "reviewer1", IsActive: true, TeamName: teamName},
			{ID: reviewer2ID, Username: "reviewer2", IsActive: true, TeamName: teamName},
			{ID: reviewer3ID, Username: "reviewer3", IsActive: true, TeamName: teamName},
		}
		_, err := env.userRepo.UpsertUsers(ctx, tx, users)
		return err
	})
	require.NoError(t, err)

	// Создаем PR
	pr := &models.PullRequest{
		ID:       "pr-1",
		Name:     "Test PR",
		AuthorID: authorID,
		Status:   &models.Status{Name: "OPEN"},
	}

	createdPR, err := env.service.CreatePR(env.ctx, pr)
	require.NoError(t, err)

	// Проверяем результат
	assert.Equal(t, "pr-1", createdPR.ID)
	assert.Equal(t, "Test PR", createdPR.Name)
	assert.Equal(t, authorID, createdPR.AuthorID)
	assert.NotNil(t, createdPR.Status)
	assert.Equal(t, "OPEN", createdPR.Status.Name)
	assert.NotZero(t, createdPR.StatusID)
	assert.Len(t, createdPR.Reviewers, 2)
	assert.NotContains(t, createdPR.Reviewers, authorID) // Автор не должен быть ревьюером

	// Проверяем, что ревьюеры действительно назначены в БД
	reviewers, err := env.prRepo.GetPullRequestReviewers(env.ctx, pr.ID)
	require.NoError(t, err)
	assert.Len(t, reviewers, 2)
	assert.NotContains(t, reviewers, authorID)
}

func TestService_CreatePR_AuthorNotFound(t *testing.T) {
	env := setupTest(t)

	pr := &models.PullRequest{
		ID:       "pr-2",
		Name:     "Test PR",
		AuthorID: "non-existent-author",
		Status:   &models.Status{Name: "OPEN"},
	}

	_, err := env.service.CreatePR(env.ctx, pr)
	require.Error(t, err)
	var notFoundErr *rpc_errors.NotFoundError
	assert.ErrorAs(t, err, &notFoundErr)
}

func TestService_CreatePR_NoAvailableReviewers(t *testing.T) {
	env := setupTest(t)

	// Создаем команду с только одним пользователем (автором)
	teamName := "lonely-team"
	err := env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := env.teamRepo.CreateTeam(ctx, tx, models.Team{TeamName: teamName})
		return err
	})
	require.NoError(t, err)

	authorID := "lonely-author"
	err = env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		users := []models.User{
			{ID: authorID, Username: "lonely", IsActive: true, TeamName: teamName},
		}
		_, err := env.userRepo.UpsertUsers(ctx, tx, users)
		return err
	})
	require.NoError(t, err)

	pr := &models.PullRequest{
		ID:       "pr-3",
		Name:     "Test PR",
		AuthorID: authorID,
		Status:   &models.Status{Name: "OPEN"},
	}

	_, err = env.service.CreatePR(env.ctx, pr)
	require.Error(t, err)
	var notFoundErr *rpc_errors.NotFoundError
	assert.ErrorAs(t, err, &notFoundErr)
}

func TestService_CreatePR_AlreadyExists(t *testing.T) {
	env := setupTest(t)

	// Подготовка данных
	teamName := "team-2"
	err := env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := env.teamRepo.CreateTeam(ctx, tx, models.Team{TeamName: teamName})
		return err
	})
	require.NoError(t, err)

	authorID := "author-2"
	reviewer1ID := "reviewer-4"
	reviewer2ID := "reviewer-5"

	err = env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		users := []models.User{
			{ID: authorID, Username: "author2", IsActive: true, TeamName: teamName},
			{ID: reviewer1ID, Username: "reviewer4", IsActive: true, TeamName: teamName},
			{ID: reviewer2ID, Username: "reviewer5", IsActive: true, TeamName: teamName},
		}
		_, err := env.userRepo.UpsertUsers(ctx, tx, users)
		return err
	})
	require.NoError(t, err)

	// Создаем первый PR
	pr1 := &models.PullRequest{
		ID:       "pr-4",
		Name:     "First PR",
		AuthorID: authorID,
		Status:   &models.Status{Name: "OPEN"},
	}
	_, err = env.service.CreatePR(env.ctx, pr1)
	require.NoError(t, err)

	// Пытаемся создать PR с тем же ID
	pr2 := &models.PullRequest{
		ID:       "pr-4",
		Name:     "Duplicate PR",
		AuthorID: authorID,
		Status:   &models.Status{Name: "OPEN"},
	}
	_, err = env.service.CreatePR(env.ctx, pr2)
	require.Error(t, err)
	var prExistsErr *rpc_errors.PRExistsError
	assert.ErrorAs(t, err, &prExistsErr)
}

func TestService_CreatePR_ReviewersAreActiveUsersOnly(t *testing.T) {
	env := setupTest(t)

	// Подготовка данных: команда с активными и неактивными пользователями
	teamName := "mixed-team"
	err := env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := env.teamRepo.CreateTeam(ctx, tx, models.Team{TeamName: teamName})
		return err
	})
	require.NoError(t, err)

	authorID := "author-3"
	activeReviewerID := "active-reviewer"
	inactiveReviewerID := "inactive-reviewer"

	err = env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		users := []models.User{
			{ID: authorID, Username: "author3", IsActive: true, TeamName: teamName},
			{ID: activeReviewerID, Username: "active", IsActive: true, TeamName: teamName},
			{ID: inactiveReviewerID, Username: "inactive", IsActive: false, TeamName: teamName},
		}
		_, err := env.userRepo.UpsertUsers(ctx, tx, users)
		return err
	})
	require.NoError(t, err)

	pr := &models.PullRequest{
		ID:       "pr-5",
		Name:     "Test PR",
		AuthorID: authorID,
		Status:   &models.Status{Name: "OPEN"},
	}

	createdPR, err := env.service.CreatePR(env.ctx, pr)
	require.NoError(t, err)

	// Проверяем, что назначен только активный ревьюер
	assert.Len(t, createdPR.Reviewers, 1)
	assert.Contains(t, createdPR.Reviewers, activeReviewerID)
	assert.NotContains(t, createdPR.Reviewers, inactiveReviewerID)
}
