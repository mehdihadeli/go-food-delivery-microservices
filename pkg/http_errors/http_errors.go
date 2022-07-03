package http_errors

import (
	"encoding/json"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type BadRequestError struct {
	*ProblemDetail
}

// NewBadRequestError New Bad Request Error
func NewBadRequestError(detail string) ProblemDetailErr {
	badRequestError := BadRequestError{&ProblemDetail{Title: constants.ErrBadRequest, Detail: detail, Status: http.StatusBadRequest, Timestamp: time.Now()}}

	return badRequestError
}

type WrongCredentialsError struct {
	*ProblemDetail
}

// NewWrongCredentialsError New Wrong Credentials Error
func NewWrongCredentialsError(detail string) ProblemDetailErr {
	wrongCredentialError := WrongCredentialsError{&ProblemDetail{Title: constants.ErrWrongCredentials, Detail: detail, Status: http.StatusUnauthorized, Timestamp: time.Now()}}

	return wrongCredentialError
}

type NotFoundError struct {
	*ProblemDetail
}

// NewNotFoundError New Not Found Error
func NewNotFoundError(detail string) ProblemDetailErr {
	notFoundError := NotFoundError{&ProblemDetail{Title: constants.ErrNotFound, Detail: detail, Status: http.StatusNotFound, Timestamp: time.Now()}}

	return notFoundError
}

type UnauthorizedError struct {
	*ProblemDetail
}

// NewUnauthorizedError New Unauthorized Error
func NewUnauthorizedError(detail string) ProblemDetailErr {
	unAuthorizeError := UnauthorizedError{&ProblemDetail{Title: constants.ErrUnauthorized, Detail: detail, Status: http.StatusUnauthorized, Timestamp: time.Now()}}

	return unAuthorizeError
}

type ForbiddenError struct {
	*ProblemDetail
}

// NewForbiddenError New Forbidden Error
func NewForbiddenError(detail string) ProblemDetailErr {
	forbiddenError := ForbiddenError{&ProblemDetail{Title: constants.ErrForbidden, Detail: detail, Status: http.StatusForbidden, Timestamp: time.Now()}}

	return forbiddenError
}

type InternalServerError struct {
	*ProblemDetail
}

// NewInternalServerError New Internal Server Error
func NewInternalServerError(detail string) ProblemDetailErr {
	internalServerError := InternalServerError{&ProblemDetail{Title: constants.ErrInternalServerError, Detail: detail, Status: http.StatusInternalServerError, Timestamp: time.Now()}}

	return internalServerError
}

type DomainError struct {
	*ProblemDetail
}

// NewDomainError New Domain Error
func NewDomainError(status int, detail string) ProblemDetailErr {
	domainError := DomainError{&ProblemDetail{Title: constants.ErrDomain, Detail: detail, Status: status, Timestamp: time.Now()}}

	return domainError
}

type ApplicationError struct {
	*ProblemDetail
}

// NewApplicationError New Application Error
func NewApplicationError(status int, detail string) ProblemDetailErr {
	applicationError := ApplicationError{&ProblemDetail{Title: constants.ErrApplication, Detail: detail, Status: status, Timestamp: time.Now()}}

	return applicationError
}

type ApiError struct {
	*ProblemDetail
}

// NewApiError New Api Error
func NewApiError(status int, detail string) ProblemDetailErr {
	apiError := ApiError{&ProblemDetail{Title: constants.ErrApi, Detail: detail, Status: status, Timestamp: time.Now()}}

	return apiError
}

// ProblemDetailErr ProblemDetail error interface
type ProblemDetailErr interface {
	GetStatus() int
	GetTitle() string
	GetDetail() string
	Error() string
	ErrBody() error
}

// ProblemDetail error struct
type ProblemDetail struct {
	Status    int       `json:"status,omitempty"`
	Title     string    `json:"title,omitempty"`
	Detail    string    `json:"detail,omitempty"`
	Type      string    `json:"type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// ErrBody Error body
func (e ProblemDetail) ErrBody() error {
	return e
}

// Error  Error() interface method
func (e ProblemDetail) Error() string {
	return fmt.Sprintf("status: %d - title: %s - detail: %v", e.Status, e.Title, e.Detail)
}

func (e ProblemDetail) GetStatus() int {
	return e.Status
}

func (e ProblemDetail) GetTitle() string {
	return e.Title
}

func (e ProblemDetail) GetDetail() string {
	return e.Detail
}

// NewProblemDetailError New ProblemDetail Error
func NewProblemDetailError(status int, title string, detail string) ProblemDetailErr {
	restError := ProblemDetail{
		Status:    status,
		Title:     title,
		Timestamp: time.Now(),
		Detail:    detail,
	}

	return restError
}

// NewProblemDetailErrorWithMessage New ProblemDetail Error With Message
func NewProblemDetailErrorWithMessage(status int, title string) ProblemDetailErr {
	return ProblemDetail{
		Status:    status,
		Title:     title,
		Timestamp: time.Now(),
	}
}

// NewProblemDetailErrorFromBytes New ProblemDetail Error From Bytes
func NewProblemDetailErrorFromBytes(bytes []byte) (ProblemDetailErr, error) {
	var apiErr ProblemDetail
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return apiErr, nil
}
