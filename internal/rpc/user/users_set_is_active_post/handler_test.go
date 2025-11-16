package users_set_is_active_post

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
	"github.com/loloneme/potential-waffle/internal/rpc/user/users_set_is_active_post/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var validSetIsActiveJSON = `{"user_id":"u1","is_active":false}`

func makeTestRequest(e *echo.Echo, body string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func TestHandler_UsersSetIsActivePost(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockuserRepo(ctrl)
		handler := New(mockRepo)

		e := echo.New()
		c, rec := makeTestRequest(e, validSetIsActiveJSON)

		expectedUser := models.User{
			ID:       "u1",
			Username: "user1",
			IsActive: false,
			TeamName: "team-1",
		}

		mockRepo.EXPECT().
			UserUpdate(gomock.Any(), gomock.Any()).
			Return(expectedUser, nil)

		err := handler.UsersSetIsActivePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "user")

		userData := response["user"].(map[string]interface{})
		assert.Equal(t, "u1", userData["user_id"])
		assert.Equal(t, false, userData["is_active"])
	})

	t.Run("bad request - invalid JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockuserRepo(ctrl)
		handler := New(mockRepo)

		e := echo.New()
		c, rec := makeTestRequest(e, "invalid json")

		err := handler.UsersSetIsActivePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, generated.ErrorResponseErrorCode("BAD_REQUEST"), response.Error.Code)
	})

	t.Run("service error - user not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockuserRepo(ctrl)
		handler := New(mockRepo)

		e := echo.New()
		c, rec := makeTestRequest(e, validSetIsActiveJSON)

		mockRepo.EXPECT().
			UserUpdate(gomock.Any(), gomock.Any()).
			Return(models.User{}, rpc_errors.NewNotFound("user not found"))

		err := handler.UsersSetIsActivePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, generated.NOTFOUND, response.Error.Code)
	})
}
