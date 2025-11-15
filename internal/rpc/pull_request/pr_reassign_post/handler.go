package pr_reassign_post

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/converter"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/pull_request"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/user"
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
		return ctx.JSON(http.StatusBadRequest, "bad request")
	}

	pr, newReviewerID, err := h.reassignPRService.ReassignReviewer(ctx.Request().Context(), input.PullRequestId, input.OldUserId)
	if err != nil {
		resp := generated.ErrorResponse{}

		switch err := err.(type) {
		case *rpc_errors.NotFoundError:
			resp.Error.Code = generated.NOTFOUND
			resp.Error.Message = err.Message

			return ctx.JSON(http.StatusNotFound, map[string]interface{}{
				"error": map[string]string{
					"code":    "NOT_FOUND",
					"message": err.Message,
				},
			})

		case *rpc_errors.PRMergedError:
			resp.Error.Code = generated.PRMERGED
			resp.Error.Message = err.Message

			return ctx.JSON(http.StatusConflict, resp)

		case *rpc_errors.NotAssignedError:
			resp.Error.Code = generated.NOTASSIGNED
			resp.Error.Message = err.Message

			return ctx.JSON(http.StatusConflict, resp)

		case *rpc_errors.NoCandidateError:
			resp.Error.Code = generated.NOCANDIDATE
			resp.Error.Message = err.Message

			return ctx.JSON(http.StatusConflict, resp)

		default:
			if errors.Is(err, pull_request.ErrPRNotFound) || errors.Is(err, user.ErrNotFound) {
				resp.Error.Code = generated.NOTFOUND
				resp.Error.Message = "resource not found"

				return ctx.JSON(http.StatusNotFound, map[string]interface{}{
					"error": map[string]string{
						"code":    "NOT_FOUND",
						"message": "resource not found",
					},
				})
			}

			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": map[string]string{
					"code":    "INTERNAL",
					"message": err.Error(),
				},
			})
		}
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"pr":          converter.ToOpenAPIPullRequest(pr),
		"replaced_by": newReviewerID,
	})
}
