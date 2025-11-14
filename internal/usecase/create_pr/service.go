package create_pr

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
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

	err := s.repo.WithTx(ctx, func(tx *sqlx.Tx) error {
		statusID, _ := s.repo.GetStatusID(ctx, tx, "OPEN")
		pr.StatusID = statusID

		created, _ := s.repo.insertPullRequest(ctx, tx, pr)
		createdPR = created

		teamName, _ := s.repo.GetTeamName(ctx, tx, pr.AuthorID)

		spec := NewGetAvailableReviewersSpecification([]string{pr.AuthorID}, teamName, 2)
		reviewers, _ := s.repo.FindReviewers(ctx, tx, spec)

		return s.repo.InsertReviewers(ctx, tx, pr.ID, reviewers)
	})

	return createdPR, err
}
