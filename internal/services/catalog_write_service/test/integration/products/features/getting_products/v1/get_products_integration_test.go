//go:build integration
// +build integration

package v1

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_products/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_products/v1/queries"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/integration"

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

func (c *getProductsIntegrationTests) Test_Should_Get_Existing_Products_List_From_DB() {
	query, err := queries.NewGetProducts(utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	queryResult, err := mediatr.Send[*queries.GetProducts, *dtos.GetProductsResponseDto](
		context.Background(),
		query,
	)

	c.Require().NoError(err)
	c.NotNil(queryResult)
	c.NotNil(queryResult.Products)
	c.NotEmpty(queryResult.Products.Items)
	c.Equal(len(c.Items), len(queryResult.Products.Items))
}
