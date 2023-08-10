//go:build unit
// +build unit

package queries

import (
	"context"
	"net/http"
	"testing"

	"emperror.dev/errors"
	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_products/v1/queries"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/mocks/testData"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/unit_test"
)

type getProductsHandlerUnitTests struct {
	*unit_test.UnitTestSharedFixture
}

func TestGetProductsUnit(t *testing.T) {
	suite.Run(
		t,
		&getProductsHandlerUnitTests{UnitTestSharedFixture: unit_test.NewUnitTestSharedFixture(t)},
	)
}

func (c *getProductsHandlerUnitTests) Test_Handle_Should_Return_Products_Successfully() {
	ctx := context.Background()

	query, err := queries.NewGetProducts(utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	items := utils.NewListResult[*models.Product](
		testData.Products,
		10,
		1,
		int64(len(testData.Products)),
	)
	c.ProductRepository.On("GetAllProducts", mock.Anything, mock.Anything).
		Once().
		Return(items, nil)

	getProductsHandler := queries.NewGetProductsHandler(c.Log, c.ProductRepository, c.Tracer)

	res, err := getProductsHandler.Handle(ctx, query)
	c.Require().NoError(err)
	c.NotNil(res)
	c.NotEmpty(res.Products)
	c.Equal(len(testData.Products), len(res.Products.Items))
	c.ProductRepository.AssertNumberOfCalls(c.T(), "GetAllProducts", 1)
}

func (c *getProductsHandlerUnitTests) Test_Handle_Should_Return_Error_For_Error_In_Fetching_Data_From_Repository() {
	ctx := context.Background()

	query, err := queries.NewGetProducts(utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	c.ProductRepository.On("GetAllProducts", mock.Anything, mock.Anything).
		Once().
		Return(nil, errors.New("error in fetching products from repository"))

	getProductsHandler := queries.NewGetProductsHandler(c.Log, c.ProductRepository, c.Tracer)

	res, err := getProductsHandler.Handle(ctx, query)
	c.Require().Error(err)
	c.True(customErrors.IsApplicationError(err, http.StatusInternalServerError))
	c.Nil(res)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "GetAllProducts", 1)
}

func (c *getProductsHandlerUnitTests) Test_Handle_Should_Return_Error_For_Mapping_List_Result() {
	ctx := context.Background()

	query, err := queries.NewGetProducts(utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	c.ProductRepository.On("GetAllProducts", mock.Anything, mock.Anything).
		Once().
		Return(nil, nil)

	getProductsHandler := queries.NewGetProductsHandler(c.Log, c.ProductRepository, c.Tracer)

	res, err := getProductsHandler.Handle(ctx, query)
	c.Require().Error(err)
	c.True(customErrors.IsApplicationError(err, http.StatusInternalServerError))
	c.Nil(res)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "GetAllProducts", 1)
}
