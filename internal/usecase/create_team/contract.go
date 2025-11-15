package create_team

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

type userRepo interface {
	UpsertUsers(ctx context.Context, tx *sqlx.Tx, users []models.User) ([]models.User, error)
}

type teamRepo interface {
	WithTx(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) error

	CreateTeam(ctx context.Context, tx *sqlx.Tx, team models.Team) (models.Team, error)
}
