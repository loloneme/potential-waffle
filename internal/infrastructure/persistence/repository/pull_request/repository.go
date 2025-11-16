package pull_request

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence"
)

const (
	tableName = "pull_requests"
	alias     = "pr"
	idField   = "pr_id"
)

const (
	reviewersTableName = "reviewers"
	statusTableName    = "statuses"
	usersTableName     = "users"
)

type Repository struct {
	db *sqlx.DB

	tableName          string
	reviewersTableName string
	statusTableName    string
	usersTableName     string

	pullRequestColumns *persistence.Columns
	reviewerColumns    *persistence.Columns
	statusColumns      *persistence.Columns
}

func NewRepository(db *sqlx.DB) *Repository {
	prCols := persistence.NewColumns(readableColumns, writableColumns, alias, idField)
	rCols := persistence.NewColumns(
		[]string{"pr_id", "reviewer_id"},
		[]string{"pr_id", "reviewer_id"},
		"r",
		"",
	)
	sCols := persistence.NewColumns(
		[]string{"status_id", "status_name"},
		[]string{"status_name"},
		"s",
		"status_id",
	)

	return &Repository{
		db: db,

		tableName:          tableName,
		reviewersTableName: reviewersTableName,
		statusTableName:    statusTableName,
		usersTableName:     usersTableName,

		pullRequestColumns: prCols,
		reviewerColumns:    rCols,
		statusColumns:      sCols,
	}
}

func (r *Repository) WithTx(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("rollback transaction: %v", err)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}
