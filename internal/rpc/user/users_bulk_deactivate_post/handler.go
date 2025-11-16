package users_bulk_deactivate_post

import (
	"net/http"

	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	rpc_errors "github.com/loloneme/potential-waffle/internal/rpc/errors"
)

type Handler struct {
	bulkDeactivateService bulkDeactivateService
}

func New(bulkDeactivateService bulkDeactivateService) *Handler {
	return &Handler{
		bulkDeactivateService: bulkDeactivateService,
	}
}

func (h *Handler) UsersBulkDeactivatePost(ctx echo.Context) error {
	var input generated.PostUsersBulkDeactivateJSONRequestBody

	if err := ctx.Bind(&input); err != nil {
		return rpc_errors.RespondBadRequest(ctx, "")
	}

	result, err := h.bulkDeactivateService.BulkDeactivateTeamUsers(ctx.Request().Context(), input.TeamName, input.UserIds)
	if err != nil {
		return rpc_errors.RespondFromError(ctx, err)
	}

	reassignments := make([]generated.Reassignment, 0, len(result.Reassignments))
	for _, r := range result.Reassignments {
		reassignments = append(reassignments, generated.Reassignment{
			PullRequestId: r.PRID,
			OldUserId:     r.OldReviewerID,
			NewUserId:     r.NewReviewerID,
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"deactivated_user_ids": result.DeactivatedUserIDs,
		"reassignments":        reassignments,
	})
}
