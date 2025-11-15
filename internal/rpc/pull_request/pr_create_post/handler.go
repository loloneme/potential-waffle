package pr_create_post

import (
	"net/http"

	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/converter"
	rpc_errors "github.com/loloneme/potential-waffle/internal/rpc/errors"
)

type Handler struct {
	createPRService createPRService
}

func New(createPRService createPRService) *Handler {
	return &Handler{
		createPRService: createPRService,
	}
}

func (h *Handler) PRCreatePost(ctx echo.Context) error {
	var input generated.PostPullRequestCreateJSONBody

	if err := ctx.Bind(&input); err != nil {
		return rpc_errors.RespondBadRequest(ctx, "")
	}

	prModel := converter.FromOpenAPIPullRequestCreate(&input, generated.PullRequestStatusOPEN)

	pullRequest, err := h.createPRService.CreatePR(ctx.Request().Context(), prModel)
	if err != nil {
		return rpc_errors.RespondFromError(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"pr": converter.ToOpenAPIPullRequest(pullRequest),
	})
}
