package pull_request

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	pr_spec "github.com/loloneme/potential-waffle/internal/infrastructure/persistence/specification/pr"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/specification/reviewer"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/specification/status"
)

var (
	st = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	ErrPRNotFound        = errors.New("pull request not found")
	ErrStatusNotFound    = errors.New("pull request status not found")
	ErrReviewersNotFound = errors.New("there are no available reviewers")

	ErrPRAlreadyExists = errors.New("pull request already exists")

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
		pr.StatusID,
	}
}

func (r *Repository) InsertPullRequest(ctx context.Context, tx *sqlx.Tx, pr *models.PullRequest) (models.PullRequest, error) {
	var created models.PullRequest

	query, args, err := st.
		Insert(r.tableName).
		Columns(r.pullRequestColumns.ForInsert()...).
		Values(getPRValues(pr)...).
		Suffix(fmt.Sprintf("ON CONFLICT (%s) DO NOTHING RETURNING *", r.pullRequestColumns.GetIDField())).
		ToSql()

	if err != nil {
		return created, fmt.Errorf("build insert pull request query: %w", err)
	}

	if err := tx.GetContext(ctx, &created, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return created, ErrPRAlreadyExists
		}
		return created, fmt.Errorf("insert pull request: %w", err)
	}

	return created, nil
}

func (r *Repository) PullRequestExists(ctx context.Context, prID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pr_id = $1)`

	if err := r.db.GetContext(ctx, &exists, query, prID); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repository) SetPullRequestStatus(ctx context.Context, prID string, statusName string) error {
	foundStatus, err := r.FindStatus(ctx, status.NewGetStatusByNameSpecification(statusName))
	if err != nil {
		return fmt.Errorf("find pull request status: %w", err)
	}

	spec := pr_spec.NewSetStatusSpecification(foundStatus, prID)
	return r.UpdatePullRequest(ctx, spec)
}

func (r *Repository) UpdatePullRequest(ctx context.Context, spec UpdateSpecification) error {
	builder := st.Update(r.tableName)

	for col, val := range spec.GetSetValues() {
		builder = builder.Set(col, val)
	}

	builder = spec.GetRule(builder)

	sqlStr, params, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("build update pr: %w", err)
	}

	_, err = r.db.ExecContext(ctx, sqlStr, params...)

	if err != nil {
		return fmt.Errorf("exec update pr: %w", err)
	}

	return nil
}

func (r *Repository) GetPRByID(ctx context.Context, prID string) (models.PullRequest, error) {
	var pr models.PullRequest

	query, args, err := st.
		Select(r.pullRequestColumns.ForSelect(nil)...).
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

func (r *Repository) GetPRByIDShort(ctx context.Context, prID string) (models.PullRequest, error) {
	var pr models.PullRequest

	query, args, err := st.
		Select(r.pullRequestColumns.ForSelect([]string{"pr_id", "pr_name", "author_id", "status_id"})...).
		From(fmt.Sprintf("%s %s", r.tableName, r.pullRequestColumns.GetAlias())).
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

	return pr, nil
}

func (r *Repository) FindPullRequests(ctx context.Context, spec FindSpecification) ([]models.PullRequest, error) {
	queryBuilder := st.
		Select(spec.GetFields()...).From(r.reviewersTableName)

	sqlStr, params, err := spec.GetRule(queryBuilder).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build find pull requests query: %w", err)
	}

	var pullRequests []models.PullRequest
	err = r.db.SelectContext(ctx, &pullRequests, sqlStr, params...)
	if err != nil {
		return nil, fmt.Errorf("exec find pull requests: %w", err)
	}

	for i, pr := range pullRequests {
		pr, err = r.GetPRByIDShort(ctx, pr.ID)
		if err != nil {
			return nil, fmt.Errorf("get pull request: %w", err)
		}
		pullRequests[i] = pr
	}

	return pullRequests, nil
}

func (r *Repository) getReviewers(ctx context.Context, prID string) ([]string, error) {
	return r.FindReviewers(ctx, reviewer.NewGetPRReviewersSpecification(prID, r.reviewersTableName))
}
