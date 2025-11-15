package pull_request_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/pull_request"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/pull_request/mocks"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/specification/status"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPullRequestRepository_GetPRByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockpullRequestRepository(ctrl)
	ctx := context.Background()
	prID := "pr-1"

	t.Run("successful get", func(t *testing.T) {
		expectedPR := models.PullRequest{
			ID:       "pr-1",
			Name:     "Test PR",
			AuthorID: "user-1",
			StatusID: 1,
		}

		mockRepo.EXPECT().
			GetPRByID(gomock.Any(), prID).
			Return(expectedPR, nil)

		pr, err := mockRepo.GetPRByID(ctx, prID)
		assert.NoError(t, err)
		assert.Equal(t, expectedPR.ID, pr.ID)
		assert.Equal(t, expectedPR.Name, pr.Name)
	})

	t.Run("PR not found", func(t *testing.T) {
		mockRepo.EXPECT().
			GetPRByID(gomock.Any(), prID).
			Return(models.PullRequest{}, pull_request.ErrPRNotFound)

		pr, err := mockRepo.GetPRByID(ctx, prID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, pull_request.ErrPRNotFound))
		assert.Equal(t, models.PullRequest{}, pr)
	})
}

func TestPullRequestRepository_InsertPullRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockpullRequestRepository(ctrl)
	ctx := context.Background()
	pr := &models.PullRequest{
		ID:       "pr-1",
		Name:     "Test PR",
		AuthorID: "user-1",
		StatusID: 1,
	}
	var tx *sqlx.Tx

	t.Run("successful insert", func(t *testing.T) {
		expectedPR := models.PullRequest{
			ID:       "pr-1",
			Name:     "Test PR",
			AuthorID: "user-1",
			StatusID: 1,
		}

		mockRepo.EXPECT().
			InsertPullRequest(gomock.Any(), gomock.Any(), pr).
			Return(expectedPR, nil)

		result, err := mockRepo.InsertPullRequest(ctx, tx, pr)
		assert.NoError(t, err)
		assert.Equal(t, expectedPR.ID, result.ID)
	})

	t.Run("PR already exists", func(t *testing.T) {
		mockRepo.EXPECT().
			InsertPullRequest(gomock.Any(), gomock.Any(), pr).
			Return(models.PullRequest{}, pull_request.ErrPRAlreadyExists)

		result, err := mockRepo.InsertPullRequest(ctx, tx, pr)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, pull_request.ErrPRAlreadyExists))
		assert.Equal(t, models.PullRequest{}, result)
	})
}

func TestPullRequestRepository_PullRequestExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockpullRequestRepository(ctrl)
	ctx := context.Background()
	prID := "pr-1"

	t.Run("PR exists", func(t *testing.T) {
		mockRepo.EXPECT().
			PullRequestExists(gomock.Any(), prID).
			Return(true, nil)

		exists, err := mockRepo.PullRequestExists(ctx, prID)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("PR does not exist", func(t *testing.T) {
		mockRepo.EXPECT().
			PullRequestExists(gomock.Any(), prID).
			Return(false, nil)

		exists, err := mockRepo.PullRequestExists(ctx, prID)
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestPullRequestRepository_SetPullRequestStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockpullRequestRepository(ctrl)
	ctx := context.Background()
	prID := "pr-1"
	statusName := "MERGED"

	t.Run("successful status update", func(t *testing.T) {
		mockRepo.EXPECT().
			SetPullRequestStatus(gomock.Any(), prID, statusName).
			Return(nil)

		err := mockRepo.SetPullRequestStatus(ctx, prID, statusName)
		assert.NoError(t, err)
	})

	t.Run("status not found", func(t *testing.T) {
		mockRepo.EXPECT().
			SetPullRequestStatus(gomock.Any(), prID, "INVALID").
			Return(pull_request.ErrStatusNotFound)

		err := mockRepo.SetPullRequestStatus(ctx, prID, "INVALID")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, pull_request.ErrStatusNotFound))
	})
}

func TestPullRequestRepository_FindStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockpullRequestRepository(ctrl)
	ctx := context.Background()

	t.Run("successful find", func(t *testing.T) {
		expectedStatus := &models.Status{
			ID:   1,
			Name: "OPEN",
		}
		spec := status.NewGetStatusByNameSpecification("OPEN")

		mockRepo.EXPECT().
			FindStatus(gomock.Any(), spec).
			Return(expectedStatus, nil)

		status, err := mockRepo.FindStatus(ctx, spec)
		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.Equal(t, expectedStatus.ID, status.ID)
		assert.Equal(t, expectedStatus.Name, status.Name)
	})

	t.Run("status not found", func(t *testing.T) {
		spec := status.NewGetStatusByNameSpecification("INVALID")

		mockRepo.EXPECT().
			FindStatus(gomock.Any(), spec).
			Return(nil, pull_request.ErrStatusNotFound)

		status, err := mockRepo.FindStatus(ctx, spec)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, pull_request.ErrStatusNotFound))
		assert.Nil(t, status)
	})
}

func TestPullRequestRepository_InsertReviewers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockpullRequestRepository(ctrl)
	ctx := context.Background()
	prID := "pr-1"
	reviewers := []string{"reviewer-1", "reviewer-2"}
	var tx *sqlx.Tx

	t.Run("successful insert", func(t *testing.T) {
		mockRepo.EXPECT().
			InsertReviewers(gomock.Any(), gomock.Any(), prID, reviewers).
			Return(nil)

		err := mockRepo.InsertReviewers(ctx, tx, prID, reviewers)
		assert.NoError(t, err)
	})

	t.Run("empty reviewers list", func(t *testing.T) {
		mockRepo.EXPECT().
			InsertReviewers(gomock.Any(), gomock.Any(), prID, []string{}).
			Return(nil)

		err := mockRepo.InsertReviewers(ctx, tx, prID, []string{})
		assert.NoError(t, err)
	})
}

func TestPullRequestRepository_GetPullRequestReviewers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockpullRequestRepository(ctrl)
	ctx := context.Background()
	prID := "pr-1"

	t.Run("successful get reviewers", func(t *testing.T) {
		expectedReviewers := []string{"reviewer-1", "reviewer-2"}

		mockRepo.EXPECT().
			GetPullRequestReviewers(gomock.Any(), prID).
			Return(expectedReviewers, nil)

		reviewers, err := mockRepo.GetPullRequestReviewers(ctx, prID)
		assert.NoError(t, err)
		assert.Equal(t, expectedReviewers, reviewers)
	})

	t.Run("no reviewers", func(t *testing.T) {
		mockRepo.EXPECT().
			GetPullRequestReviewers(gomock.Any(), prID).
			Return([]string{}, nil)

		reviewers, err := mockRepo.GetPullRequestReviewers(ctx, prID)
		assert.NoError(t, err)
		assert.Empty(t, reviewers)
	})
}

func TestPullRequestRepository_GetAvailableReviewers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockpullRequestRepository(ctrl)
	ctx := context.Background()
	teamName := "team-1"
	excludeIDs := []string{"user-1"}
	limit := 2

	t.Run("successful get available reviewers", func(t *testing.T) {
		expectedReviewers := []string{"user-2", "user-3"}

		mockRepo.EXPECT().
			GetAvailableReviewers(gomock.Any(), teamName, excludeIDs, limit).
			Return(expectedReviewers, nil)

		reviewers, err := mockRepo.GetAvailableReviewers(ctx, teamName, excludeIDs, limit)
		assert.NoError(t, err)
		assert.Equal(t, expectedReviewers, reviewers)
	})

	t.Run("no available reviewers", func(t *testing.T) {
		mockRepo.EXPECT().
			GetAvailableReviewers(gomock.Any(), teamName, excludeIDs, limit).
			Return(nil, pull_request.ErrReviewersNotFound)

		reviewers, err := mockRepo.GetAvailableReviewers(ctx, teamName, excludeIDs, limit)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, pull_request.ErrReviewersNotFound))
		assert.Nil(t, reviewers)
	})
}

func TestPullRequestRepository_ReassignReviewer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockpullRequestRepository(ctrl)
	ctx := context.Background()
	prID := "pr-1"
	oldReviewerID := "reviewer-1"
	newReviewerID := "reviewer-2"
	var tx *sqlx.Tx

	t.Run("successful reassign", func(t *testing.T) {
		mockRepo.EXPECT().
			ReassignReviewer(gomock.Any(), gomock.Any(), prID, oldReviewerID, newReviewerID).
			Return(nil)

		err := mockRepo.ReassignReviewer(ctx, tx, prID, oldReviewerID, newReviewerID)
		assert.NoError(t, err)
	})

	t.Run("reviewer not assigned", func(t *testing.T) {
		mockRepo.EXPECT().
			ReassignReviewer(gomock.Any(), gomock.Any(), prID, "reviewer-999", newReviewerID).
			Return(pull_request.ErrReviewerNotAssigned)

		err := mockRepo.ReassignReviewer(ctx, tx, prID, "reviewer-999", newReviewerID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, pull_request.ErrReviewerNotAssigned))
	})
}

func TestPullRequestRepository_WithTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockpullRequestRepository(ctrl)
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
