package converter

import (
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

func ToOpenAPITeam(team models.Team) generated.Team {
	members := make([]generated.TeamMember, len(team.Members))
	for i, m := range team.Members {
		members[i] = generated.TeamMember{
			IsActive: m.IsActive,
			UserId:   m.ID,
			Username: m.Username,
		}
	}
	return generated.Team{
		Members:  members,
		TeamName: team.TeamName,
	}
}

func ToModelTeam(team generated.Team) *models.Team {
	users := make([]models.User, len(team.Members))
	for i, m := range team.Members {
		users[i] = models.User{
			ID:       m.UserId,
			Username: m.Username,
			IsActive: m.IsActive,
			TeamName: team.TeamName,
		}
	}

	return &models.Team{
		TeamName: team.TeamName,
		Members:  users,
	}
}
