package user

import sq "github.com/Masterminds/squirrel"

type GetUserTeamNameSpecification struct {
	userID string
}

func NewGetUserTeamNameSpecification(userID string) *GetUserTeamNameSpecification {
	return &GetUserTeamNameSpecification{userID: userID}
}

func (s *GetUserTeamNameSpecification) GetFields() []string {
	return []string{"team_name"}
}

func (s *GetUserTeamNameSpecification) GetRule(builder sq.SelectBuilder) sq.SelectBuilder {
	return builder.Where(sq.Eq{"user_id": s.userID})
}
