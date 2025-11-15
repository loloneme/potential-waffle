package reassign_pr

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/pull_request"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/user"
	rpc_errors "github.com/loloneme/potential-waffle/internal/rpc/errors"
)

const (
	numberOfReviewers = 1
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

func (s *Service) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (models.PullRequest, string, error) {
	pr, err := s.prRepo.GetPRByID(ctx, prID)
	if err != nil {
		if errors.Is(err, pull_request.ErrPRNotFound) {
			return models.PullRequest{}, "", rpc_errors.NewNotFound("PR not found")
		}
		return models.PullRequest{}, "", fmt.Errorf("get PR: %w", err)
	}

	if pr.Status.Name == "MERGED" {
		return models.PullRequest{}, "", rpc_errors.NewPRMerged("cannot reassign on merged PR")
	}

	teamName, err := s.userRepo.GetUserTeamName(ctx, oldReviewerID)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return models.PullRequest{}, "", rpc_errors.NewNotFound("reviewer not found")
		}
		return models.PullRequest{}, "", fmt.Errorf("get reviewer team: %w", err)
	}

	currentReviewers, err := s.prRepo.GetPullRequestReviewers(ctx, prID)
	if err != nil {
		return models.PullRequest{}, "", fmt.Errorf("get current reviewers: %w", err)
	}

	// Проверяем, что oldReviewerID действительно назначен на PR
	isAssigned := false
	for _, reviewerID := range currentReviewers {
		if reviewerID == oldReviewerID {
			isAssigned = true
			break
		}
	}
	if !isAssigned {
		return models.PullRequest{}, "", rpc_errors.NewNotAssigned("reviewer is not assigned to this PR")
	}

	excludeIDs := []string{oldReviewerID, pr.AuthorID}
	excludeIDs = append(excludeIDs, currentReviewers...)

	availableReviewers, err := s.prRepo.GetAvailableReviewers(ctx, teamName, excludeIDs, numberOfReviewers)
	if err != nil {
		if errors.Is(err, pull_request.ErrReviewersNotFound) {
			return models.PullRequest{}, "", rpc_errors.NewNoCandidate("no available reviewers in team")
		}
		return models.PullRequest{}, "", fmt.Errorf("get available reviewers: %w", err)
	}

	if len(availableReviewers) == 0 {
		return models.PullRequest{}, "", rpc_errors.NewNoCandidate("no available reviewers in team")
	}

	newReviewerID := availableReviewers[0]

	err = s.prRepo.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		err := s.prRepo.ReassignReviewer(ctx, tx, prID, oldReviewerID, newReviewerID)
		if err != nil {
			if errors.Is(err, pull_request.ErrReviewerNotAssigned) {
				return rpc_errors.NewNotAssigned("reviewer is not assigned to this PR")
			}
			return fmt.Errorf("reassign reviewer: %w", err)
		}
		return nil
	})

	if err != nil {
		return models.PullRequest{}, "", err
	}

	var reassignedPR models.PullRequest
	reassignedPR, err = s.prRepo.GetPRByID(ctx, prID)
	if err != nil {
		if errors.Is(err, pull_request.ErrPRNotFound) {
			return models.PullRequest{}, "", rpc_errors.NewNotFound("PR not found")
		}
		return models.PullRequest{}, "", fmt.Errorf("get updated PR: %w", err)
	}

	return reassignedPR, newReviewerID, nil
}
