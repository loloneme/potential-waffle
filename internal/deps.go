package internal

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/postgres"
)

func NewDatabaseConnection(ctx context.Context) (*sqlx.DB, error) {
	return postgres.NewFromConfig(ctx)
}
