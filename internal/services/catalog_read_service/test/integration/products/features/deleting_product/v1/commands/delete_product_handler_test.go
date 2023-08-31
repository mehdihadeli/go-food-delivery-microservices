//go:build integration
// +build integration

package commands

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/deleting_products/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/shared/test_fixture/integration"
)

type deleteProductIntegrationTests struct {
	*integration.IntegrationTestSharedFixture
}

func TestDeleteProductIntegration(t *testing.T) {
	suite.Run(
		t,
		&deleteProductIntegrationTests{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *deleteProductIntegrationTests) Test_Delete_Product_Command_Handler() {
	productId, err := uuid.FromString(c.Items[0].ProductId)
	c.NoError(err)

	command, err := commands.NewDeleteProduct(productId)
	c.NoError(err)
	result, err := mediatr.Send[*commands.DeleteProduct, *mediatr.Unit](
		context.Background(),
		command,
	)

	c.NoError(err)
	c.NotNil(result)
}

func (c *deleteProductIntegrationTests) Test_Delete_Product_Command_Handler_Should_Return_Error_Not_Valid_UUID() {
	command, err := commands.NewDeleteProduct(uuid.UUID{})
	c.Assert().Nil(command)
	c.Require().Error(err)
}
