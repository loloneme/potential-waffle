package team

import (
	"context"

	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

type teamRepository interface {
	CreateTeam(ctx context.Context, team models.Team) (models.Team, error)
	FindTeamByID(ctx context.Context, teamName string) (models.Team, error)
}
