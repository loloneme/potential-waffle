package team_get_get

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/rpc/team/team_get_get/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_TeamGetGet(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockuserRepo(ctrl)
		mockTeamRepo := mocks.NewMockteamRepo(ctrl)
		handler := New(mockUserRepo, mockTeamRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/team/get?team_name=team-1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		params := generated.GetTeamGetParams{
			TeamName: "team-1",
		}

		mockTeamRepo.EXPECT().
			Exists(gomock.Any(), "team-1").
			Return(true, nil)

		expectedUsers := []models.User{
			{ID: "u1", Username: "user1", IsActive: true, TeamName: "team-1"},
			{ID: "u2", Username: "user2", IsActive: true, TeamName: "team-1"},
		}

		mockUserRepo.EXPECT().
			Find(gomock.Any(), gomock.Any()).
			Return(expectedUsers, nil)

		err := handler.TeamGetGet(c, params)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response generated.Team
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "team-1", response.TeamName)
		assert.Len(t, response.Members, 2)
	})

	t.Run("team not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockuserRepo(ctrl)
		mockTeamRepo := mocks.NewMockteamRepo(ctrl)
		handler := New(mockUserRepo, mockTeamRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/team/get?team_name=team-1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		params := generated.GetTeamGetParams{
			TeamName: "team-1",
		}

		mockTeamRepo.EXPECT().
			Exists(gomock.Any(), "team-1").
			Return(false, nil)

		err := handler.TeamGetGet(c, params)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, generated.NOTFOUND, response.Error.Code)
	})

	t.Run("internal error - team exists check fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockuserRepo(ctrl)
		mockTeamRepo := mocks.NewMockteamRepo(ctrl)
		handler := New(mockUserRepo, mockTeamRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/team/get?team_name=team-1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		params := generated.GetTeamGetParams{
			TeamName: "team-1",
		}

		mockTeamRepo.EXPECT().
			Exists(gomock.Any(), "team-1").
			Return(false, assert.AnError)

		err := handler.TeamGetGet(c, params)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
