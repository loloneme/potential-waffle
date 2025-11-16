//go:generate mockgen -source=contract.go -destination=mocks/contract.go -package=mocks

package bulk_deactivate_team

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/pull_request"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/user"
)

type userRepo interface {
	BulkDeactivateUsers(ctx context.Context, tx *sqlx.Tx, teamName string, userIDs []string) ([]string, error)
	Find(ctx context.Context, spec user.FindSpecification) ([]models.User, error)
}

type prRepo interface {
	WithTx(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) error
	GetOpenPRsWithReviewers(ctx context.Context, reviewerIDs []string) (map[string][]string, error)
	GetOpenPRsWithFullInfo(ctx context.Context, deactivatedReviewerIDs []string) (map[string]pull_request.PRFullInfo, error)
	BulkReassignReviewers(ctx context.Context, tx *sqlx.Tx, reassignments []pull_request.PRReassignments) error
}
