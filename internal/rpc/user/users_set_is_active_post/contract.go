package users_set_is_active_post

import (
	"context"

	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/user"
)

type userRepo interface {
	UserUpdate(ctx context.Context, spec user.UpdateSpecification) (models.User, error)
}
