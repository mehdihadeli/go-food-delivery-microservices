//go:build integration
// +build integration

package queries

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/getting_products/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/getting_products/v1/queries"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/shared/test_fixture/integration"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
)

type getProductsIntegrationTests struct {
	*integration.IntegrationTestSharedFixture
}

func TestGetProductsIntegration(t *testing.T) {
	suite.Run(
		t,
		&getProductsIntegrationTests{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *getProductsIntegrationTests) Test_Get_Products_Query_Handler() {
	ctx := context.Background()
	query := queries.NewGetProducts(utils.NewListQuery(10, 1))
	queryResult, err := mediatr.Send[*queries.GetProducts, *dtos.GetProductsResponseDto](
		ctx,
		query,
	)

	c.NoError(err)
	c.NotNil(queryResult)
	c.NotNil(queryResult.Products)
	c.NotEmpty(queryResult.Products.Items)
}
