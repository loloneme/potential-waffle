package pr_reassign_post

import (
	"net/http"

	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/converter"
	rpc_errors "github.com/loloneme/potential-waffle/internal/rpc/errors"
)

type Handler struct {
	reassignPRService reassignPRService
}

func New(reassignPRService reassignPRService) *Handler {
	return &Handler{
		reassignPRService: reassignPRService,
	}
}

func (h *Handler) PRReassignPost(ctx echo.Context) error {
	var input generated.PostPullRequestReassignJSONBody

	if err := ctx.Bind(&input); err != nil {
		return rpc_errors.RespondBadRequest(ctx, "")
	}

	pr, newReviewerID, err := h.reassignPRService.ReassignReviewer(ctx.Request().Context(), input.PullRequestId, input.OldUserId)
	if err != nil {
		return rpc_errors.RespondFromError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"pr":          converter.ToOpenAPIPullRequest(pr),
		"replaced_by": newReviewerID,
	})
}
