package pr_create_post

import (
	"context"

	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

type createPRService interface {
	CreatePR(ctx context.Context, pr *models.PullRequest) (models.PullRequest, error)
}
