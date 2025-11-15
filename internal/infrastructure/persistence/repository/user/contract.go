package user

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

type userRepository interface {
	UpsertUsers(ctx context.Context, tx *sqlx.Tx, users []models.User) ([]models.User, error)
	GetUserByID(ctx context.Context, userID string) (models.User, error)
	Find(ctx context.Context, spec FindSpecification) ([]models.User, error)
	UserUpdate(ctx context.Context, spec UpdateSpecification) (models.User, error)

	GetUserTeamName(ctx context.Context, userID string) (string, error)
}
