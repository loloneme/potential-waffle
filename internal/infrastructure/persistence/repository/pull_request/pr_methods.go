package pull_request

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	pr_spec "github.com/loloneme/potential-waffle/internal/infrastructure/persistence/specification/pr"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/specification/reviewer"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/specification/status"
)

const (
	reviewersTable = "pull_request_reviewers"
	reviewersAlias = "prr"
)

var (
	st = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	ErrPRNotFound        = errors.New("pull request not found")
	ErrStatusNotFound    = errors.New("pull request status not found")
	ErrReviewersNotFound = errors.New("there are no available reviewers")

	ErrReviewerNotAssigned = errors.New("reviewer not assigned to pull request")
)

type FindSpecification interface {
	GetFields() []string
	GetRule(s sq.SelectBuilder) sq.SelectBuilder
}

type UpdateSpecification interface {
	GetSetValues() map[string]interface{}
	GetRule(builder sq.UpdateBuilder) sq.UpdateBuilder
	GetReturningFields() []string
}

type contextExecutor interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...any) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func getPRValues(pr *models.PullRequest) []interface{} {
	return []interface{}{
		pr.ID,
		pr.Name,
		pr.AuthorID,
		pr.Status,
	}
}

func (r *Repository) InsertPullRequest(ctx context.Context, tx *sqlx.Tx, pr *models.PullRequest) (models.PullRequest, error) {
	var created models.PullRequest

	query, args, err := st.
		Insert(r.tableName).
		Columns(r.pullRequestColumns.ForInsert()...).
		Values(getPRValues(pr)...).
		Suffix(fmt.Sprintf("RETURNING *")).
		ToSql()

	if err != nil {
		return created, fmt.Errorf("build insert pull request query: %w", err)
	}

	if err := tx.GetContext(ctx, &created, query, args...); err != nil {
		return created, fmt.Errorf("exec insert pull request query: %w", err)
	}

	return created, nil
}

func (r *Repository) MergePullRequest(ctx context.Context, prID string, statusID int64) (models.PullRequest, error) {
	spec := pr_spec.NewSetStatusSpecification(statusID, prID)
	return r.UpdatePullRequest(ctx, spec)
}

func (r *Repository) UpdatePullRequest(ctx context.Context, spec UpdateSpecification) (models.PullRequest, error) {
	var pr models.PullRequest

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
		return pr, fmt.Errorf("build update pr: %w", err)
	}

	err = r.db.GetContext(ctx, &pr, sqlStr, params...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pr, ErrPRNotFound
		}
		return pr, fmt.Errorf("exec update pr: %w", err)
	}

	return pr, nil
}

func (r *Repository) GetPRByID(ctx context.Context, prID string) (models.PullRequest, error) {
	var pr models.PullRequest

	query, args, err := st.
		Select(r.pullRequestColumns.ForSelect()...).
		From(r.tableName).
		Where(sq.Eq{r.pullRequestColumns.GetIDField(): prID}).
		ToSql()

	if err != nil {
		return pr, fmt.Errorf("build get pull request query: %w", err)
	}

	if err := r.db.GetContext(ctx, &pr, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pr, ErrPRNotFound
		}
		return pr, fmt.Errorf("exec get pull request query: %w", err)
	}

	pr.Status, err = r.FindStatus(ctx, status.NewGetStatusByIDSpecification(pr.StatusID))
	if err != nil {
		return models.PullRequest{}, fmt.Errorf("get pull request status: %w", err)
	}

	pr.Reviewers, err = r.getReviewers(ctx, pr.ID)
	if err != nil {
		return models.PullRequest{}, fmt.Errorf("get pull request reviewers: %w", err)
	}

	return pr, nil
}

func (r *Repository) getReviewers(ctx context.Context, prID string) ([]string, error) {
	return r.FindReviewers(ctx, reviewer.NewGetPRReviewersSpecification(prID, r.reviewersTableName))
}
