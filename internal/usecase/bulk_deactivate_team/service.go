package bulk_deactivate_team

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/pull_request"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/team"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/user"
	user_spec "github.com/loloneme/potential-waffle/internal/infrastructure/persistence/specification/user"
	rpc_errors "github.com/loloneme/potential-waffle/internal/rpc/errors"
)

const (
	numberOfReviewers = 1
)

type ReassignmentResult struct {
	PRID          string
	OldReviewerID string
	NewReviewerID string
}

type BulkDeactivateResult struct {
	DeactivatedUserIDs []string
	Reassignments      []ReassignmentResult
}

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

func (s *Service) BulkDeactivateTeamUsers(ctx context.Context, teamName string, userIDs []string) (BulkDeactivateResult, error) {
	var res BulkDeactivateResult
	if len(userIDs) == 0 {
		return res, nil
	}

	users, err := s.userRepo.Find(ctx, user_spec.NewGetUsersByTeamNameSpec(teamName))
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return res, rpc_errors.NewNotFound("team not found or has no users")
		}
		return res, fmt.Errorf("find team users: %w", err)
	}

	activeTeamUserMap := make(map[string]bool)
	activeTeamUserIDs := make([]string, 0, len(users))
	for _, u := range users {
		if u.IsActive {
			activeTeamUserMap[u.ID] = true
			activeTeamUserIDs = append(activeTeamUserIDs, u.ID)
		}
	}

	var validUserIDs []string
	for _, userID := range userIDs {
		if activeTeamUserMap[userID] {
			validUserIDs = append(validUserIDs, userID)
		}
	}

	if len(validUserIDs) == 0 {
		return res, nil
	}

	prsInfo, err := s.prRepo.GetOpenPRsWithFullInfo(ctx, validUserIDs)
	if err != nil {
		return res, fmt.Errorf("get open PRs with full info: %w", err)
	}

	var reassignments []ReassignmentResult
	err = s.prRepo.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		deactivatedUserIDs, err := s.userRepo.BulkDeactivateUsers(ctx, tx, teamName, validUserIDs)
		if err != nil {
			return fmt.Errorf("bulk deactivate users: %w", err)
		}

		bulkReassignments := make([]pull_request.PRReassignments, 0)

		for prID, prInfo := range prsInfo {
			if len(prInfo.DeactivatedReviewers) == 0 {
				continue
			}

			excludeSet := make(map[string]bool)
			excludeSet[prInfo.AuthorID] = true
			for _, reviewerID := range prInfo.AllReviewers {
				excludeSet[reviewerID] = true
			}

			neededReviewers := len(prInfo.DeactivatedReviewers)
			availableReviewers := make([]string, 0, neededReviewers)
			for _, userID := range activeTeamUserIDs {
				if !excludeSet[userID] {
					availableReviewers = append(availableReviewers, userID)
					if len(availableReviewers) >= neededReviewers {
						break
					}
				}
			}

			reviewerMap := make(map[string]string)
			availableIdx := 0

			for _, oldReviewerID := range prInfo.DeactivatedReviewers {
				if availableIdx >= len(availableReviewers) {
					break
				}

				newReviewerID := availableReviewers[availableIdx]
				availableIdx++

				reviewerMap[oldReviewerID] = newReviewerID
				reassignments = append(reassignments, ReassignmentResult{
					PRID:          prID,
					OldReviewerID: oldReviewerID,
					NewReviewerID: newReviewerID,
				})

				excludeSet[newReviewerID] = true
			}

			if len(reviewerMap) > 0 {
				bulkReassignments = append(bulkReassignments, pull_request.PRReassignments{
					PRID:          prID,
					Reassignments: reviewerMap,
				})
			}
		}

		if len(bulkReassignments) > 0 {
			if err := s.prRepo.BulkReassignReviewers(ctx, tx, bulkReassignments); err != nil {
				return fmt.Errorf("bulk reassign reviewers: %w", err)
			}
		}

		res.DeactivatedUserIDs = deactivatedUserIDs
		res.Reassignments = reassignments
		return nil
	})
	if err != nil {
		if errors.Is(err, team.ErrNotFound) {
			return BulkDeactivateResult{}, rpc_errors.NewNotFound("team not found")
		}
		return BulkDeactivateResult{}, err
	}

	return res, nil
}
