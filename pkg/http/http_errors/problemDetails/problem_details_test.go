package problemDetails

import (
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_Domain_Err(t *testing.T) {
	domainErr := NewDomainProblemDetail(http.StatusBadRequest, "Order with id '1' already completed", "stack")

	assert.Equal(t, "Order with id '1' already completed", domainErr.GetDetail())
	assert.Equal(t, "Domain Model Error", domainErr.GetTitle())
	assert.Equal(t, "stack", domainErr.GetStackTrace())
	assert.Equal(t, "https://httpstatuses.io/400", domainErr.GetType())
	assert.Equal(t, 400, domainErr.GetStatus())
}

func Test_Application_Err(t *testing.T) {
	applicationErr := NewApplicationProblemDetail(http.StatusBadRequest, "application error", "stack")

	assert.Equal(t, "application error", applicationErr.GetDetail())
	assert.Equal(t, "Application Service Error", applicationErr.GetTitle())
	assert.Equal(t, "stack", applicationErr.GetStackTrace())
	assert.Equal(t, "https://httpstatuses.io/400", applicationErr.GetType())
	assert.Equal(t, 400, applicationErr.GetStatus())
}

func Test_BadRequest_Err(t *testing.T) {
	badRequestError := NewBadRequestProblemDetail("bad-request error", "stack")

	assert.Equal(t, "bad-request error", badRequestError.GetDetail())
	assert.Equal(t, "Bad Request", badRequestError.GetTitle())
	assert.Equal(t, "stack", badRequestError.GetStackTrace())
	assert.Equal(t, "https://httpstatuses.io/400", badRequestError.GetType())
	assert.Equal(t, 400, badRequestError.GetStatus())
}

func Test_Parse_Error(t *testing.T) {
	// Bad-Request ProblemDetail
	badRequestError := errors.Wrap(customErrors.NewBadRequestError("bad-request error"), "bad request error")
	badRequestPrb := ParseError(badRequestError)
	assert.NotNil(t, badRequestPrb)
	assert.Equal(t, badRequestPrb.GetStatus(), 400)

	// NotFound ProblemDetail
	notFoundError := customErrors.NewNotFoundError("notfound error")
	notfoundPrb := ParseError(notFoundError)
	assert.NotNil(t, notFoundError)
	assert.Equal(t, notfoundPrb.GetStatus(), 404)
}
