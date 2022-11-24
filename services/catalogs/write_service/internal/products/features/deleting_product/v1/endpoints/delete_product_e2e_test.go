package endpoints

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/mocks/testData"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/e2e"
)

type deleteProductE2ETest struct {
	*e2e.E2ETestFixture
	*e2e.E2ETestSharedFixture
}

func TestDeleteProductE2E(t *testing.T) {
	suite.Run(t, &deleteProductE2ETest{E2ETestSharedFixture: e2e.NewE2ETestSharedFixture(t)})
}

func (c *deleteProductE2ETest) Test_Should_Return_No_Content_Status_With_Valid_Input() {
	testUtils.SkipCI(c.T())

	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.Cfg.Http.BasePathAddress())

	id := testData.Products[0].ProductId

	expect.DELETE("products/{id}").
		WithContext(c.Ctx).
		WithPath("id", id.String()).
		Expect().
		Status(http.StatusNoContent)
}

// Input validations
func (c *deleteProductE2ETest) Test_Should_Return_Not_Found_Status_With_Invalid_Id() {
	testUtils.SkipCI(c.T())

	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.Cfg.Http.BasePathAddress())

	id := uuid.NewV4()

	expect.DELETE("products/{id}").
		WithContext(c.Ctx).
		WithPath("id", id.String()).
		Expect().
		Status(http.StatusNotFound)
}

func (c *deleteProductE2ETest) SetupTest() {
	c.T().Log("SetupTest")
	c.E2ETestFixture = e2e.NewE2ETestFixture(c.E2ETestSharedFixture)

	e := NewDeleteProductEndpoint(c.ProductEndpointBase)
	e.MapRoute()

	c.E2ETestFixture.Run()
}

func (c *deleteProductE2ETest) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}
