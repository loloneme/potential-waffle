package team

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

var st = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var (
	ErrNotFound = errors.New("team not found")
)

func getValues(team models.Team) []interface{} {
	return []interface{}{
		team.TeamName,
	}
}

func (r *Repository) CreateTeam(ctx context.Context, team models.Team) (models.Team, error) {
	var res models.Team

	sqlStr, params, err := st.
		Insert(r.tableName).
		Columns(r.columns.ForInsert()...).
		Values(getValues(team)).
		Suffix("RETURNING *").
		ToSql()

	if err != nil {
		return res, fmt.Errorf("create team: %w", err)
	}

	if err = r.db.GetContext(ctx, &res, sqlStr, params...); err != nil {
		return res, fmt.Errorf("get team: %w", err)
	}
	return res, nil
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
