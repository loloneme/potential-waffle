package errors

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/pull_request"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/team"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/user"
)

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func NewNotFound(message string) *NotFoundError {
	return &NotFoundError{Message: message}
}

type PRExistsError struct {
	Message string
}

func (e *PRExistsError) Error() string {
	return e.Message
}

func NewPRExists(message string) *PRExistsError {
	return &PRExistsError{Message: message}
}

type PRMergedError struct {
	Message string
}

func (e *PRMergedError) Error() string {
	return e.Message
}

func NewPRMerged(message string) *PRMergedError {
	return &PRMergedError{Message: message}
}

type NotAssignedError struct {
	Message string
}

func (e *NotAssignedError) Error() string {
	return e.Message
}

func NewNotAssigned(message string) *NotAssignedError {
	return &NotAssignedError{Message: message}
}

type NoCandidateError struct {
	Message string
}

func (e *NoCandidateError) Error() string {
	return e.Message
}

func NewNoCandidate(message string) *NoCandidateError {
	return &NoCandidateError{Message: message}
}

type TeamExistsError struct {
	Message string
}

func (e *TeamExistsError) Error() string {
	return e.Message
}

func NewTeamExists(message string) *TeamExistsError {
	return &TeamExistsError{Message: message}
}

func RespondBadRequest(ctx echo.Context, message string) error {
	if message == "" {
		message = "bad request"
	}
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.ErrorResponseErrorCode("BAD_REQUEST")
	resp.Error.Message = message
	return ctx.JSON(http.StatusBadRequest, resp)
}

func RespondInternal(ctx echo.Context, message string) error {
	if message == "" {
		message = "internal server error"
	}
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.ErrorResponseErrorCode("INTERNAL")
	resp.Error.Message = message
	return ctx.JSON(http.StatusInternalServerError, resp)
}

func RespondNotFound(ctx echo.Context, message string) error {
	if message == "" {
		message = "resource not found"
	}
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.NOTFOUND
	resp.Error.Message = message
	return ctx.JSON(http.StatusNotFound, resp)
}

func RespondPRExists(ctx echo.Context, message string) error {
	if message == "" {
		message = "PR already exists"
	}
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.PREXISTS
	resp.Error.Message = message
	return ctx.JSON(http.StatusConflict, resp)
}

func RespondPRMerged(ctx echo.Context, message string) error {
	if message == "" {
		message = "cannot reassign on merged PR"
	}
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.PRMERGED
	resp.Error.Message = message
	return ctx.JSON(http.StatusConflict, resp)
}

func RespondNotAssigned(ctx echo.Context, message string) error {
	if message == "" {
		message = "reviewer is not assigned to this PR"
	}
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.NOTASSIGNED
	resp.Error.Message = message
	return ctx.JSON(http.StatusConflict, resp)
}

func RespondNoCandidate(ctx echo.Context, message string) error {
	if message == "" {
		message = "no active replacement candidate in team"
	}
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.NOCANDIDATE
	resp.Error.Message = message
	return ctx.JSON(http.StatusConflict, resp)
}

func RespondTeamExists(ctx echo.Context, message string) error {
	if message == "" {
		message = "team_name already exists"
	}
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.TEAMEXISTS
	resp.Error.Message = message
	return ctx.JSON(http.StatusBadRequest, resp)
}

func RespondFromError(ctx echo.Context, err error) error {
	if err == nil {
		return RespondInternal(ctx, "unknown error")
	}

	switch e := err.(type) {
	case *NotFoundError:
		return RespondNotFound(ctx, e.Message)
	case *PRExistsError:
		return RespondPRExists(ctx, e.Message)
	case *PRMergedError:
		return RespondPRMerged(ctx, e.Message)
	case *NotAssignedError:
		return RespondNotAssigned(ctx, e.Message)
	case *NoCandidateError:
		return RespondNoCandidate(ctx, e.Message)
	case *TeamExistsError:
		return RespondTeamExists(ctx, e.Message)
	}

	if errors.Is(err, pull_request.ErrPRNotFound) || errors.Is(err, user.ErrNotFound) || errors.Is(err, team.ErrNotFound) {
		return RespondNotFound(ctx, "resource not found")
	}
	if errors.Is(err, pull_request.ErrPRAlreadyExists) {
		return RespondPRExists(ctx, "PR already exists")
	}
	if errors.Is(err, team.ErrAlreadyExists) {
		return RespondTeamExists(ctx, "team_name already exists")
	}

	return RespondInternal(ctx, err.Error())
}
