package pr_merge_post

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/converter"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/pull_request"
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
		return ctx.JSON(http.StatusBadRequest, "bad request")
	}

	pr, err := h.mergePRService.MergePullRequest(ctx.Request().Context(),
		input.PullRequestId,
		string(generated.PullRequestStatusMERGED))

	if err != nil {
		resp := generated.ErrorResponse{}

		switch {
		case errors.Is(err, pull_request.ErrPRNotFound):
			resp.Error.Code = generated.NOTFOUND
			resp.Error.Message = "resource not found"

			return ctx.JSON(http.StatusNotFound, map[string]interface{}{
				"error": map[string]string{
					"code":    "NOT_FOUND",
					"message": "resource not found",
				},
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": map[string]string{
					"code":    "INTERNAL",
					"message": err.Error(),
				},
			})
		}
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"pr": converter.ToOpenAPIPullRequest(pr),
	})
}
