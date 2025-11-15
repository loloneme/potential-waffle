package reassign_pr_test

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
	"github.com/loloneme/potential-waffle/internal/usecase/reassign_pr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testEnv struct {
	ctx      context.Context
	db       *sqlx.DB
	userRepo *user.Repository
	prRepo   *pull_request.Repository
	teamRepo *team.Repository
	service  *reassign_pr.Service
}

func setupTest(t *testing.T) *testEnv {
	t.Helper()

	ctx := context.Background()

	db, err := test_env.NewTestDatabaseConnection(ctx)
	require.NoError(t, err)

	userRepo := user.NewRepository(db)
	prRepo := pull_request.NewRepository(db)
	teamRepo := team.NewRepository(db)
	service := reassign_pr.New(userRepo, prRepo)

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

func TestService_ReassignReviewer_Successful(t *testing.T) {
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

		// Назначаем только одного ревьюера (согласно numberOfReviewers = 1)
		reviewers := []string{reviewer1ID}
		return env.prRepo.InsertReviewers(ctx, tx, prID, reviewers)
	})
	require.NoError(t, err)

	// Переназначаем ревьюера
	reassignedPR, newReviewerID, err := env.service.ReassignReviewer(env.ctx, prID, reviewer1ID)
	require.NoError(t, err)

	// Проверяем результат
	assert.Equal(t, prID, reassignedPR.ID)
	assert.NotEqual(t, reviewer1ID, newReviewerID)
	assert.Contains(t, []string{reviewer2ID, reviewer3ID}, newReviewerID)
	assert.NotContains(t, reassignedPR.Reviewers, reviewer1ID)
	assert.Contains(t, reassignedPR.Reviewers, newReviewerID)

	// Проверяем, что ревьюер действительно переназначен в БД
	reviewers, err := env.prRepo.GetPullRequestReviewers(env.ctx, prID)
	require.NoError(t, err)
	assert.Len(t, reviewers, 1)
	assert.NotContains(t, reviewers, reviewer1ID)
	assert.Contains(t, reviewers, newReviewerID)
}

func TestService_ReassignReviewer_PRNotFound(t *testing.T) {
	env := setupTest(t)

	nonExistentPRID := "non-existent-pr"
	reviewerID := "reviewer-1"

	_, _, err := env.service.ReassignReviewer(env.ctx, nonExistentPRID, reviewerID)
	require.Error(t, err)
	var notFoundErr *rpc_errors.NotFoundError
	assert.ErrorAs(t, err, &notFoundErr)
}

func TestService_ReassignReviewer_PRAlreadyMerged(t *testing.T) {
	env := setupTest(t)

	// Подготовка данных
	teamName := "test-team-2"
	err := env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := env.teamRepo.CreateTeam(ctx, tx, models.Team{TeamName: teamName})
		return err
	})
	require.NoError(t, err)

	authorID := "author-2"
	reviewer1ID := "reviewer-4"

	err = env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		users := []models.User{
			{ID: authorID, Username: "author2", IsActive: true, TeamName: teamName},
			{ID: reviewer1ID, Username: "reviewer4", IsActive: true, TeamName: teamName},
		}
		_, err := env.userRepo.UpsertUsers(ctx, tx, users)
		return err
	})
	require.NoError(t, err)

	// Создаем PR со статусом MERGED
	prID := "pr-2"
	err = env.prRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		pr := &models.PullRequest{
			ID:       prID,
			Name:     "Merged PR",
			AuthorID: authorID,
			Status:   &models.Status{Name: "MERGED"},
		}

		foundStatus, err := env.prRepo.FindStatus(ctx, status.NewGetStatusByNameSpecification("MERGED"))
		require.NoError(t, err)
		pr.StatusID = foundStatus.ID

		_, err = env.prRepo.InsertPullRequest(ctx, tx, pr)
		if err != nil {
			return err
		}

		reviewers := []string{reviewer1ID}
		return env.prRepo.InsertReviewers(ctx, tx, prID, reviewers)
	})
	require.NoError(t, err)

	// Пытаемся переназначить ревьюера для мерженного PR
	_, _, err = env.service.ReassignReviewer(env.ctx, prID, reviewer1ID)
	require.Error(t, err)
	var mergedErr *rpc_errors.PRMergedError
	assert.ErrorAs(t, err, &mergedErr)
}

func TestService_ReassignReviewer_ReviewerNotFound(t *testing.T) {
	env := setupTest(t)

	// Подготовка данных
	teamName := "test-team-3"
	err := env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := env.teamRepo.CreateTeam(ctx, tx, models.Team{TeamName: teamName})
		return err
	})
	require.NoError(t, err)

	authorID := "author-3"
	err = env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		users := []models.User{
			{ID: authorID, Username: "author3", IsActive: true, TeamName: teamName},
		}
		_, err := env.userRepo.UpsertUsers(ctx, tx, users)
		return err
	})
	require.NoError(t, err)

	// Создаем PR
	prID := "pr-3"
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
		return err
	})
	require.NoError(t, err)

	// Пытаемся переназначить несуществующего ревьюера
	nonExistentReviewerID := "non-existent-reviewer"
	_, _, err = env.service.ReassignReviewer(env.ctx, prID, nonExistentReviewerID)
	require.Error(t, err)
	var notFoundErr *rpc_errors.NotFoundError
	assert.ErrorAs(t, err, &notFoundErr)
}

func TestService_ReassignReviewer_ReviewerNotAssigned(t *testing.T) {
	env := setupTest(t)

	// Подготовка данных
	teamName := "test-team-4"
	err := env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := env.teamRepo.CreateTeam(ctx, tx, models.Team{TeamName: teamName})
		return err
	})
	require.NoError(t, err)

	authorID := "author-4"
	reviewer1ID := "reviewer-5"
	reviewer2ID := "reviewer-6"

	err = env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		users := []models.User{
			{ID: authorID, Username: "author4", IsActive: true, TeamName: teamName},
			{ID: reviewer1ID, Username: "reviewer5", IsActive: true, TeamName: teamName},
			{ID: reviewer2ID, Username: "reviewer6", IsActive: true, TeamName: teamName},
		}
		_, err := env.userRepo.UpsertUsers(ctx, tx, users)
		return err
	})
	require.NoError(t, err)

	// Создаем PR с одним ревьюером
	prID := "pr-4"
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

		// Назначаем только reviewer1ID
		reviewers := []string{reviewer1ID}
		return env.prRepo.InsertReviewers(ctx, tx, prID, reviewers)
	})
	require.NoError(t, err)

	// Пытаемся переназначить reviewer2ID, который не назначен на PR
	_, _, err = env.service.ReassignReviewer(env.ctx, prID, reviewer2ID)
	require.Error(t, err)
	var notAssignedErr *rpc_errors.NotAssignedError
	assert.ErrorAs(t, err, &notAssignedErr)
}

func TestService_ReassignReviewer_NoAvailableCandidates(t *testing.T) {
	env := setupTest(t)

	// Подготовка данных: команда с только одним активным пользователем (кроме автора)
	teamName := "lonely-team"
	err := env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := env.teamRepo.CreateTeam(ctx, tx, models.Team{TeamName: teamName})
		return err
	})
	require.NoError(t, err)

	authorID := "author-5"
	reviewer1ID := "reviewer-7"

	err = env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		users := []models.User{
			{ID: authorID, Username: "author5", IsActive: true, TeamName: teamName},
			{ID: reviewer1ID, Username: "reviewer7", IsActive: true, TeamName: teamName},
		}
		_, err := env.userRepo.UpsertUsers(ctx, tx, users)
		return err
	})
	require.NoError(t, err)

	// Создаем PR
	prID := "pr-5"
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

		reviewers := []string{reviewer1ID}
		return env.prRepo.InsertReviewers(ctx, tx, prID, reviewers)
	})
	require.NoError(t, err)

	// Пытаемся переназначить ревьюера, но нет других доступных кандидатов
	// (только автор и текущий ревьюер в команде)
	_, _, err = env.service.ReassignReviewer(env.ctx, prID, reviewer1ID)
	require.Error(t, err)
	var noCandidateErr *rpc_errors.NoCandidateError
	assert.ErrorAs(t, err, &noCandidateErr)
}

func TestService_ReassignReviewer_ExcludesInactiveUsers(t *testing.T) {
	env := setupTest(t)

	// Подготовка данных: команда с активными и неактивными пользователями
	teamName := "mixed-team"
	err := env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := env.teamRepo.CreateTeam(ctx, tx, models.Team{TeamName: teamName})
		return err
	})
	require.NoError(t, err)

	authorID := "author-6"
	reviewer1ID := "reviewer-8"
	activeReviewerID := "active-reviewer"
	inactiveReviewerID := "inactive-reviewer"

	err = env.teamRepo.WithTx(env.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		users := []models.User{
			{ID: authorID, Username: "author6", IsActive: true, TeamName: teamName},
			{ID: reviewer1ID, Username: "reviewer8", IsActive: true, TeamName: teamName},
			{ID: activeReviewerID, Username: "active", IsActive: true, TeamName: teamName},
			{ID: inactiveReviewerID, Username: "inactive", IsActive: false, TeamName: teamName},
		}
		_, err := env.userRepo.UpsertUsers(ctx, tx, users)
		return err
	})
	require.NoError(t, err)

	// Создаем PR
	prID := "pr-6"
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

		reviewers := []string{reviewer1ID}
		return env.prRepo.InsertReviewers(ctx, tx, prID, reviewers)
	})
	require.NoError(t, err)

	// Переназначаем ревьюера
	reassignedPR, newReviewerID, err := env.service.ReassignReviewer(env.ctx, prID, reviewer1ID)
	require.NoError(t, err)

	// Проверяем, что назначен только активный ревьюер
	assert.Equal(t, activeReviewerID, newReviewerID)
	assert.Contains(t, reassignedPR.Reviewers, activeReviewerID)
	assert.NotContains(t, reassignedPR.Reviewers, inactiveReviewerID)
}
