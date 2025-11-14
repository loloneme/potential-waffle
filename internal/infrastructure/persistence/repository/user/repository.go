package user

import (
	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence"
)

const (
	tableName = "users"
	alias     = "u"
	idField   = "user_id"
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
