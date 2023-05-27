package commands

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

    testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"
    "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/features/creating_product/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/test_fixture/integration"
)

type createProductIntegrationTests struct {
	*integration.IntegrationTestFixture
	*integration.IntegrationTestSharedFixture
}

func TestCreateProductIntegration(t *testing.T) {
	suite.Run(t, &createProductIntegrationTests{IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t)})
}

func (c *createProductIntegrationTests) Test_Should_Create_New_Product_To_DB() {
	testUtils.SkipCI(c.T())

	command, err := NewCreateProduct(uuid.NewV4().String(), gofakeit.Name(), gofakeit.AdjectiveDescriptive(), gofakeit.Price(150, 6000), time.Now())
	c.Require().NoError(err)

	result, err := mediatr.Send[*CreateProduct, *dtos.CreateProductResponseDto](c.Ctx, command)
	c.Require().NoError(err)

	c.Assert().NotNil(result)
	c.Assert().Equal(command.Id, result.Id)

	createdProduct, err := c.IntegrationTestFixture.MongoProductRepository.GetProductById(c.Ctx, result.Id)
	c.Require().NoError(err)
	c.Assert().NotNil(createdProduct)
}

//func Test_Create_Product_Command_Handler(t *testing.T) {
//	utils.SkipCI(t)
//	fixture := integration_events.NewIntegrationTestFixture()
//
//	err := mediatr.RegisterRequestHandler[*CreateProduct, *creating_product.CreateProductResponseDto](NewCreateProductHandler(fixture.Log, fixture.Cfg, fixture.MongoProductRepository, fixture.RedisProductRepository))
//	assert.NoError(t, err)
//
//	fixture.Run()
//	defer fixture.Cleanup()
//
//	command := NewCreateProduct(gofakeit.UUID(), gofakeit.Name(), gofakeit.AdjectiveDescriptive(), gofakeit.Price(150, 6000), time.Now())
//	result, err := mediatr.Send[*CreateProduct, *creating_product.CreateProductResponseDto](context.Background(), command)
//	assert.NoError(t, err)
//
//	assert.NotNil(t, result)
//	assert.NotEmpty(t, result.Id)
//}
