package errors

import (
	"fmt"

	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
)

// NotFoundError представляет ошибку NOT_FOUND
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

// ToErrorResponse возвращает ErrorResponse для NOT_FOUND
func (e *NotFoundError) ToErrorResponse() generated.ErrorResponse {
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.NOTFOUND
	resp.Error.Message = e.Message
	return resp
}

// NewNotFound создаёт новую ошибку NotFoundError
func NewNotFound(message string) *NotFoundError {
	return &NotFoundError{Message: message}
}

// PRExistsError представляет ошибку PR_EXISTS
type PRExistsError struct {
	Message string
}

func (e *PRExistsError) Error() string {
	return e.Message
}

// ToErrorResponse возвращает ErrorResponse для PR_EXISTS
func (e *PRExistsError) ToErrorResponse() generated.ErrorResponse {
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.PREXISTS
	resp.Error.Message = e.Message
	return resp
}

// NewPRExists создаёт новую ошибку PRExistsError
func NewPRExists(message string) *PRExistsError {
	return &PRExistsError{Message: message}
}

// PRMergedError представляет ошибку PR_MERGED
type PRMergedError struct {
	Message string
}

func (e *PRMergedError) Error() string {
	return e.Message
}

// ToErrorResponse возвращает ErrorResponse для PR_MERGED
func (e *PRMergedError) ToErrorResponse() generated.ErrorResponse {
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.PRMERGED
	resp.Error.Message = e.Message
	return resp
}

// NewPRMerged создаёт новую ошибку PRMergedError
func NewPRMerged(message string) *PRMergedError {
	return &PRMergedError{Message: message}
}

// NotAssignedError представляет ошибку NOT_ASSIGNED
type NotAssignedError struct {
	Message string
}

func (e *NotAssignedError) Error() string {
	return e.Message
}

// ToErrorResponse возвращает ErrorResponse для NOT_ASSIGNED
func (e *NotAssignedError) ToErrorResponse() generated.ErrorResponse {
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.NOTASSIGNED
	resp.Error.Message = e.Message
	return resp
}

// NewNotAssigned создаёт новую ошибку NotAssignedError
func NewNotAssigned(message string) *NotAssignedError {
	return &NotAssignedError{Message: message}
}

// NoCandidateError представляет ошибку NO_CANDIDATE
type NoCandidateError struct {
	Message string
}

func (e *NoCandidateError) Error() string {
	return e.Message
}

// ToErrorResponse возвращает ErrorResponse для NO_CANDIDATE
func (e *NoCandidateError) ToErrorResponse() generated.ErrorResponse {
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.NOCANDIDATE
	resp.Error.Message = e.Message
	return resp
}

// NewNoCandidate создаёт новую ошибку NoCandidateError
func NewNoCandidate(message string) *NoCandidateError {
	return &NoCandidateError{Message: message}
}

// TeamExistsError представляет ошибку TEAM_EXISTS
type TeamExistsError struct {
	Message string
}

func (e *TeamExistsError) Error() string {
	return e.Message
}

// ToErrorResponse возвращает ErrorResponse для TEAM_EXISTS
func (e *TeamExistsError) ToErrorResponse() generated.ErrorResponse {
	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.TEAMEXISTS
	resp.Error.Message = e.Message
	return resp
}

// NewTeamExists создаёт новую ошибку TeamExistsError
func NewTeamExists(message string) *TeamExistsError {
	return &TeamExistsError{Message: message}
}

// ErrorResponseConverter интерфейс для конвертации ошибок в ErrorResponse
type ErrorResponseConverter interface {
	ToErrorResponse() generated.ErrorResponse
}

// ToErrorResponse конвертирует ошибку в ErrorResponse
// Если ошибка реализует ErrorResponseConverter, использует его метод
// Иначе возвращает общую ошибку NOT_FOUND
func ToErrorResponse(err error) generated.ErrorResponse {
	if converter, ok := err.(ErrorResponseConverter); ok {
		return converter.ToErrorResponse()
	}

	resp := generated.ErrorResponse{}
	resp.Error.Code = generated.NOTFOUND
	resp.Error.Message = fmt.Sprintf("Internal error: %s", err.Error())
	return resp
}
