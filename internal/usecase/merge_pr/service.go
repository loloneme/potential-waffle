package merge_pr

import (
	"context"
	"errors"

	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/pull_request"
	rpc_errors "github.com/loloneme/potential-waffle/internal/rpc/errors"
)

type Service struct {
	prRepo prRepo
}

func New(prRepo prRepo) *Service {
	return &Service{
		prRepo: prRepo,
	}
}

func (s *Service) MergePullRequest(ctx context.Context, prID string, mergeStatus string) (models.PullRequest, error) {
	if exists, err := s.prRepo.PullRequestExists(ctx, prID); err != nil {
		return models.PullRequest{}, err
	} else if !exists {
		return models.PullRequest{}, rpc_errors.NewNotFound("PR not found")
	}

	err := s.prRepo.SetPullRequestStatus(ctx, prID, mergeStatus)
	if err != nil {
		return models.PullRequest{}, err
	}

	pr, err := s.prRepo.GetPRByID(ctx, prID)
	if err != nil {
		if errors.Is(err, pull_request.ErrPRNotFound) {
			return models.PullRequest{}, rpc_errors.NewNotFound("PR not found")
		}
		return models.PullRequest{}, err
	}

	return pr, nil
}
