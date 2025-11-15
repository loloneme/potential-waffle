package team_add_post

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	rpc_errors "github.com/loloneme/potential-waffle/internal/rpc/errors"
	"github.com/loloneme/potential-waffle/internal/rpc/team/team_add_post/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	validTeamJSON = `{"team_name":"team-1","members":[{"user_id":"u1","username":"user1","is_active":true},{"user_id":"u2","username":"user2","is_active":true}]}`
)

func makeTestRequest(e *echo.Echo, body string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodPost, "/team/add", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func TestHandler_TeamAddPost(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		// Setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockcreateTeamService(ctrl)
		handler := New(mockService)

		e := echo.New()
		c, rec := makeTestRequest(e, validTeamJSON)

		expectedTeam := models.Team{
			TeamName: "team-1",
			Members: []models.User{
				{ID: "u1", Username: "user1", IsActive: true, TeamName: "team-1"},
				{ID: "u2", Username: "user2", IsActive: true, TeamName: "team-1"},
			},
		}

		mockService.EXPECT().
			CreateTeam(gomock.Any(), gomock.Any()).
			Return(expectedTeam, nil)

		// Execute
		err := handler.TeamAddPost(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "team")

		teamData := response["team"].(map[string]interface{})
		assert.Equal(t, "team-1", teamData["team_name"])
		assert.NotNil(t, teamData["members"])
	})

	t.Run("bad request - invalid JSON", func(t *testing.T) {
		// Setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockcreateTeamService(ctrl)
		handler := New(mockService)

		e := echo.New()
		c, rec := makeTestRequest(e, "invalid json")

		// Execute
		err := handler.TeamAddPost(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, generated.ErrorResponseErrorCode("BAD_REQUEST"), response.Error.Code)
	})

	t.Run("service error - team already exists", func(t *testing.T) {
		// Setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockcreateTeamService(ctrl)
		handler := New(mockService)

		e := echo.New()
		c, rec := makeTestRequest(e, validTeamJSON)

		mockService.EXPECT().
			CreateTeam(gomock.Any(), gomock.Any()).
			Return(models.Team{}, rpc_errors.NewTeamExists("team_name already exists"))

		// Execute
		err := handler.TeamAddPost(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, generated.TEAMEXISTS, response.Error.Code)
	})
}
