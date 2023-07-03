//go:build integration
// +build integration

package commands

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/creating_product/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/creating_product/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/shared/test_fixture/integration"
)

type createProductIntegrationTests struct {
	*integration.IntegrationTestSharedFixture
}

func TestCreateProductIntegration(t *testing.T) {
	suite.Run(
		t,
		&createProductIntegrationTests{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *createProductIntegrationTests) Test_Should_Create_New_Product_To_DB() {
	ctx := context.Background()
	command, err := commands.NewCreateProduct(
		uuid.NewV4().String(),
		gofakeit.Name(),
		gofakeit.AdjectiveDescriptive(),
		gofakeit.Price(150, 6000),
		time.Now(),
	)
	c.Require().NoError(err)

	result, err := mediatr.Send[*commands.CreateProduct, *dtos.CreateProductResponseDto](
		ctx,
		command,
	)
	c.Require().NoError(err)

	c.Assert().NotNil(result)
	c.Assert().Equal(command.Id, result.Id)

	createdProduct, err := c.ProductRepository.GetProductById(
		ctx,
		result.Id,
	)
	c.Require().NoError(err)
	c.Assert().NotNil(createdProduct)
}
