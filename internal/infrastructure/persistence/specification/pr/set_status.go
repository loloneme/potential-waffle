package pr_spec

import sq "github.com/Masterminds/squirrel"

type SetStatusSpecification struct {
	statusID      int64
	pullRequestID string
}

func NewSetStatusSpecification(statusID int64, pullRequestID string) *SetStatusSpecification {
	return &SetStatusSpecification{
		statusID:      statusID,
		pullRequestID: pullRequestID,
	}
}

func (s *SetStatusSpecification) GetSetValues() map[string]interface{} {
	return map[string]interface{}{
		"status_id": s.statusID,
		"merged_at": sq.Expr("COALESCE(merged_at, NOW())"),
	}
}

func (s *SetStatusSpecification) GetRule(builder sq.UpdateBuilder) sq.UpdateBuilder {
	return builder.Where(sq.Eq{"pr_id": s.pullRequestID})
}

func (s *SetStatusSpecification) GetReturningFields() []string {
	return []string{"*"}
}
