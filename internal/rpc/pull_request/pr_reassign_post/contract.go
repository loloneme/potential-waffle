package pr_reassign_post

import (
	"context"

	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

type reassignPRService interface {
	ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (models.PullRequest, string, error)
}
