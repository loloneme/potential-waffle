package users_get_review_get

import (
	"net/http"

	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/converter"
	pr_spec "github.com/loloneme/potential-waffle/internal/infrastructure/persistence/specification/pr"
)

type Handler struct {
	prRepo prRepo
}

func New(prRepo prRepo) *Handler {
	return &Handler{
		prRepo: prRepo,
	}
}

func (h *Handler) UsersGetReviewGet(ctx echo.Context, params generated.GetUsersGetReviewParams) error {
	pullRequests, err := h.prRepo.FindPullRequests(ctx.Request().Context(), pr_spec.NewGetPRByReviewerSpecification(params.UserId))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "internal server error")
	}

	shortPRs := make([]generated.PullRequestShort, len(pullRequests))
	for i, pr := range pullRequests {
		shortPR := converter.ToOpenAPIPullRequestShort(pr)
		if shortPR != nil {
			shortPRs[i] = *shortPR
		}
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"user_id":       params.UserId,
		"pull_requests": shortPRs,
	})
}
