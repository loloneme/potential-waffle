//go:generate mockgen -source=contract.go -destination=mocks/contract.go -package=mocks

package users_get_review_get

import (
	"context"

	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/pull_request"
)

type prRepo interface {
	FindPullRequests(ctx context.Context, spec pull_request.FindSpecification) ([]models.PullRequest, error)
}
