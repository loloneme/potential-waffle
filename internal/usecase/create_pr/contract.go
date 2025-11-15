package create_pr

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/pull_request"
)

type userRepo interface {
	GetUserTeamName(ctx context.Context, userID string) (string, error)
}

type prRepo interface {
	WithTx(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) error

	InsertPullRequest(ctx context.Context, tx *sqlx.Tx, pr *models.PullRequest) (models.PullRequest, error)
	FindStatus(ctx context.Context, spec pull_request.FindSpecification) (*models.Status, error)

	GetAvailableReviewers(ctx context.Context, teamName string, excludeIDs []string, limit int) ([]string, error)

	InsertReviewers(ctx context.Context, tx *sqlx.Tx, prID string, reviewers []string) error
}
