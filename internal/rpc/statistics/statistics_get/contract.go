//go:generate mockgen -source=contract.go -destination=mocks/contract.go -package=mocks

package statistics_get

import (
	"context"

	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/pull_request"
)

type prRepo interface {
	GetStatistics(ctx context.Context) (*pull_request.Statistics, error)
}
