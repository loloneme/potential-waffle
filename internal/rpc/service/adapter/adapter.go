package adapter

import (
	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
)

type Adapter struct {
}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) PostPullRequestCreate(ctx echo.Context) error {
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.NOTFOUND
	resp.Error.Message = "not implemented"

	return ctx.JSON(501, resp)
}

func (a *Adapter) PostPullRequestMerge(ctx echo.Context) error {
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.NOTFOUND
	resp.Error.Message = "not implemented"

	return ctx.JSON(501, resp)
}

func (a *Adapter) PostPullRequestReassign(ctx echo.Context) error {
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.NOTFOUND
	resp.Error.Message = "not implemented"

	return ctx.JSON(501, resp)
}

func (a *Adapter) PostTeamAdd(ctx echo.Context) error {
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.NOTFOUND
	resp.Error.Message = "not implemented"

	return ctx.JSON(501, resp)
}

func (a *Adapter) GetTeamGet(ctx echo.Context, params generated.GetTeamGetParams) error {
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.NOTFOUND
	resp.Error.Message = "not implemented"

	return ctx.JSON(501, resp)
}

func (a *Adapter) GetUsersGetReview(ctx echo.Context, params generated.GetUsersGetReviewParams) error {
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.NOTFOUND
	resp.Error.Message = "not implemented"

	return ctx.JSON(501, resp)
}

func (a *Adapter) PostUsersSetIsActive(ctx echo.Context) error {
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.NOTFOUND
	resp.Error.Message = "not implemented"

	return ctx.JSON(501, resp)
}
