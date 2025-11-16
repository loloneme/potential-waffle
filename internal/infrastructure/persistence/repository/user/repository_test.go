package user_test

import (
	"context"
	"errors"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/user"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/user/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserRepository_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockuserRepository(ctrl)
	ctx := context.Background()
	userID := "user-1"

	t.Run("successful get", func(t *testing.T) {
		expectedUser := models.User{
			ID:       "user-1",
			Username: "testuser",
			IsActive: true,
			TeamName: "team-1",
		}

		mockRepo.EXPECT().
			GetUserByID(gomock.Any(), userID).
			Return(expectedUser, nil)

		u, err := mockRepo.GetUserByID(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, u.ID)
		assert.Equal(t, expectedUser.Username, u.Username)
		assert.Equal(t, expectedUser.IsActive, u.IsActive)
		assert.Equal(t, expectedUser.TeamName, u.TeamName)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.EXPECT().
			GetUserByID(gomock.Any(), userID).
			Return(models.User{}, user.ErrNotFound)

		u, err := mockRepo.GetUserByID(ctx, userID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, user.ErrNotFound))
		assert.Equal(t, models.User{}, u)
	})
}

func TestUserRepository_GetUserTeamName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockuserRepository(ctrl)
	ctx := context.Background()
	userID := "user-1"

	t.Run("successful get team name", func(t *testing.T) {
		expectedTeamName := "team-1"

		mockRepo.EXPECT().
			GetUserTeamName(gomock.Any(), userID).
			Return(expectedTeamName, nil)

		teamName, err := mockRepo.GetUserTeamName(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedTeamName, teamName)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.EXPECT().
			GetUserTeamName(gomock.Any(), userID).
			Return("", user.ErrNotFound)

		teamName, err := mockRepo.GetUserTeamName(ctx, userID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, user.ErrNotFound))
		assert.Empty(t, teamName)
	})
}

func TestUserRepository_UpsertUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockuserRepository(ctrl)
	ctx := context.Background()
	var tx *sqlx.Tx

	t.Run("successful upsert", func(t *testing.T) {
		users := []models.User{
			{
				ID:       "user-1",
				Username: "testuser1",
				IsActive: true,
				TeamName: "team-1",
			},
			{
				ID:       "user-2",
				Username: "testuser2",
				IsActive: false,
				TeamName: "team-1",
			},
		}
		expectedUsers := []models.User{
			{
				ID:       "user-1",
				Username: "testuser1",
				IsActive: true,
				TeamName: "team-1",
			},
			{
				ID:       "user-2",
				Username: "testuser2",
				IsActive: false,
				TeamName: "team-1",
			},
		}

		mockRepo.EXPECT().
			UpsertUsers(gomock.Any(), gomock.Any(), users).
			Return(expectedUsers, nil)

		result, err := mockRepo.UpsertUsers(ctx, tx, users)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedUsers), len(result))
		assert.Equal(t, expectedUsers[0].ID, result[0].ID)
	})

	t.Run("empty users list", func(t *testing.T) {
		mockRepo.EXPECT().
			UpsertUsers(gomock.Any(), gomock.Any(), []models.User{}).
			Return([]models.User{}, nil)

		result, err := mockRepo.UpsertUsers(ctx, tx, []models.User{})
		assert.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestUserRepository_Find(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockuserRepository(ctrl)
	ctx := context.Background()

	t.Run("successful find", func(t *testing.T) {
		expectedUsers := []models.User{
			{
				ID:       "user-1",
				Username: "testuser1",
				IsActive: true,
				TeamName: "team-1",
			},
			{
				ID:       "user-2",
				Username: "testuser2",
				IsActive: false,
				TeamName: "team-1",
			},
		}
		spec := &mockFindSpec{fields: []string{"user_id", "username", "is_active", "team_name"}}

		mockRepo.EXPECT().
			Find(gomock.Any(), spec).
			Return(expectedUsers, nil)

		users, err := mockRepo.Find(ctx, spec)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedUsers), len(users))
		assert.Equal(t, expectedUsers[0].ID, users[0].ID)
	})

	t.Run("no users found", func(t *testing.T) {
		spec := &mockFindSpec{fields: []string{"user_id", "username", "is_active", "team_name"}}

		mockRepo.EXPECT().
			Find(gomock.Any(), spec).
			Return([]models.User{}, user.ErrNotFound)

		users, err := mockRepo.Find(ctx, spec)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, user.ErrNotFound))
		assert.Empty(t, users)
	})
}

func TestUserRepository_UserUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockuserRepository(ctrl)
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		expectedUser := models.User{
			ID:       "user-1",
			Username: "testuser",
			IsActive: false,
			TeamName: "team-1",
		}
		spec := &mockUpdateSpec{
			setValues:       map[string]interface{}{"is_active": false},
			returningFields: []string{"user_id", "username", "is_active", "team_name"},
		}

		mockRepo.EXPECT().
			UserUpdate(gomock.Any(), spec).
			Return(expectedUser, nil)

		u, err := mockRepo.UserUpdate(ctx, spec)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, u.ID)
		assert.Equal(t, expectedUser.IsActive, u.IsActive)
	})

	t.Run("user not found", func(t *testing.T) {
		spec := &mockUpdateSpec{
			setValues:       map[string]interface{}{"is_active": false},
			returningFields: []string{"user_id", "username", "is_active", "team_name"},
		}

		mockRepo.EXPECT().
			UserUpdate(gomock.Any(), spec).
			Return(models.User{}, user.ErrNotFound)

		u, err := mockRepo.UserUpdate(ctx, spec)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, user.ErrNotFound))
		assert.Equal(t, models.User{}, u)
	})
}

type mockFindSpec struct {
	fields []string
}

func (m *mockFindSpec) GetFields() []string {
	return m.fields
}

func (m *mockFindSpec) GetRule(s sq.SelectBuilder) sq.SelectBuilder {
	return s
}

type mockUpdateSpec struct {
	setValues       map[string]interface{}
	returningFields []string
}

func (m *mockUpdateSpec) GetSetValues() map[string]interface{} {
	return m.setValues
}

func (m *mockUpdateSpec) GetRule(builder sq.UpdateBuilder) sq.UpdateBuilder {
	return builder
}

func (m *mockUpdateSpec) GetReturningFields() []string {
	return m.returningFields
}
