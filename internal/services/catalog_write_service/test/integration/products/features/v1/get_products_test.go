//go:build integration
// +build integration

package v1

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-mediatr"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_products/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_products/v1/queries"
)

func (p *productFeaturesIntegrationTestsv1) Test_Get_Products() {
	p.T().Run("Should_Get_Existing_Products_List_From_DB", func(t *testing.T) {
		query, err := queries.NewGetProducts(utils.NewListQuery(10, 1))
		p.Require().NoError(err)

		queryResult, err := mediatr.Send[*queries.GetProducts, *dtos.GetProductsResponseDto](
			context.Background(),
			query,
		)

		p.Require().NoError(err)
		p.NotNil(queryResult)
		p.NotNil(queryResult.Products)
		p.NotEmpty(queryResult.Products.Items)
		p.Equal(len(p.Items), len(queryResult.Products.Items))
	})
}
