package merge_pr

import (
	"context"

	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

type prRepo interface {
	PullRequestExists(ctx context.Context, prID string) (bool, error)
	SetPullRequestStatus(ctx context.Context, prID string, statusName string) error

	GetPRByID(ctx context.Context, prID string) (models.PullRequest, error)
}
