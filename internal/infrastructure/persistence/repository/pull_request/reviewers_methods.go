package pull_request

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/specification/reviewer"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/specification/status"
)

type PRFullInfo struct {
	AuthorID             string
	AllReviewers         []string
	DeactivatedReviewers []string
}

func (r *Repository) GetAvailableReviewers(ctx context.Context, teamName string, excludeIDs []string, limit int) ([]string, error) {
	return r.FindReviewers(ctx, reviewer.NewGetAvailableReviewersSpecification(teamName, excludeIDs, limit, r.usersTableName))
}

func (r *Repository) GetPullRequestReviewers(ctx context.Context, prID string) ([]string, error) {
	return r.FindReviewers(ctx, reviewer.NewGetPRReviewersSpecification(prID, r.reviewersTableName))
}

func (r *Repository) InsertReviewers(ctx context.Context, tx *sqlx.Tx, prID string, reviewers []string) error {
	if len(reviewers) == 0 {
		return nil
	}

	builder := st.
		Insert(r.reviewersTableName).
		Columns(r.reviewerColumns.ForInsert()...)

	for _, r := range reviewers {
		builder = builder.Values(prID, r)
	}

	query, args, err := builder.
		Suffix(r.reviewerColumns.OnConflict()).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *Repository) FindReviewers(ctx context.Context, spec FindSpecification) ([]string, error) {
	queryBuilder := st.Select(spec.GetFields()...)

	sqlStr, params, err := spec.GetRule(queryBuilder).ToSql()
	if err != nil {
		return nil, err
	}

	var res []string
	err = r.db.SelectContext(ctx, &res, sqlStr, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return res, nil
}

func (r *Repository) ReassignReviewer(ctx context.Context, tx *sqlx.Tx, prID, oldReviewerID, newReviewerID string) error {
	removed, err := r.deleteReviewer(ctx, tx, prID, oldReviewerID)
	if err != nil {
		return err
	}

	if !removed {
		return ErrReviewerNotAssigned
	}

	if err := r.InsertReviewers(ctx, tx, prID, []string{newReviewerID}); err != nil {
		return err
	}

	return nil
}

func (r *Repository) deleteReviewer(ctx context.Context, tx *sqlx.Tx, prID, reviewerID string) (bool, error) {
	query, args, err := st.
		Delete(r.reviewersTableName).
		Where(sq.Eq{
			"pr_id":       prID,
			"reviewer_id": reviewerID,
		}).
		ToSql()
	if err != nil {
		return false, err
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

func (r *Repository) GetOpenPRsWithReviewers(ctx context.Context, reviewerIDs []string) (map[string][]string, error) {
	if len(reviewerIDs) == 0 {
		return make(map[string][]string), nil
	}

	openStatus, err := r.FindStatus(ctx, status.NewGetStatusByNameSpecification("OPEN"))
	if err != nil {
		return nil, err
	}

	queryBuilder := st.
		Select("r.pr_id", "r.reviewer_id").
		From(fmt.Sprintf("%s r", r.reviewersTableName)).
		Join(fmt.Sprintf("%s pr ON r.pr_id = pr.%s", r.tableName, r.pullRequestColumns.GetIDField())).
		Where(sq.Eq{"pr.status_id": openStatus.ID}).
		Where(sq.Eq{"r.reviewer_id": reviewerIDs})

	sqlStr, params, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	type prReviewerRow struct {
		PRID       string `db:"pr_id"`
		ReviewerID string `db:"reviewer_id"`
	}

	var rows []prReviewerRow
	err = r.db.SelectContext(ctx, &rows, sqlStr, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make(map[string][]string), nil
		}
		return nil, err
	}

	result := make(map[string][]string)
	for _, row := range rows {
		if reviewers, exists := result[row.PRID]; exists {
			result[row.PRID] = append(reviewers, row.ReviewerID)
		} else {
			result[row.PRID] = []string{row.ReviewerID}
		}
	}

	return result, nil
}

func (r *Repository) GetOpenPRsWithFullInfo(ctx context.Context, deactivatedReviewerIDs []string) (map[string]PRFullInfo, error) {
	if len(deactivatedReviewerIDs) == 0 {
		return make(map[string]PRFullInfo), nil
	}

	openStatus, err := r.FindStatus(ctx, status.NewGetStatusByNameSpecification("OPEN"))
	if err != nil {
		return nil, err
	}

	queryBuilder := st.
		Select([]string{"pr.pr_id", "pr.author_id", "r.reviewer_id"}...).
		From(fmt.Sprintf("%s pr", r.tableName)).
		Join(fmt.Sprintf("%s r ON r.pr_id = pr.%s", r.reviewersTableName, r.pullRequestColumns.GetIDField())).
		Where(sq.Eq{"pr.status_id": openStatus.ID}).
		Where(sq.Eq{"r.reviewer_id": deactivatedReviewerIDs})

	sqlStr, params, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	type prRow struct {
		PRID       string `db:"pr_id"`
		AuthorID   string `db:"author_id"`
		ReviewerID string `db:"reviewer_id"`
	}

	var rows []prRow
	err = r.db.SelectContext(ctx, &rows, sqlStr, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make(map[string]PRFullInfo), nil
		}
		return nil, err
	}

	prInfoMap := make(map[string]PRFullInfo)
	deactivatedSet := make(map[string]bool)
	for _, id := range deactivatedReviewerIDs {
		deactivatedSet[id] = true
	}

	for _, row := range rows {
		info, exists := prInfoMap[row.PRID]
		if !exists {
			info = PRFullInfo{
				AuthorID:             row.AuthorID,
				AllReviewers:         []string{},
				DeactivatedReviewers: []string{},
			}
		}

		if deactivatedSet[row.ReviewerID] {
			info.DeactivatedReviewers = append(info.DeactivatedReviewers, row.ReviewerID)
		}
		prInfoMap[row.PRID] = info
	}

	if len(prInfoMap) > 0 {
		prIDs := make([]string, 0, len(prInfoMap))
		for prID := range prInfoMap {
			prIDs = append(prIDs, prID)
		}

		queryBuilder = st.
			Select(r.reviewerColumns.ForSelect([]string{"*"})...).
			From(r.reviewersTableName).
			Where(sq.Eq{"pr_id": prIDs})

		sqlStr, params, err = queryBuilder.ToSql()
		if err != nil {
			return nil, err
		}

		var reviewerRows []models.Reviewer
		err = r.db.SelectContext(ctx, &reviewerRows, sqlStr, params...)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		reviewersByPR := make(map[string][]string)
		for _, row := range reviewerRows {
			reviewersByPR[row.PullRequestID] = append(reviewersByPR[row.PullRequestID], row.ReviewerID)
		}

		for prID, info := range prInfoMap {
			info.AllReviewers = reviewersByPR[prID]
			prInfoMap[prID] = info
		}
	}

	return prInfoMap, nil
}

func (r *Repository) BulkReassignReviewers(ctx context.Context, tx *sqlx.Tx, reassignments []PRReassignments) error {
	if len(reassignments) == 0 {
		return nil
	}

	for _, prReassignments := range reassignments {
		reviewerMap := prReassignments.Reassignments
		prID := prReassignments.PRID

		if len(reviewerMap) == 0 {
			continue
		}

		oldReviewerIDs := make([]string, 0, len(reviewerMap))
		newReviewerSet := make(map[string]bool)
		newReviewerIDs := make([]string, 0, len(reviewerMap))

		for oldID, newID := range reviewerMap {
			oldReviewerIDs = append(oldReviewerIDs, oldID)
			if !newReviewerSet[newID] {
				newReviewerSet[newID] = true
				newReviewerIDs = append(newReviewerIDs, newID)
			}
		}

		query, args, err := st.
			Delete(r.reviewersTableName).
			Where(sq.Eq{"pr_id": prID}).
			Where(sq.Eq{"reviewer_id": oldReviewerIDs}).
			ToSql()
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return err
		}

		if len(newReviewerIDs) > 0 {
			if err := r.InsertReviewers(ctx, tx, prID, newReviewerIDs); err != nil {
				return err
			}
		}
	}

	return nil
}
