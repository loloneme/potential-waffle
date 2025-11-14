package status

import sq "github.com/Masterminds/squirrel"

type GetStatusByIDSpecification struct {
	statusID int64
}

func NewGetStatusByIDSpecification(statusID int64) *GetStatusByIDSpecification {
	return &GetStatusByIDSpecification{statusID: statusID}
}

func (s *GetStatusByIDSpecification) GetFields() []string {
	return []string{"*"}
}

func (s *GetStatusByIDSpecification) GetRule(builder sq.SelectBuilder) sq.SelectBuilder {
	return builder.Where(sq.Eq{"status_id": s.statusID})
}
