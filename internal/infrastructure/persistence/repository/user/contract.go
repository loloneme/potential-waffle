package user

import (
	"context"

	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

type userRepository interface {
	UpsertUsers(ctx context.Context, users []models.User) ([]models.User, error)
	GetUserByID(ctx context.Context, userID string) (models.User, error)
	Find(ctx context.Context, spec FindSpecification) ([]models.User, error)
	UserUpdate(ctx context.Context, spec UpdateSpecification) (models.User, error)
}
