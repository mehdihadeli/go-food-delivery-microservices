//go:build unit
// +build unit

package v1

import (
	"net/http"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/cqrs"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	gettingproductsv1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/gettingproducts/v1"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/gettingproducts/v1/dtos"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"

	"github.com/stretchr/testify/suite"
)

type getProductsHandlerUnitTests struct {
	*unittest.UnitTestSharedFixture
	handler cqrs.RequestHandlerWithRegisterer[*gettingproductsv1.GetProducts, *dtos.GetProductsResponseDto]
}

func TestGetProductsUnit(t *testing.T) {
	suite.Run(
		t,
		&getProductsHandlerUnitTests{
			UnitTestSharedFixture: unittest.NewUnitTestSharedFixture(t),
		},
	)
}

func (c *getProductsHandlerUnitTests) SetupTest() {
	// call base SetupTest hook before running child hook
	c.UnitTestSharedFixture.SetupTest()
	c.handler = gettingproductsv1.NewGetProductsHandler(
		fxparams.ProductHandlerParams{
			CatalogsDBContext: c.CatalogDBContext,
			Tracer:            c.Tracer,
			RabbitmqProducer:  c.Bus,
			Log:               c.Log,
		})
}

func (c *getProductsHandlerUnitTests) TearDownTest() {
	// call base TearDownTest hook before running child hook
	c.UnitTestSharedFixture.TearDownTest()
}

func (c *getProductsHandlerUnitTests) Test_Handle_Should_Return_Products_Successfully() {
	query, err := gettingproductsv1.NewGetProducts(utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	res, err := c.handler.Handle(c.Ctx, query)
	c.Require().NoError(err)
	c.NotNil(res)
	c.NotEmpty(res.Products)
	c.Equal(len(c.Products), len(res.Products.Items))
}

func (c *getProductsHandlerUnitTests) Test_Handle_Should_Return_Error_For_Mapping_List_Result() {
	query, err := gettingproductsv1.NewGetProducts(utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	mapper.ClearMappings()

	res, err := c.handler.Handle(c.Ctx, query)
	c.Require().Error(err)
	c.True(customErrors.IsApplicationError(err, http.StatusInternalServerError))
	c.Nil(res)
}
