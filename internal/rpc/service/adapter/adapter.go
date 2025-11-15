package adapter

import (
	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/rpc/pull_request/pr_create_post"
	"github.com/loloneme/potential-waffle/internal/rpc/pull_request/pr_merge_post"
	"github.com/loloneme/potential-waffle/internal/rpc/pull_request/pr_reassign_post"
	"github.com/loloneme/potential-waffle/internal/rpc/team/team_add_post"
	"github.com/loloneme/potential-waffle/internal/rpc/team/team_get_get"
	"github.com/loloneme/potential-waffle/internal/rpc/user/users_get_review_get"
	"github.com/loloneme/potential-waffle/internal/rpc/user/users_set_is_active_post"
)

type Adapter struct {
	createTeamHandler *team_add_post.Handler
	getTeamHandler    *team_get_get.Handler

	createPullRequestHandler   *pr_create_post.Handler
	mergePullRequestHandler    *pr_merge_post.Handler
	reassignPullRequestHandler *pr_reassign_post.Handler

	getUsersReviewHandler *users_get_review_get.Handler
	setIsActiveHandler    *users_set_is_active_post.Handler
}

func NewAdapter(
	createTeamHandler *team_add_post.Handler,
	getTeamHandler *team_get_get.Handler,
	createPullRequestHandler *pr_create_post.Handler,
	mergePullRequestHandler *pr_merge_post.Handler,
	reassignPullRequestHandler *pr_reassign_post.Handler,
	getUsersReviewHandler *users_get_review_get.Handler,
	setIsActiveHandler *users_set_is_active_post.Handler,
) *Adapter {
	return &Adapter{
		createTeamHandler:          createTeamHandler,
		getTeamHandler:             getTeamHandler,
		createPullRequestHandler:   createPullRequestHandler,
		mergePullRequestHandler:    mergePullRequestHandler,
		reassignPullRequestHandler: reassignPullRequestHandler,
		getUsersReviewHandler:      getUsersReviewHandler,
		setIsActiveHandler:         setIsActiveHandler,
	}
}

func (a *Adapter) PostPullRequestCreate(ctx echo.Context) error {
	return a.createPullRequestHandler.PRCreatePost(ctx)
}

func (a *Adapter) PostPullRequestMerge(ctx echo.Context) error {
	return a.mergePullRequestHandler.PRMergePost(ctx)
}

func (a *Adapter) PostPullRequestReassign(ctx echo.Context) error {
	return a.reassignPullRequestHandler.PRReassignPost(ctx)
}

func (a *Adapter) PostTeamAdd(ctx echo.Context) error {
	return a.createTeamHandler.TeamAddPost(ctx)
}

func (a *Adapter) GetTeamGet(ctx echo.Context, params generated.GetTeamGetParams) error {
	return a.getTeamHandler.TeamGetGet(ctx, params)
}

func (a *Adapter) GetUsersGetReview(ctx echo.Context, params generated.GetUsersGetReviewParams) error {
	return a.getUsersReviewHandler.UsersGetReviewGet(ctx, params)
}

func (a *Adapter) PostUsersSetIsActive(ctx echo.Context) error {
	return a.setIsActiveHandler.UsersSetIsActivePost(ctx)
}
