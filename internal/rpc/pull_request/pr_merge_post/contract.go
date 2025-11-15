//go:generate mockgen -source=contract.go -destination=mocks/contract.go -package=mocks

package pr_merge_post

import (
	"context"

	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

type MergePRService interface {
	MergePullRequest(ctx context.Context, prID string, mergeStatus string) (models.PullRequest, error)
}
