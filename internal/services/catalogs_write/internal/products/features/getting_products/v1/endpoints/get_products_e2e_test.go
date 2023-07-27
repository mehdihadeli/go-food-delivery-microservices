//go:build.sh e2e
// +build.sh e2e

package endpoints

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"

	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/test_fixtures/e2e"
)

type getProductsE2ETest struct {
	*e2e.E2ETestFixture
	*e2e.E2ETestSharedFixture
}

func TestGetProductsE2e(t *testing.T) {
	suite.Run(t, &getProductsE2ETest{E2ETestSharedFixture: e2e.NewE2ETestSharedFixture(t)})
}

func (c *getProductsE2ETest) Test_Should_Return_Ok_Status() {
	testUtils.SkipCI(c.T())

	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.Cfg.Http.BasePathAddress())

	expect.GET("products").
		WithContext(c.Ctx).
		Expect().
		Status(http.StatusOK)
}

func (c *getProductsE2ETest) SetupTest() {
	c.T().Log("SetupTest")
	c.E2ETestFixture = e2e.NewE2ETestFixture(c.E2ETestSharedFixture)
	e := NewGetProductsEndpoint(c.ProductEndpointBase)
	e.MapRoute()

	c.E2ETestFixture.Run()
}

func (c *getProductsE2ETest) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}
