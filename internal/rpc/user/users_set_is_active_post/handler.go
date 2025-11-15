package users_set_is_active_post

import (
	"net/http"

	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/converter"
	user_spec "github.com/loloneme/potential-waffle/internal/infrastructure/persistence/specification/user"
	rpc_errors "github.com/loloneme/potential-waffle/internal/rpc/errors"
)

type Handler struct {
	userRepo userRepo
}

func New(userRepo userRepo) *Handler {
	return &Handler{
		userRepo: userRepo,
	}
}

func (h *Handler) UsersSetIsActivePost(ctx echo.Context) error {
	var input generated.PostUsersSetIsActiveJSONBody

	if err := ctx.Bind(&input); err != nil {
		return rpc_errors.RespondBadRequest(ctx, "")
	}

	spec := user_spec.NewSetIsActiveSpecification(input.UserId, input.IsActive)
	updatedUser, err := h.userRepo.UserUpdate(ctx.Request().Context(), spec)
	if err != nil {
		return rpc_errors.RespondFromError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"user": converter.ToUser(updatedUser),
	})
}
