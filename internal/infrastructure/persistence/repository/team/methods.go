package team

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

var st = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var (
	ErrNotFound = errors.New("team not found")

	ErrAlreadyExists = errors.New("team already exists")
)

func getValues(team models.Team) []interface{} {
	return []interface{}{
		team.TeamName,
	}
}

func (r *Repository) CreateTeam(ctx context.Context, tx *sqlx.Tx, team models.Team) (models.Team, error) {
	if exists, err := r.Exists(ctx, team.TeamName); err != nil {
		return models.Team{}, err
	} else if exists {
		return team, nil
	}
	var res models.Team

	sqlStr, params, err := st.
		Insert(r.tableName).
		Columns(r.columns.ForInsert()...).
		Values(getValues(team)...).
		Suffix("RETURNING *").
		ToSql()

	if err != nil {
		return res, fmt.Errorf("create team: %w", err)
	}

	if err = tx.GetContext(ctx, &res, sqlStr, params...); err != nil {
		return res, fmt.Errorf("get team: %w", err)
	}
	return res, nil
}

func (r *Repository) Exists(ctx context.Context, teamName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`

	if err := r.db.GetContext(ctx, &exists, query, teamName); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repository) FindTeamByID(ctx context.Context, teamName string) (models.Team, error) {
	var res models.Team

	sqlStr, params, err := st.
		Select(r.tableName).
		From(r.tableName).
		Where(sq.Eq{r.tableName: teamName}).
		ToSql()

	if err != nil {
		return res, fmt.Errorf("create team: %w", err)
	}

	err = r.db.GetContext(ctx, &res, sqlStr, params...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, ErrNotFound
		}
		return res, fmt.Errorf("exec get team: %w", err)
	}

	return res, nil
}
