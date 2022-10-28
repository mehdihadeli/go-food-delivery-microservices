package v1

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/e2e"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type createProductE2ETest struct {
	*e2e.E2ETestFixture
	*e2e.E2ETestSharedFixture
}

func TestCreateProductE2e(t *testing.T) {
	suite.Run(t, &createProductE2ETest{E2ETestSharedFixture: e2e.NewE2ETestSharedFixture(t)})
}

func (c *createProductE2ETest) Test_Should_Return_Created_Status_With_Valid_Input() {
	utils.SkipCI(c.T())

	request := dtos.CreateProductRequestDto{
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(100, 1000),
		Name:        gofakeit.Name(),
	}

	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.Cfg.Http.BasePathAddress())

	expect.POST("products").
		WithContext(c.Ctx).
		WithJSON(request).
		Expect().
		Status(http.StatusCreated)
}

func (c *createProductE2ETest) Test_Should_Return_Bad_Request_Status_With_Invalid_Input() {
	utils.SkipCI(c.T())

	request := dtos.CreateProductRequestDto{
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       0,
		Name:        gofakeit.Name(),
	}

	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.Cfg.Http.BasePathAddress())

	expect.POST("products").
		WithContext(c.Ctx).
		WithJSON(request).
		Expect().
		Status(http.StatusBadRequest)
}

func (c *createProductE2ETest) SetupTest() {
	c.T().Log("SetupTest")
	c.E2ETestFixture = e2e.NewE2ETestFixture(c.E2ETestSharedFixture)
	e := NewCreteProductEndpoint(c.ProductEndpointBase)
	e.MapRoute()

	c.E2ETestFixture.Run()
}

func (c *createProductE2ETest) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}
