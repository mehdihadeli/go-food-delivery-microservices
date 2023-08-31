//go:build integration
// +build integration

package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_product_by_id/v1/dtos"
	getProductByIdQuery "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_product_by_id/v1/queries"
)

func (p *productFeaturesIntegrationTestsv1) Test_Get_Product_By_Id() {
	ctx := context.Background()

	p.T().Run("Should_Returns_Existing_Product_From_DB_With_Correct_Properties", func(t *testing.T) {
		id := p.Items[0].ProductId
		query, err := getProductByIdQuery.NewGetProductById(id)
		p.Require().NoError(err)

		result, err := mediatr.Send[*getProductByIdQuery.GetProductById, *dtos.GetProductByIdResponseDto](
			ctx,
			query,
		)
		p.Require().NoError(err)
		p.NotNil(result)
		p.NotNil(result.Product)
		p.Equal(id, result.Product.ProductId)
	})

	p.T().Run("Should_Returns_NotFound_Error_When_Record_DoesNot_Exists", func(t *testing.T) {
		id := uuid.NewV4()
		query, err := getProductByIdQuery.NewGetProductById(id)
		p.Require().NoError(err)

		result, err := mediatr.Send[*getProductByIdQuery.GetProductById, *dtos.GetProductByIdResponseDto](
			ctx,
			query,
		)
		p.Require().Error(err)
		p.True(customErrors.IsNotFoundError(err))
		p.True(customErrors.IsApplicationError(err, http.StatusNotFound))
		p.Nil(result)
	})
}
