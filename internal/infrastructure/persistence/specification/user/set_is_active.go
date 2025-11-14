package user

import sq "github.com/Masterminds/squirrel"

type SetIsActiveSpecification struct {
	userID   string
	isActive bool
}

func newSetIsActiveSpecification(userID string, isActive bool) *SetIsActiveSpecification {
	return &SetIsActiveSpecification{
		userID:   userID,
		isActive: isActive,
	}
}

func (s *SetIsActiveSpecification) GetSetValues() map[string]interface{} {
	return map[string]interface{}{
		"is_active": s.isActive,
	}
}

func (s *SetIsActiveSpecification) GetRule(builder sq.UpdateBuilder) sq.UpdateBuilder {
	return builder.Where(sq.Eq{"user_id": s.userID})
}

func (s *SetIsActiveSpecification) GetReturningFields() []string {
	return []string{"*"}
}
