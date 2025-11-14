package create_pr

import (
	"context"

	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/user"
)

type userRepo interface {
	Find(ctx context.Context, spec user.FindSpecification) ([]models.User, error)
}

type prRepo interface {
}
