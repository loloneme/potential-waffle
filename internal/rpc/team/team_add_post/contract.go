package team_add_post

import (
	"context"

	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

type createTeamService interface {
	CreateTeam(ctx context.Context, team *models.Team) (models.Team, error)
}
