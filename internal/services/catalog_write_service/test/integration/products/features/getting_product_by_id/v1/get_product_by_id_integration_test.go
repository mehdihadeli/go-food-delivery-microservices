//go:build integration
// +build integration

package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_product_by_id/v1/dtos"
	getProductByIdQuery "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_product_by_id/v1/queries"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/integration"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
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

func (c *getProductByIdIntegrationTests) Test_Should_Returns_Existing_Product_From_DB_With_Correct_Properties() {
	ctx := context.Background()

	id := c.Items[0].ProductId
	query, err := getProductByIdQuery.NewGetProductById(id)
	c.Require().NoError(err)

	result, err := mediatr.Send[*getProductByIdQuery.GetProductById, *dtos.GetProductByIdResponseDto](
		ctx,
		query,
	)

	c.Require().NoError(err)
	c.NotNil(result)
	c.NotNil(result.Product)
	c.Equal(id, result.Product.ProductId)
}

func (c *getProductByIdIntegrationTests) Test_Should_Returns_NotFound_Error_When_Record_DoesNot_Exists() {
	ctx := context.Background()

	id := uuid.NewV4()
	query, err := getProductByIdQuery.NewGetProductById(id)
	c.Require().NoError(err)

	result, err := mediatr.Send[*getProductByIdQuery.GetProductById, *dtos.GetProductByIdResponseDto](
		ctx,
		query,
	)

	c.Require().Error(err)
	c.True(customErrors.IsNotFoundError(err))
	c.True(customErrors.IsApplicationError(err, http.StatusNotFound))
	c.Nil(result)
}
