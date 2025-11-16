package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	user_spec "github.com/loloneme/potential-waffle/internal/infrastructure/persistence/specification/user"
)

var st = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var ErrNotFound = errors.New("user not found")

type FindSpecification interface {
	GetFields() []string
	GetRule(s sq.SelectBuilder) sq.SelectBuilder
}

type UpdateSpecification interface {
	GetSetValues() map[string]interface{}
	GetRule(builder sq.UpdateBuilder) sq.UpdateBuilder
	GetReturningFields() []string
}

func GetValues(user models.User) []interface{} {
	return []interface{}{
		user.ID,
		user.Username,
		user.IsActive,
		user.TeamName,
	}
}

func (r *Repository) GetUserByID(ctx context.Context, userID string) (models.User, error) {
	var user models.User

	sqlStr, params, err := st.
		Select(r.columns.ForSelect(nil)...).
		From(r.tableName).
		Where(sq.Eq{r.columns.GetIDField(): userID}).
		ToSql()
	if err != nil {
		return user, err
	}

	err = r.db.GetContext(ctx, &user, sqlStr, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrNotFound
		}
		return user, err
	}

	return user, nil
}

func (r *Repository) GetUserTeamName(ctx context.Context, userID string) (string, error) {
	var res string

	queryBuilder := st.
		Select("team_name").From(r.tableName).
		Where(sq.Eq{"user_id": userID})

	sqlStr, params, err := queryBuilder.ToSql()
	if err != nil {
		return "", err
	}

	err = r.db.GetContext(ctx, &res, sqlStr, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, ErrNotFound
		}
		return res, err
	}

	return res, nil
}

func (r *Repository) UpsertUsers(ctx context.Context, tx *sqlx.Tx, users []models.User) ([]models.User, error) {
	if len(users) == 0 {
		return nil, nil
	}

	builder := st.
		Insert(r.tableName).
		Columns(r.columns.ForInsert()...)

	for _, user := range users {
		builder = builder.Values(GetValues(user)...)
	}

	sqlStr, params, err := builder.
		Suffix(r.columns.OnConflict() + " RETURNING *").
		ToSql()
	if err != nil {
		return nil, err
	}

	var result []models.User
	if err := tx.SelectContext(ctx, &result, sqlStr, params...); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) GetUpsertUserQuery(user models.User) (string, []interface{}, error) {
	queryBuilder := st.
		Insert(r.tableName).
		Columns(append([]string{r.columns.GetIDField()}, r.columns.ForInsert()...)...).
		Values(GetValues(user))

	suffix := r.columns.OnConflict()

	return queryBuilder.Suffix(suffix).ToSql()
}

func (r *Repository) Find(ctx context.Context, spec FindSpecification) ([]models.User, error) {
	queryBuilder := st.
		Select(spec.GetFields()...).From(r.tableName)

	sqlStr, params, err := spec.GetRule(queryBuilder).ToSql()
	if err != nil {
		return nil, err
	}

	var users []models.User
	err = r.db.SelectContext(ctx, &users, sqlStr, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return users, nil
}

func (r *Repository) UserUpdate(ctx context.Context, spec UpdateSpecification) (models.User, error) {
	var user models.User

	builder := st.Update(r.tableName)

	for col, val := range spec.GetSetValues() {
		builder = builder.Set(col, val)
	}

	builder = spec.GetRule(builder)

	if returning := spec.GetReturningFields(); len(returning) > 0 {
		builder = builder.Suffix(fmt.Sprintf("RETURNING %s", strings.Join(returning, ",")))
	}

	sqlStr, params, err := builder.ToSql()
	if err != nil {
		return user, err
	}

	err = r.db.GetContext(ctx, &user, sqlStr, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrNotFound
		}
		return user, err
	}

	return user, nil
}

func (r *Repository) BulkDeactivateUsers(ctx context.Context, tx *sqlx.Tx, teamName string, userIDs []string) ([]string, error) {
	spec := user_spec.NewBulkDeactivateTeamUsersSpecification(teamName, userIDs)

	builder := st.Update(r.tableName)

	for col, val := range spec.GetSetValues() {
		builder = builder.Set(col, val)
	}

	builder = spec.GetRule(builder)

	if returning := spec.GetReturningFields(); len(returning) > 0 {
		builder = builder.Suffix(fmt.Sprintf("RETURNING %s", strings.Join(returning, ",")))
	}

	sqlStr, params, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var deactivatedUserIDs []string
	err = tx.SelectContext(ctx, &deactivatedUserIDs, sqlStr, params...)
	if err != nil {
		return nil, err
	}

	return deactivatedUserIDs, nil
}
