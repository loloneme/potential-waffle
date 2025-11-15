package create_team

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

type Service struct {
	teamRepo teamRepo
	userRepo userRepo
}

func New(teamRepo teamRepo, userRepo userRepo) *Service {
	return &Service{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (s *Service) CreateTeam(ctx context.Context, team *models.Team) (models.Team, error) {
	var createdTeam models.Team

	err := s.teamRepo.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		created, err := s.teamRepo.CreateTeam(ctx, tx, *team)
		if err != nil {
			return err
		}

		createdTeam = created

		members, err := s.userRepo.UpsertUsers(ctx, tx, team.Members)
		if err != nil {
			return err
		}

		createdTeam.Members = members

		return nil
	})

	return createdTeam, err
}
