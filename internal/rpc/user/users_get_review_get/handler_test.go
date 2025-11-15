package users_get_review_get

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
	"github.com/loloneme/potential-waffle/internal/rpc/user/users_get_review_get/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_UsersGetReviewGet(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		// Setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockprRepo(ctrl)
		handler := New(mockRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/users/getReview?user_id=u1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		params := generated.GetUsersGetReviewParams{
			UserId: "u1",
		}

		expectedPRs := []models.PullRequest{
			{
				ID:       "pr-1001",
				Name:     "Add search",
				AuthorID: "u2",
				Status: &models.Status{
					Name: "OPEN",
				},
			},
			{
				ID:       "pr-1002",
				Name:     "Fix bug",
				AuthorID: "u3",
				Status: &models.Status{
					Name: "OPEN",
				},
			},
		}

		mockRepo.EXPECT().
			FindPullRequests(gomock.Any(), gomock.Any()).
			Return(expectedPRs, nil)

		// Execute
		err := handler.UsersGetReviewGet(c, params)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "u1", response["user_id"])
		assert.NotNil(t, response["pull_requests"])

		pullRequests := response["pull_requests"].([]interface{})
		assert.Len(t, pullRequests, 2)
	})

	t.Run("empty list", func(t *testing.T) {
		// Setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockprRepo(ctrl)
		handler := New(mockRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/users/getReview?user_id=u1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		params := generated.GetUsersGetReviewParams{
			UserId: "u1",
		}

		mockRepo.EXPECT().
			FindPullRequests(gomock.Any(), gomock.Any()).
			Return([]models.PullRequest{}, nil)

		// Execute
		err := handler.UsersGetReviewGet(c, params)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "u1", response["user_id"])

		pullRequests := response["pull_requests"].([]interface{})
		assert.Len(t, pullRequests, 0)
	})

	t.Run("internal error", func(t *testing.T) {
		// Setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockprRepo(ctrl)
		handler := New(mockRepo)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/users/getReview?user_id=u1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		params := generated.GetUsersGetReviewParams{
			UserId: "u1",
		}

		mockRepo.EXPECT().
			FindPullRequests(gomock.Any(), gomock.Any()).
			Return(nil, assert.AnError)

		// Execute
		err := handler.UsersGetReviewGet(c, params)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
