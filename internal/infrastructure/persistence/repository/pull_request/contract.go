//go:generate mockgen -source=contract.go -destination=mocks/contract.go -package=mocks

package pull_request

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

type pullRequestRepository interface {
	WithTx(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) error

	InsertPullRequest(ctx context.Context, tx *sqlx.Tx, pr *models.PullRequest) (models.PullRequest, error)
	SetPullRequestStatus(ctx context.Context, prID string, statusName string) error
	UpdatePullRequest(ctx context.Context, spec UpdateSpecification) error
	GetPRByID(ctx context.Context, prID string) (models.PullRequest, error)
	PullRequestExists(ctx context.Context, prID string) (bool, error)

	FindStatus(ctx context.Context, spec FindSpecification) (*models.Status, error)

	InsertReviewers(ctx context.Context, tx *sqlx.Tx, prID string, reviewers []string) error
	ReassignReviewer(ctx context.Context, tx *sqlx.Tx, prID, oldReviewerID, newReviewerID string) error
	GetAvailableReviewers(ctx context.Context, teamName string, excludeIDs []string, limit int) ([]string, error)
	GetPullRequestReviewers(ctx context.Context, prID string) ([]string, error)
}
