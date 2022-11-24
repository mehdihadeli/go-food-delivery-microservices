package queries

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/suite"

	testUtils "github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products/v1/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/mocks/testData"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/integration"
)

type getProductsIntegrationTests struct {
	*integration.IntegrationTestFixture
	*integration.IntegrationTestSharedFixture
}

func TestGetProductsIntegration(t *testing.T) {
	suite.Run(t, &getProductsIntegrationTests{IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t)})
}

func (c *getProductsIntegrationTests) Test_Should_Delete_Product_From_DB() {
	testUtils.SkipCI(c.T())

	query, err := NewGetProducts(utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	queryResult, err := mediatr.Send[*GetProducts, *dtos.GetProductsResponseDto](context.Background(), query)

	c.Require().NoError(err)
	c.NotNil(queryResult)
	c.NotNil(queryResult.Products)
	c.NotEmpty(queryResult.Products.Items)
	c.Equal(len(testData.Products), len(queryResult.Products.Items))
}

func (c *getProductsIntegrationTests) SetupTest() {
	c.T().Log("SetupTest")
	c.IntegrationTestFixture = integration.NewIntegrationTestFixture(c.IntegrationTestSharedFixture)
	err := mediatr.RegisterRequestHandler[*GetProducts, *dtos.GetProductsResponseDto](NewGetProductsHandler(c.Log, c.Cfg, c.ProductRepository))
	c.Require().NoError(err)

	c.IntegrationTestFixture.Run()
}

func (c *getProductsIntegrationTests) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}
