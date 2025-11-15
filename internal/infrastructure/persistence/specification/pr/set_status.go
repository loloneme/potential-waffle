package pr_spec

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

type SetStatusSpecification struct {
	Status        *models.Status
	pullRequestID string
}

func NewSetStatusSpecification(status *models.Status, pullRequestID string) *SetStatusSpecification {
	return &SetStatusSpecification{
		Status:        status,
		pullRequestID: pullRequestID,
	}
}

func (s *SetStatusSpecification) GetSetValues() map[string]interface{} {
	result := map[string]interface{}{
		"status_id": s.Status.ID,
	}
	if s.Status.Name == "MERGED" {
		result["merged_at"] = sq.Expr("COALESCE(merged_at, NOW())")
	}
	return result
}

func (s *SetStatusSpecification) GetRule(builder sq.UpdateBuilder) sq.UpdateBuilder {
	return builder.Where(sq.Eq{"pr_id": s.pullRequestID})
}

func (s *SetStatusSpecification) GetReturningFields() []string {
	return []string{"*"}
}
