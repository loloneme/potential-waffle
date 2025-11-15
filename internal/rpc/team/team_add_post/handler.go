package team_add_post

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/converter"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/team"
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
		return ctx.JSON(http.StatusBadRequest, "bad request")
	}

	teamModel := converter.ToModelTeam(input)

	createdTeam, err := h.createTeamService.CreateTeam(ctx.Request().Context(), teamModel)
	if err != nil {
		resp := generated.ErrorResponse{}

		switch {
		case errors.Is(err, team.ErrNotFound):
			resp.Error.Code = generated.NOTFOUND
			resp.Error.Message = "Team not found"

			return ctx.JSON(404, resp)

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
		"team": converter.ToOpenAPITeam(createdTeam),
	})
}
