package pr_create_post

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
	"github.com/loloneme/potential-waffle/internal/rpc/pull_request/pr_create_post/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var validPRJSON = `{"pull_request_id":"pr-1001","pull_request_name":"Add search","author_id":"u1"}`

func makeTestRequest(e *echo.Echo, body string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func TestHandler_PRCreatePost(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockcreatePRService(ctrl)
		handler := New(mockService)

		e := echo.New()
		c, rec := makeTestRequest(e, validPRJSON)

		expectedPR := models.PullRequest{
			ID:       "pr-1001",
			Name:     "Add search",
			AuthorID: "u1",
			Status: &models.Status{
				Name: "OPEN",
			},
			Reviewers: []string{"u2", "u3"},
		}

		mockService.EXPECT().
			CreatePR(gomock.Any(), gomock.Any()).
			Return(expectedPR, nil)

		err := handler.PRCreatePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "pr")

		prData := response["pr"].(map[string]interface{})
		assert.Equal(t, "pr-1001", prData["pull_request_id"])
		assert.Equal(t, "Add search", prData["pull_request_name"])
		assert.Equal(t, "u1", prData["author_id"])
		assert.Equal(t, "OPEN", prData["status"])
	})

	t.Run("bad request - invalid JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockcreatePRService(ctrl)
		handler := New(mockService)

		e := echo.New()
		c, rec := makeTestRequest(e, "invalid json")

		err := handler.PRCreatePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, generated.ErrorResponseErrorCode("BAD_REQUEST"), response.Error.Code)
	})

	t.Run("service error - not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockcreatePRService(ctrl)
		handler := New(mockService)

		e := echo.New()
		c, rec := makeTestRequest(e, validPRJSON)

		mockService.EXPECT().
			CreatePR(gomock.Any(), gomock.Any()).
			Return(models.PullRequest{}, rpc_errors.NewNotFound("author not found"))

		err := handler.PRCreatePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, generated.NOTFOUND, response.Error.Code)
		assert.Contains(t, response.Error.Message, "author not found")
	})

	t.Run("service error - PR already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockcreatePRService(ctrl)
		handler := New(mockService)

		e := echo.New()
		c, rec := makeTestRequest(e, validPRJSON)

		mockService.EXPECT().
			CreatePR(gomock.Any(), gomock.Any()).
			Return(models.PullRequest{}, rpc_errors.NewPRExists("PR already exists"))

		err := handler.PRCreatePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, generated.PREXISTS, response.Error.Code)
		assert.Contains(t, response.Error.Message, "PR already exists")
	})
}
