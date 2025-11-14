package pull_request

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

type pullRequestRepository interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error

	InsertPullRequest(ctx context.Context, tx *sqlx.Tx, pr *models.PullRequest) (models.PullRequest, error)
	MergePullRequest(ctx context.Context, prID string, statusID int64) (models.PullRequest, error)
	UpdatePullRequest(ctx context.Context, spec UpdateSpecification) (models.PullRequest, error)
	GetPRByID(ctx context.Context, prID string) (models.PullRequest, error)

	FindStatus(ctx context.Context, spec FindSpecification) (*models.Status, error)

	InsertReviewers(ctx context.Context, tx *sqlx.Tx, prID string, reviewers []string) error
	FindReviewers(ctx context.Context, spec FindSpecification) ([]string, error)
	ReassignReviewer(ctx context.Context, tx *sqlx.Tx, prID, oldReviewerID, newReviewerID string) error
}
