//go:generate mockgen -source=contract.go -destination=mocks/contract.go -package=mocks

package team_get_get

import (
	"context"

	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/user"
)

type userRepo interface {
	Find(ctx context.Context, spec user.FindSpecification) ([]models.User, error)
}

type teamRepo interface {
	Exists(ctx context.Context, teamName string) (bool, error)
}
