package pull_request

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func (r *Repository) InsertReviewers(ctx context.Context, tx *sqlx.Tx, prID string, reviewers []string) error {
	if len(reviewers) == 0 {
		return nil
	}

	builder := st.
		Insert(reviewersTable).
		Columns(r.reviewerColumns.ForInsert()...)

	for _, reviewer := range reviewers {
		builder = builder.Values(prID, reviewer)
	}

	query, args, err := builder.
		Suffix(r.reviewerColumns.OnConflict()).
		ToSql()

	if err != nil {
		return fmt.Errorf("build insert reviewers query: %w", err)
	}

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("exec insert reviewers query: %w", err)
	}

	return nil
}

func (r *Repository) FindReviewers(ctx context.Context, spec FindSpecification) ([]string, error) {
	queryBuilder := st.Select(spec.GetFields()...)

	sqlStr, params, err := spec.GetRule(queryBuilder).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build find reviewers query: %w", err)
	}

	var res []string
	err = r.db.SelectContext(ctx, &res, sqlStr, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("exec find users: %w", err)
	}
	return res, nil
}

func (r *Repository) ReassignReviewer(ctx context.Context, tx *sqlx.Tx, prID, oldReviewerID, newReviewerID string) error {
	removed, err := r.deleteReviewer(ctx, tx, prID, oldReviewerID)

	if err != nil {
		_ = tx.Rollback()

		return err
	}

	if !removed {
		_ = tx.Rollback()

		return ErrReviewerNotAssigned
	}

	if err := r.InsertReviewers(ctx, tx, prID, []string{newReviewerID}); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) deleteReviewer(ctx context.Context, tx *sqlx.Tx, prID, reviewerID string) (bool, error) {
	query, args, err := st.
		Delete(reviewersTable).
		Where(sq.Eq{
			"pr_id":       prID,
			"reviewer_id": reviewerID,
		}).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("build delete reviewer query: %w", err)
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("exec delete reviewer query: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("delete reviewer rows affected: %w", err)
	}

	return affected > 0, nil
}
