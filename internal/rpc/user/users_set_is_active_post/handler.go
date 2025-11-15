package users_set_is_active_post

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/converter"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/user"
	user_spec "github.com/loloneme/potential-waffle/internal/infrastructure/persistence/specification/user"
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
		return ctx.JSON(http.StatusBadRequest, "bad request")
	}

	spec := user_spec.NewSetIsActiveSpecification(input.UserId, input.IsActive)
	updatedUser, err := h.userRepo.UserUpdate(ctx.Request().Context(), spec)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]interface{}{
				"error": map[string]string{
					"code":    "NOT_FOUND",
					"message": "user not found",
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

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"user": converter.ToUser(updatedUser),
	})
}
