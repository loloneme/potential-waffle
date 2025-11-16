package pull_request

import (
	"context"
)

type UserAssignmentStats struct {
	UserID string `db:"user_id"`
	Count  int    `db:"count"`
}

type PRAssignmentStats struct {
	PullRequestID string `db:"pull_request_id"`
	Count         int    `db:"count"`
}

type Statistics struct {
	AssignmentsByUser []UserAssignmentStats
	AssignmentsByPR   []PRAssignmentStats
}

func (r *Repository) GetStatistics(ctx context.Context) (*Statistics, error) {
	assignmentsByUser, err := r.getAssignmentsByUser(ctx)
	if err != nil {
		return nil, err
	}

	assignmentsByPR, err := r.getAssignmentsByPR(ctx)
	if err != nil {
		return nil, err
	}

	return &Statistics{
		AssignmentsByUser: assignmentsByUser,
		AssignmentsByPR:   assignmentsByPR,
	}, nil
}

func (r *Repository) getAssignmentsByUser(ctx context.Context) ([]UserAssignmentStats, error) {
	query, args, err := st.
		Select("reviewer_id AS user_id", "COUNT(*) AS count").
		From(r.reviewersTableName).
		GroupBy("reviewer_id").
		OrderBy("count DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	var stats []UserAssignmentStats
	err = r.db.SelectContext(ctx, &stats, query, args...)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *Repository) getAssignmentsByPR(ctx context.Context) ([]PRAssignmentStats, error) {
	query, args, err := st.
		Select("pr_id AS pull_request_id", "COUNT(*) AS count").
		From(r.reviewersTableName).
		GroupBy("pr_id").
		OrderBy("count DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	var stats []PRAssignmentStats
	err = r.db.SelectContext(ctx, &stats, query, args...)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
