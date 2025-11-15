package team

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence"
)

const (
	tableName = "teams"
	alias     = "t"
	idField   = "team_name"
)

type Repository struct {
	db        *sqlx.DB
	tableName string
	columns   *persistence.Columns
}

func NewRepository(db *sqlx.DB) *Repository {
	cols := persistence.NewColumns(readableColumns, writableColumns, alias, idField)

	return &Repository{
		db:        db,
		tableName: tableName,
		columns:   cols,
	}
}

func (r *Repository) WithTx(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}
