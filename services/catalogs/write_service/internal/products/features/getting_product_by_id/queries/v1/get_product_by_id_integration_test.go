package v1

import (
	"net/http"
	"testing"

	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"

	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/utils"
	gettingProductByIdDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/mocks/testData"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/integration"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type getProductByIdIntegrationTests struct {
	*integration.IntegrationTestFixture
	*integration.IntegrationTestSharedFixture
}

func TestGetProductByIdIntegration(t *testing.T) {
	suite.Run(t, &getProductByIdIntegrationTests{IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t)})
}

func (c *getProductByIdIntegrationTests) Test_Should_Delete_Product_From_DB() {
	utils.SkipCI(c.T())

	id := testData.Products[0].ProductId
	query := NewGetProductById(id)
	result, err := mediatr.Send[*GetProductById, *gettingProductByIdDtos.GetProductByIdResponseDto](c.Ctx, query)

	c.Require().NoError(err)
	c.NotNil(result)
	c.NotNil(result.Product)
	c.Equal(id, result.Product.ProductId)
}

func (c *getProductByIdIntegrationTests) Test_Should_Returns_NotFound_Error_When_Record_DoesNot_Exists() {
	utils.SkipCI(c.T())

	id := uuid.NewV4()
	query := NewGetProductById(id)
	result, err := mediatr.Send[*GetProductById, *gettingProductByIdDtos.GetProductByIdResponseDto](c.Ctx, query)

	c.Require().Error(err)
	c.True(customErrors.IsNotFoundError(err))
	c.True(customErrors.IsApplicationError(err, http.StatusNotFound))
	c.Nil(result)
}

func (c *getProductByIdIntegrationTests) SetupTest() {
	c.T().Log("SetupTest")
	c.IntegrationTestFixture = integration.NewIntegrationTestFixture(c.IntegrationTestSharedFixture)
	err := mediatr.RegisterRequestHandler[*GetProductById, *gettingProductByIdDtos.GetProductByIdResponseDto](NewGetProductByIdHandler(c.Log, c.Cfg, c.ProductRepository))
	c.Require().NoError(err)

	c.IntegrationTestFixture.Run()
}

func (c *getProductByIdIntegrationTests) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}
