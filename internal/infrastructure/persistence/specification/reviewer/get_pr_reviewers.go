package reviewer

import sq "github.com/Masterminds/squirrel"

type GetPRReviewersSpecification struct {
	PullRequestID string
	FromTable     string
}

func NewGetPRReviewersSpecification(PullRequestID string, FromTable string) *GetPRReviewersSpecification {
	return &GetPRReviewersSpecification{
		PullRequestID: PullRequestID,
		FromTable:     FromTable,
	}
}

func (s *GetPRReviewersSpecification) GetRule(builder sq.SelectBuilder) sq.SelectBuilder {
	return builder.From(s.FromTable).Where(sq.Eq{"pr_id": s.PullRequestID}).OrderBy("reviewer_id ASC")
}

func (s *GetPRReviewersSpecification) GetFields() []string {
	return []string{"reviewer_id"}
}
