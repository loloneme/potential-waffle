//go:generate mockgen -source=contract.go -destination=mocks/contract.go -package=mocks

package team

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

type teamRepository interface {
	WithTx(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) error

	CreateTeam(ctx context.Context, tx *sqlx.Tx, team models.Team) (models.Team, error)
	FindTeamByID(ctx context.Context, teamName string) (models.Team, error)
}
