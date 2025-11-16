//go:generate mockgen -source=contract.go -destination=mocks/contract.go -package=mocks

package users_bulk_deactivate_post

import (
	"context"

	"github.com/loloneme/potential-waffle/internal/usecase/bulk_deactivate_team"
)

type bulkDeactivateService interface {
	BulkDeactivateTeamUsers(ctx context.Context, teamName string, userIDs []string) (bulk_deactivate_team.BulkDeactivateResult, error)
}
