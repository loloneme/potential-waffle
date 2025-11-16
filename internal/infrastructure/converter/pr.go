package converter

import (
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

func ToOpenAPIPullRequest(pr models.PullRequest) *generated.PullRequest {
	res := &generated.PullRequest{
		AssignedReviewers: pr.Reviewers,
		AuthorId:          pr.AuthorID,
		CreatedAt:         pr.CreatedAt,
		PullRequestId:     pr.ID,
		PullRequestName:   pr.Name,
	}

	if pr.Status != nil {
		res.Status = ToStatusEnum(pr.Status.Name)
	}

	if pr.MergedAt != nil {
		res.MergedAt = pr.MergedAt
	}
	return res
}

func ToOpenAPIPullRequestShort(pr models.PullRequest) *generated.PullRequestShort {
	return &generated.PullRequestShort{
		AuthorId:        pr.AuthorID,
		PullRequestId:   pr.ID,
		PullRequestName: pr.Name,
		Status:          ToShortStatusEnum(pr.Status.Name),
	}
}

func FromOpenAPIPullRequestCreate(pr *generated.PostPullRequestCreateJSONBody, status generated.PullRequestStatus) *models.PullRequest {
	return &models.PullRequest{
		ID:       pr.PullRequestId,
		Name:     pr.PullRequestName,
		AuthorID: pr.AuthorId,
		Status: &models.Status{
			Name: string(status),
		},
	}
}

func ToStatusEnum(name string) generated.PullRequestStatus {
	switch name {
	case "MERGED":
		return generated.PullRequestStatusMERGED
	case "OPEN":
		return generated.PullRequestStatusOPEN
	default:
		return ""
	}
}

func ToShortStatusEnum(name string) generated.PullRequestShortStatus {
	switch name {
	case "MERGED":
		return generated.PullRequestShortStatusMERGED
	case "OPEN":
		return generated.PullRequestShortStatusOPEN
	default:
		return generated.PullRequestShortStatusOPEN
	}
}
