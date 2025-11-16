package user

import sq "github.com/Masterminds/squirrel"

type BulkDeactivateTeamUsersSpecification struct {
	teamName string
	userIDs  []string
}

func NewBulkDeactivateTeamUsersSpecification(teamName string, userIDs []string) *BulkDeactivateTeamUsersSpecification {
	return &BulkDeactivateTeamUsersSpecification{
		teamName: teamName,
		userIDs:  userIDs,
	}
}

func (s *BulkDeactivateTeamUsersSpecification) GetSetValues() map[string]interface{} {
	return map[string]interface{}{
		"is_active": false,
	}
}

func (s *BulkDeactivateTeamUsersSpecification) GetRule(builder sq.UpdateBuilder) sq.UpdateBuilder {
	builder = builder.
		Where(sq.Eq{"team_name": s.teamName}).
		Where(sq.Eq{"is_active": true})

	if len(s.userIDs) > 0 {
		builder = builder.Where(sq.Eq{"user_id": s.userIDs})
	}

	return builder
}

func (s *BulkDeactivateTeamUsersSpecification) GetReturningFields() []string {
	return []string{"user_id"}
}
