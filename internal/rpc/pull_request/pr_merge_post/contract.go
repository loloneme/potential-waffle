package pr_merge_post

import (
	"context"

	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

type MergePRService interface {
	MergePullRequest(ctx context.Context, prID string, mergeStatus string) (models.PullRequest, error)
}
