package statistics_get

import (
	"net/http"

	"github.com/labstack/echo/v4"
	rpc_errors "github.com/loloneme/potential-waffle/internal/rpc/errors"
)

type Handler struct {
	prRepo prRepo
}

func New(prRepo prRepo) *Handler {
	return &Handler{
		prRepo: prRepo,
	}
}

type UserAssignmentStat struct {
	UserID string `json:"user_id"`
	Count  int    `json:"count"`
}

type PRAssignmentStat struct {
	PullRequestID string `json:"pull_request_id"`
	Count         int    `json:"count"`
}

type StatisticsResponse struct {
	AssignmentsByUser []UserAssignmentStat `json:"assignments_by_user"`
	AssignmentsByPR   []PRAssignmentStat   `json:"assignments_by_pr"`
}

func (h *Handler) StatisticsGet(ctx echo.Context) error {
	stats, err := h.prRepo.GetStatistics(ctx.Request().Context())
	if err != nil {
		return rpc_errors.RespondInternal(ctx, "")
	}

	assignmentsByUser := make([]UserAssignmentStat, len(stats.AssignmentsByUser))
	for i, stat := range stats.AssignmentsByUser {
		assignmentsByUser[i] = UserAssignmentStat{
			UserID: stat.UserID,
			Count:  stat.Count,
		}
	}

	assignmentsByPR := make([]PRAssignmentStat, len(stats.AssignmentsByPR))
	for i, stat := range stats.AssignmentsByPR {
		assignmentsByPR[i] = PRAssignmentStat{
			PullRequestID: stat.PullRequestID,
			Count:         stat.Count,
		}
	}

	response := StatisticsResponse{
		AssignmentsByUser: assignmentsByUser,
		AssignmentsByPR:   assignmentsByPR,
	}

	return ctx.JSON(http.StatusOK, response)
}
