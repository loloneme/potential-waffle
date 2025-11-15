package create_pr

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/specification/status"
)

const (
	numberOfReviewers = 2
)

type Service struct {
	userRepo userRepo
	prRepo   prRepo
}

func New(userRepo userRepo, prRepo prRepo) *Service {
	return &Service{
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

func (s *Service) CreatePR(ctx context.Context, pr *models.PullRequest) (models.PullRequest, error) {
	var createdPR models.PullRequest

	err := s.prRepo.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		foundStatus, err := s.prRepo.FindStatus(ctx, status.NewGetStatusByNameSpecification(pr.Status.Name))
		if err != nil {
			return fmt.Errorf("find status: %w", err)
		}
		pr.StatusID = foundStatus.ID

		created, err := s.prRepo.InsertPullRequest(ctx, tx, pr)
		if err != nil {
			return fmt.Errorf("insert pull request: %w", err)
		}
		createdPR = created
		createdPR.Status = foundStatus

		teamName, err := s.userRepo.GetUserTeamName(ctx, pr.AuthorID)
		if err != nil {
			return fmt.Errorf("get author team: %w", err)
		}

		reviewers, err := s.prRepo.GetAvailableReviewers(ctx, teamName, []string{pr.AuthorID}, numberOfReviewers)
		if err != nil {
			return fmt.Errorf("get available reviewers: %w", err)
		}

		if err := s.prRepo.InsertReviewers(ctx, tx, pr.ID, reviewers); err != nil {
			return fmt.Errorf("insert reviewers: %w", err)
		}

		createdPR.Reviewers = reviewers

		return nil
	})

	return createdPR, err
}
