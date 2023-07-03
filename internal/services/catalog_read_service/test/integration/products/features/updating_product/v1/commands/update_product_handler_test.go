//go:build integration
// +build integration

package commands

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/updating_products/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/shared/test_fixture/integration"
)

type updateProductIntegrationTests struct {
	*integration.IntegrationTestSharedFixture
}

func TestUpdateProductIntegration(t *testing.T) {
	suite.Run(
		t,
		&updateProductIntegrationTests{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *updateProductIntegrationTests) Test_Update_Product_Command_Handler() {
	ctx := context.Background()
	productId, err := uuid.FromString(c.Items[0].ProductId)
	c.Require().NoError(err)

	command := commands.NewUpdateProduct(
		productId,
		gofakeit.Name(),
		gofakeit.AdjectiveDescriptive(),
		gofakeit.Price(150, 6000),
	)

	result, err := mediatr.Send[*commands.UpdateProduct, *mediatr.Unit](ctx, command)
	c.Require().NoError(err)
	c.NotNil(result)
}
