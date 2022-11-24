package endpoints

import (
	"context"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product/v1/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/mocks/testData"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/e2e"
)

//// we could also run the server on docker and then send rest call to the api
//func Test_Update_Product_E2E(t *testing.T) {
//	test.SkipCI(t)
//	fixture := e2e.NewE2ETestFixture()
//
//	e := NewUpdateProductEndpoint(delivery.NewProductEndpointBase(fixture.InfrastructureConfigurations, fixture.V1.ProductsGroup, fixture.Bus, fixture.CatalogsMetrics))
//	e.MapRoute()
//
//	fixture.Run()
//	defer fixture.Cleanup()
//
//	id, err := uuid.FromString("49a8e487-945b-4050-9a4c-a9242247cb48")
//	if err != nil {
//		return
//	}
//	request := updating_product.UpdateProductRequestDto{
//		Description: gofakeit.AdjectiveDescriptive(),
//		Price:       gofakeit.Price(100, 1000),
//		Name:        gofakeit.Name(),
//	}
//
//	// create httpexpect instance
//	expect := httpexpect.New(t, fixture.HttpServer.URL)
//
//	expect.PUT("/api/v1/products/{id}").
//		WithContext(context.Background()).
//		WithJSON(request).
//		WithPath("id", id.String()).
//		Expect().
//		Status(http.StatusNoContent)
//}

type updateProductE2ETest struct {
	*e2e.E2ETestFixture
	*e2e.E2ETestSharedFixture
}

func TestCreateProductE2E(t *testing.T) {
	suite.Run(t, &updateProductE2ETest{E2ETestSharedFixture: e2e.NewE2ETestSharedFixture(t)})
}

func (c *updateProductE2ETest) Test_Should_Return_NoContent_Status_With_Valid_Input() {
	testUtils.SkipCI(c.T())

	id := testData.Products[0].ProductId

	request := dtos.UpdateProductRequestDto{
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(100, 1000),
		Name:        gofakeit.Name(),
	}

	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.Cfg.Http.BasePathAddress())

	expect.PUT("products/{id}").
		WithContext(context.Background()).
		WithJSON(request).
		WithPath("id", id.String()).
		Expect().
		Status(http.StatusNoContent)
}

// Input validations
func (c *updateProductE2ETest) Test_Should_Return_Bad_Request_Status_With_Invalid_Input() {
	testUtils.SkipCI(c.T())

	id := testData.Products[0].ProductId

	request := dtos.UpdateProductRequestDto{
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       0,
		Name:        gofakeit.Name(),
	}

	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.Cfg.Http.BasePathAddress())

	expect.PUT("products/{id}").
		WithContext(context.Background()).
		WithJSON(request).
		WithPath("id", id.String()).
		Expect().
		Status(http.StatusBadRequest)
}

func (c *updateProductE2ETest) SetupTest() {
	c.T().Log("SetupTest")
	c.E2ETestFixture = e2e.NewE2ETestFixture(c.E2ETestSharedFixture)
	e := NewUpdateProductEndpoint(c.ProductEndpointBase)
	e.MapRoute()

	c.E2ETestFixture.Run()
}

func (c *updateProductE2ETest) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}
