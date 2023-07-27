//go:build.sh integration
// +build.sh integration

package getProductByIdQuery

import (
	"net/http"
	"testing"

	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/features/getting_product_by_id/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/mocks/testData"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/test_fixtures/integration"
)

type getProductByIdIntegrationTests struct {
	*integration.IntegrationTestFixture
	*integration.IntegrationTestSharedFixture
}

func TestGetProductByIdIntegration(t *testing.T) {
	suite.Run(
		t,
		&getProductByIdIntegrationTests{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *getProductByIdIntegrationTests) Test_Should_Returns_Existing_Product_From_DB_With_Correct_Properties() {
	testUtils.SkipCI(c.T())

	id := testData.Products[0].ProductId
	query, err := NewGetProductById(id)
	c.Require().NoError(err)

	result, err := mediatr.Send[*GetProductById, *dtos.GetProductByIdResponseDto](c.Ctx, query)

	c.Require().NoError(err)
	c.NotNil(result)
	c.NotNil(result.Product)
	c.Equal(id, result.Product.ProductId)
}

func (c *getProductByIdIntegrationTests) Test_Should_Returns_NotFound_Error_When_Record_DoesNot_Exists() {
	testUtils.SkipCI(c.T())

	id := uuid.NewV4()
	query, err := NewGetProductById(id)
	c.Require().NoError(err)

	result, err := mediatr.Send[*GetProductById, *dtos.GetProductByIdResponseDto](c.Ctx, query)

	c.Require().Error(err)
	c.True(customErrors.IsNotFoundError(err))
	c.True(customErrors.IsApplicationError(err, http.StatusNotFound))
	c.Nil(result)
}

func (c *getProductByIdIntegrationTests) SetupTest() {
	c.T().Log("SetupTest")
	c.IntegrationTestFixture = integration.NewIntegrationTestFixture(c.IntegrationTestSharedFixture)
	err := mediatr.RegisterRequestHandler[*GetProductById, *dtos.GetProductByIdResponseDto](
		NewGetProductByIdHandler(c.Log, c.Cfg, c.ProductRepository),
	)
	c.Require().NoError(err)

	c.IntegrationTestFixture.Run()
}

func (c *getProductByIdIntegrationTests) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}
