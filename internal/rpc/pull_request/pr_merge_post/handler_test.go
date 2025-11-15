package pr_merge_post

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
	"github.com/loloneme/potential-waffle/internal/rpc/pull_request/pr_merge_post/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	validMergeJSON = `{"pull_request_id":"pr-1001"}`
)

func makeTestRequest(e *echo.Echo, body string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func TestHandler_PRMergePost(t *testing.T) {
	t.Run("successful merge", func(t *testing.T) {
		// Setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockMergePRService(ctrl)
		handler := New(mockService)

		e := echo.New()
		c, rec := makeTestRequest(e, validMergeJSON)

		expectedPR := models.PullRequest{
			ID:       "pr-1001",
			Name:     "Add search",
			AuthorID: "u1",
			Status: &models.Status{
				Name: "MERGED",
			},
		}

		mockService.EXPECT().
			MergePullRequest(gomock.Any(), "pr-1001", "MERGED").
			Return(expectedPR, nil)

		// Execute
		err := handler.PRMergePost(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "pr")

		prData := response["pr"].(map[string]interface{})
		assert.Equal(t, "pr-1001", prData["pull_request_id"])
		assert.Equal(t, "MERGED", prData["status"])
	})

	t.Run("bad request - invalid JSON", func(t *testing.T) {
		// Setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockMergePRService(ctrl)
		handler := New(mockService)

		e := echo.New()
		c, rec := makeTestRequest(e, "invalid json")

		// Execute
		err := handler.PRMergePost(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, generated.ErrorResponseErrorCode("BAD_REQUEST"), response.Error.Code)
	})

	t.Run("service error - PR not found", func(t *testing.T) {
		// Setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockMergePRService(ctrl)
		handler := New(mockService)

		e := echo.New()
		c, rec := makeTestRequest(e, validMergeJSON)

		mockService.EXPECT().
			MergePullRequest(gomock.Any(), "pr-1001", "MERGED").
			Return(models.PullRequest{}, rpc_errors.NewNotFound("PR not found"))

		// Execute
		err := handler.PRMergePost(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, generated.NOTFOUND, response.Error.Code)
	})
}
