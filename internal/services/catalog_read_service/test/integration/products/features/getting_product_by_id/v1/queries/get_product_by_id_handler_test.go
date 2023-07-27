//go:build integration
// +build integration

package queries

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/get_product_by_id/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/get_product_by_id/v1/queries"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/shared/test_fixture/integration"
)

type getProductByIdIntegrationTests struct {
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

func (c *getProductByIdIntegrationTests) Test_Get_Product_By_Id_Query_Handler() {
	ctx := context.Background()
	id, err := uuid.FromString(c.Items[0].Id)
	c.Require().NoError(err)

	query := queries.NewGetProductById(id)
	result, err := mediatr.Send[*queries.GetProductById, *dtos.GetProductByIdResponseDto](
		ctx,
		query,
	)

	c.NotNil(result.Product)
	c.Equal(result.Product.Id, id.String())
}
