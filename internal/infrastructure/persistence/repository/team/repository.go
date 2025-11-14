package team

import (
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
