package v1

import (
	"context"
	"testing"

	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/mocks/testData"

	utils2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/utils"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products/dtos"
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
	utils2.SkipCI(c.T())

	query := NewGetProducts(utils.NewListQuery(10, 1))
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
