package httpErrors

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type DetailError struct {
	Message  string
	InnerErr error
}

type ValidationError struct {
	*ProblemDetail
}

func NewValidationError(err validator.ValidationErrors) ProblemDetailErr {
	detail := DetailError{
		Message: err.Error(),
	}
	validationError := ConflictError{&ProblemDetail{Title: constants.ErrBadRequest, Detail: detail, Status: http.StatusBadRequest, Timestamp: time.Now()}}

	return validationError
}

type ConflictError struct {
	*ProblemDetail
}

// NewConflictError New Conflict Error
func NewConflictError(err error, detailMessage string) ProblemDetailErr {
	detail := DetailError{
		Message:  detailMessage,
		InnerErr: err,
	}

	conflictError := ConflictError{&ProblemDetail{Title: constants.ErrConflict, Detail: detail, Status: http.StatusConflict, Timestamp: time.Now()}}

	return conflictError
}

type BadRequestError struct {
	*ProblemDetail
}

// NewBadRequestError New Bad Request Error
func NewBadRequestError(err error, detailMessage string) ProblemDetailErr {
	detail := DetailError{
		Message:  detailMessage,
		InnerErr: err,
	}

	badRequestError := BadRequestError{&ProblemDetail{Title: constants.ErrBadRequest, Detail: detail, Status: http.StatusBadRequest, Timestamp: time.Now()}}

	return badRequestError
}

type WrongCredentialsError struct {
	*ProblemDetail
}

// NewWrongCredentialsError New Wrong Credentials Error
func NewWrongCredentialsError(err error, detailMessage string) ProblemDetailErr {
	detail := DetailError{
		Message:  detailMessage,
		InnerErr: err,
	}

	wrongCredentialError := WrongCredentialsError{&ProblemDetail{Title: constants.ErrWrongCredentials, Detail: detail, Status: http.StatusUnauthorized, Timestamp: time.Now()}}

	return wrongCredentialError
}

type NotFoundError struct {
	*ProblemDetail
}

// NewNotFoundError New Not Found Error
func NewNotFoundError(err error, detailMessage string) ProblemDetailErr {
	detail := DetailError{
		Message:  detailMessage,
		InnerErr: err,
	}

	notFoundError := NotFoundError{&ProblemDetail{Title: constants.ErrNotFound, Detail: detail, Status: http.StatusNotFound, Timestamp: time.Now()}}

	return notFoundError
}

type UnauthorizedError struct {
	*ProblemDetail
}

// NewUnauthorizedError New Unauthorized Error
func NewUnauthorizedError(err error, detailMessage string) ProblemDetailErr {
	detail := DetailError{
		Message:  detailMessage,
		InnerErr: err,
	}

	unAuthorizeError := UnauthorizedError{&ProblemDetail{Title: constants.ErrUnauthorized, Detail: detail, Status: http.StatusUnauthorized, Timestamp: time.Now()}}

	return unAuthorizeError
}

type ForbiddenError struct {
	*ProblemDetail
}

// NewForbiddenError New Forbidden Error
func NewForbiddenError(err error, detailMessage string) ProblemDetailErr {
	detail := DetailError{
		Message:  detailMessage,
		InnerErr: err,
	}

	forbiddenError := ForbiddenError{&ProblemDetail{Title: constants.ErrForbidden, Detail: detail, Status: http.StatusForbidden, Timestamp: time.Now()}}

	return forbiddenError
}

type InternalServerError struct {
	*ProblemDetail
}

// NewInternalServerError New Internal Server Error
func NewInternalServerError(err error, detailMessage string) ProblemDetailErr {
	detail := DetailError{
		Message:  detailMessage,
		InnerErr: err,
	}

	internalServerError := InternalServerError{&ProblemDetail{Title: constants.ErrInternalServerError, Detail: detail, Status: http.StatusInternalServerError, Timestamp: time.Now()}}

	return internalServerError
}

type DomainError struct {
	*ProblemDetail
}

// NewDomainErrorWithStatus New Domain Error with specific status code
func NewDomainErrorWithStatus(err error, status int, detailMessage string) ProblemDetailErr {
	detail := DetailError{
		Message:  detailMessage,
		InnerErr: err,
	}

	domainError := DomainError{&ProblemDetail{Title: constants.ErrDomain, Detail: detail, Status: status, Timestamp: time.Now()}}

	return domainError
}

// NewDomainError New Domain Error
func NewDomainError(err error, detailMessage string) ProblemDetailErr {
	detail := DetailError{
		Message:  detailMessage,
		InnerErr: err,
	}

	domainError := DomainError{&ProblemDetail{Title: constants.ErrDomain, Detail: detail, Status: http.StatusBadRequest, Timestamp: time.Now()}}

	return domainError
}

type ApplicationError struct {
	*ProblemDetail
}

// NewApplicationErrorWithStatus New Application Error with specific status code
func NewApplicationErrorWithStatus(err error, status int, detailMessage string) ProblemDetailErr {
	detail := DetailError{
		Message:  detailMessage,
		InnerErr: err,
	}

	applicationError := ApplicationError{&ProblemDetail{Title: constants.ErrApplication, Detail: detail, Status: status, Timestamp: time.Now()}}

	return applicationError
}

// NewApplicationError New Application Error with specific status code
func NewApplicationError(err error, status int, detailMessage string) ProblemDetailErr {
	detail := DetailError{
		Message:  detailMessage,
		InnerErr: err,
	}

	applicationError := ApplicationError{&ProblemDetail{Title: constants.ErrApplication, Detail: detail, Status: http.StatusBadRequest, Timestamp: time.Now()}}

	return applicationError
}

type ApiError struct {
	*ProblemDetail
}

// NewApiError New Api Error
func NewApiError(err error, status int, detailMessage string) ProblemDetailErr {
	detail := DetailError{
		Message:  detailMessage,
		InnerErr: err,
	}

	apiError := ApiError{&ProblemDetail{Title: constants.ErrApi, Detail: detail, Status: status, Timestamp: time.Now()}}

	return apiError
}

// ProblemDetailErr ProblemDetail error interface
type ProblemDetailErr interface {
	GetStatus() int
	GetTitle() string
	GetDetail() DetailError
	GetDetailError() string
	Error() string
	ErrBody() error
}

// ProblemDetail error struct
type ProblemDetail struct {
	Status    int         `json:"status,omitempty"`
	Title     string      `json:"title,omitempty"`
	Detail    DetailError `json:"detail,omitempty"`
	Type      string      `json:"type,omitempty"`
	Timestamp time.Time   `json:"timestamp,omitempty"`
}

// ErrBody Error body
func (e *ProblemDetail) ErrBody() error {
	return e
}

// Error  Error() interface method
func (e *ProblemDetail) Error() string {
	return fmt.Sprintf("status: %d - title: %s - %s", e.Status, e.Title, e.GetDetailError())
}

func (e *ProblemDetail) GetStatus() int {
	return e.Status
}

func (e *ProblemDetail) GetTitle() string {
	return e.Title
}

func (e *ProblemDetail) GetDetail() DetailError {
	return e.Detail
}

func (e *ProblemDetail) GetDetailError() string {
	if e.Detail.InnerErr != nil {
		return fmt.Sprintf("detail message: %s  \n innerException: %s", e.Detail.Message, e.Detail.InnerErr.Error())
	}
	return fmt.Sprintf("detail message: %s", e.Detail.Message)
}

// NewProblemDetailError New ProblemDetail Error
func NewProblemDetailError(err error, status int, title string, detailMessage string) ProblemDetailErr {
	detail := DetailError{
		Message:  detailMessage,
		InnerErr: err,
	}

	problemDetail := &ProblemDetail{
		Status:    status,
		Title:     title,
		Timestamp: time.Now(),
		Detail:    detail,
	}

	return problemDetail
}

// NewProblemDetailErrorWithMessage New ProblemDetail Error With Message
func NewProblemDetailErrorWithMessage(status int, title string) ProblemDetailErr {
	return &ProblemDetail{
		Status:    status,
		Title:     title,
		Timestamp: time.Now(),
	}
}

// NewProblemDetailErrorFromBytes New ProblemDetail Error From Bytes
func NewProblemDetailErrorFromBytes(bytes []byte) (ProblemDetailErr, error) {
	var apiErr *ProblemDetail
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return apiErr, nil
}
