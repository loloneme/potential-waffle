package reassign_pr

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

type userRepo interface {
	GetUserTeamName(ctx context.Context, userID string) (string, error)
}

type prRepo interface {
	WithTx(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) error

	GetPRByID(ctx context.Context, prID string) (models.PullRequest, error)
	GetAvailableReviewers(ctx context.Context, teamName string, excludeIDs []string, limit int) ([]string, error)
	GetPullRequestReviewers(ctx context.Context, prID string) ([]string, error)
	ReassignReviewer(ctx context.Context, tx *sqlx.Tx, prID, oldReviewerID, newReviewerID string) error
}
