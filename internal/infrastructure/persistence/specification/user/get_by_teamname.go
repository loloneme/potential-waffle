package user

import sq "github.com/Masterminds/squirrel"

type GetUsersByTeamNameSpec struct {
	teamName string
}

func NewGetUsersByTeamNameSpec(teamName string) *GetUsersByTeamNameSpec {
	return &GetUsersByTeamNameSpec{teamName: teamName}
}

func (s *GetUsersByTeamNameSpec) GetFields() []string {
	return []string{"user_id", "username", "is_active"}
}

func (s *GetUsersByTeamNameSpec) GetRule(builder sq.SelectBuilder) sq.SelectBuilder {
	return builder.Where(sq.Eq{"team_name": s.teamName})
}
