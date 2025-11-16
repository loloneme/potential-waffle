package pr_merge_post

import (
	"net/http"

	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/converter"
	rpc_errors "github.com/loloneme/potential-waffle/internal/rpc/errors"
)

type Handler struct {
	mergePRService MergePRService
}

func New(MergePRService MergePRService) *Handler {
	return &Handler{
		mergePRService: MergePRService,
	}
}

func (h *Handler) PRMergePost(ctx echo.Context) error {
	var input generated.PostPullRequestMergeJSONBody

	if err := ctx.Bind(&input); err != nil {
		return rpc_errors.RespondBadRequest(ctx, "")
	}

	pr, err := h.mergePRService.MergePullRequest(ctx.Request().Context(),
		input.PullRequestId,
		string(generated.PullRequestStatusMERGED))
	if err != nil {
		return rpc_errors.RespondFromError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"pr": converter.ToOpenAPIPullRequest(pr),
	})
}
