package reviewer

import sq "github.com/Masterminds/squirrel"

type GetAvailableReviewersSpecification struct {
	ExcludeIDs []string
	TeamName   string
	Limit      int
	FromTable  string
}

func NewGetAvailableReviewersSpecification(teamName string, excludeIDs []string, limit int, fromTable string) *GetAvailableReviewersSpecification {
	return &GetAvailableReviewersSpecification{
		ExcludeIDs: excludeIDs,
		TeamName:   teamName,
		Limit:      limit,

		FromTable: fromTable,
	}
}

func (s *GetAvailableReviewersSpecification) GetRule(builder sq.SelectBuilder) sq.SelectBuilder {
	builder = builder.From(s.FromTable).
		Where(sq.Eq{"team_name": s.TeamName}).
		Where(sq.Eq{"is_active": true})

	if len(s.ExcludeIDs) > 0 {
		builder = builder.Where(sq.NotEq{"user_id": s.ExcludeIDs})
	}

	return builder.Limit(uint64(s.Limit))
}

func (s *GetAvailableReviewersSpecification) GetFields() []string {
	return []string{"user_id"}
}
