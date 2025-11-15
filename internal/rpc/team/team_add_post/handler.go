package team_add_post

import (
	"net/http"

	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/converter"
	rpc_errors "github.com/loloneme/potential-waffle/internal/rpc/errors"
)

type Handler struct {
	createTeamService createTeamService
}

func New(createTeamService createTeamService) *Handler {
	return &Handler{
		createTeamService: createTeamService,
	}
}

func (h *Handler) TeamAddPost(ctx echo.Context) error {
	var input generated.PostTeamAddJSONRequestBody

	if err := ctx.Bind(&input); err != nil {
		return rpc_errors.RespondBadRequest(ctx, "")
	}

	teamModel := converter.ToModelTeam(input)

	createdTeam, err := h.createTeamService.CreateTeam(ctx.Request().Context(), teamModel)
	if err != nil {
		return rpc_errors.RespondFromError(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"team": converter.ToOpenAPITeam(createdTeam),
	})
}
