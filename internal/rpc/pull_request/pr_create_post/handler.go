package pr_create_post

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/converter"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/pull_request"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/user"
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
		return ctx.JSON(http.StatusBadRequest, "bad request")
	}

	prModel := &models.PullRequest{
		ID:       input.PullRequestId,
		Name:     input.PullRequestName,
		AuthorID: input.AuthorId,
		Status: &models.Status{
			Name: string(generated.PullRequestStatusOPEN),
		},
	}

	pullRequest, err := h.createPRService.CreatePR(ctx.Request().Context(), prModel)
	if err != nil {
		resp := generated.ErrorResponse{}

		switch {
		case errors.Is(err, pull_request.ErrPRAlreadyExists):
			resp.Error.Code = generated.PREXISTS
			resp.Error.Message = "PR already exists"

			return ctx.JSON(409, resp)

		case errors.Is(err, pull_request.ErrReviewersNotFound) || errors.Is(err, user.ErrNotFound):
			resp.Error.Code = generated.PREXISTS
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

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"pr": converter.ToOpenAPIPullRequest(pullRequest),
	})
}
