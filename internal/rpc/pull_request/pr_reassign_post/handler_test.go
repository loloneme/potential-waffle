package pr_reassign_post

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
	"github.com/loloneme/potential-waffle/internal/rpc/pull_request/pr_reassign_post/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	validReassignJSON = `{"pull_request_id":"pr-1001","old_user_id":"u2"}`
)

func makeTestRequest(e *echo.Echo, body string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func TestHandler_PRReassignPost(t *testing.T) {
	t.Run("successful reassign", func(t *testing.T) {
		// Setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockreassignPRService(ctrl)
		handler := New(mockService)

		e := echo.New()
		c, rec := makeTestRequest(e, validReassignJSON)

		expectedPR := models.PullRequest{
			ID:       "pr-1001",
			Name:     "Add search",
			AuthorID: "u1",
			Status: &models.Status{
				Name: "OPEN",
			},
			Reviewers: []string{"u3", "u4"},
		}

		mockService.EXPECT().
			ReassignReviewer(gomock.Any(), "pr-1001", "u2").
			Return(expectedPR, "u4", nil)

		// Execute
		err := handler.PRReassignPost(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "pr")
		assert.Contains(t, response, "replaced_by")
		assert.Equal(t, "u4", response["replaced_by"])

		prData := response["pr"].(map[string]interface{})
		assert.Equal(t, "pr-1001", prData["pull_request_id"])
	})

	t.Run("bad request - invalid JSON", func(t *testing.T) {
		// Setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockreassignPRService(ctrl)
		handler := New(mockService)

		e := echo.New()
		c, rec := makeTestRequest(e, "invalid json")

		// Execute
		err := handler.PRReassignPost(c)

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

		mockService := mocks.NewMockreassignPRService(ctrl)
		handler := New(mockService)

		e := echo.New()
		c, rec := makeTestRequest(e, validReassignJSON)

		mockService.EXPECT().
			ReassignReviewer(gomock.Any(), "pr-1001", "u2").
			Return(models.PullRequest{}, "", rpc_errors.NewNotFound("PR not found"))

		// Execute
		err := handler.PRReassignPost(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, generated.NOTFOUND, response.Error.Code)
	})

	t.Run("service error - PR merged", func(t *testing.T) {
		// Setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockreassignPRService(ctrl)
		handler := New(mockService)

		e := echo.New()
		c, rec := makeTestRequest(e, validReassignJSON)

		mockService.EXPECT().
			ReassignReviewer(gomock.Any(), "pr-1001", "u2").
			Return(models.PullRequest{}, "", rpc_errors.NewPRMerged("cannot reassign on merged PR"))

		// Execute
		err := handler.PRReassignPost(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)

		var response generated.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, generated.PRMERGED, response.Error.Code)
	})
}
