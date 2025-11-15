package team_get_get

import (
	"net/http"

	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/converter"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	user_spec "github.com/loloneme/potential-waffle/internal/infrastructure/persistence/specification/user"
)

type Handler struct {
	userRepo userRepo
	teamRepo teamRepo
}

func New(userRepo userRepo, teamRepo teamRepo) *Handler {
	return &Handler{
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (h *Handler) TeamGetGet(ctx echo.Context, params generated.GetTeamGetParams) error {
	teamExists, err := h.teamRepo.Exists(ctx.Request().Context(), params.TeamName)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "internal server error")
	}

	if !teamExists {
		return ctx.JSON(http.StatusNotFound, "team not found")
	}

	spec := user_spec.NewGetUsersByTeamNameSpec(params.TeamName)
	users, err := h.userRepo.Find(ctx.Request().Context(), spec)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "internal server error")
	}

	team := models.Team{
		TeamName: params.TeamName,
		Members:  users,
	}

	return ctx.JSON(http.StatusOK, converter.ToOpenAPITeam(team))
}
