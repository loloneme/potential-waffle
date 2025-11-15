package pr_spec

import sq "github.com/Masterminds/squirrel"

type GetPRByReviewerSpecification struct {
	ReviewerID string
}

func NewGetPRByReviewerSpecification(reviewerID string) *GetPRByReviewerSpecification {
	return &GetPRByReviewerSpecification{
		ReviewerID: reviewerID,
	}
}

func (s *GetPRByReviewerSpecification) GetRule(builder sq.SelectBuilder) sq.SelectBuilder {
	return builder.Where(sq.Eq{"reviewer_id": s.ReviewerID})
}

func (s *GetPRByReviewerSpecification) GetFields() []string {
	return []string{"pr_id"}
}
