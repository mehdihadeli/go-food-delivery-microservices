//go:build.sh e2e
// +build.sh e2e

package endpoints

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/mocks/testData"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/test_fixtures/e2e"
)

type getProductByIdE2ETest struct {
	*e2e.E2ETestFixture
	*e2e.E2ETestSharedFixture
}

func TestGetProductByIdE2E(t *testing.T) {
	suite.Run(t, &getProductByIdE2ETest{E2ETestSharedFixture: e2e.NewE2ETestSharedFixture(t)})
}

func (c *getProductByIdE2ETest) Test_Should_Return_Ok_Status_With_Valid_Id() {
	testUtils.SkipCI(c.T())

	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.Cfg.Http.BasePathAddress())

	id := testData.Products[0].ProductId

	expect.GET("products/{id}").
		WithPath("id", id.String()).
		WithContext(c.Ctx).
		Expect().
		Status(http.StatusOK)
}

// Input validations
func (c *getProductByIdE2ETest) Test_Should_Return_NotFound_Status_With_Invalid_Id() {
	testUtils.SkipCI(c.T())

	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.Cfg.Http.BasePathAddress())

	id := uuid.NewV4()

	expect.GET("products/{id}").
		WithPath("id", id.String()).
		WithContext(c.Ctx).
		Expect().
		Status(http.StatusNotFound)
}

func (c *getProductByIdE2ETest) SetupTest() {
	c.T().Log("SetupTest")
	c.E2ETestFixture = e2e.NewE2ETestFixture(c.E2ETestSharedFixture)
	e := NewGetProductByIdEndpoint(c.ProductEndpointBase)
	e.MapRoute()

	c.E2ETestFixture.Run()
}

func (c *getProductByIdE2ETest) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}
