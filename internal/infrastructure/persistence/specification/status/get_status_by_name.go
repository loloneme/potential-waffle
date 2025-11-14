package status

import sq "github.com/Masterminds/squirrel"

type GetStatusByNameSpecification struct {
	statusName string
}

func NewGetStatusByNameSpecification(statusName string) *GetStatusByNameSpecification {
	return &GetStatusByNameSpecification{statusName: statusName}
}

func (s *GetStatusByNameSpecification) GetFields() []string {
	return []string{"*"}
}

func (s *GetStatusByNameSpecification) GetRule(builder sq.SelectBuilder) sq.SelectBuilder {
	return builder.Where(sq.Eq{"status_name": s.statusName})
}
