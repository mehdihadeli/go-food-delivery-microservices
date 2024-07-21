package v1

import (
	"net/http"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/cqrs"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	searchingproductsv1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/searchingproduct/v1"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/searchingproduct/v1/dtos"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"

	"github.com/stretchr/testify/suite"
)

type searchProductsHandlerUnitTests struct {
	*unittest.UnitTestSharedFixture
	handler cqrs.RequestHandlerWithRegisterer[*searchingproductsv1.SearchProducts, *dtos.SearchProductsResponseDto]
}

func TestSearchProductsUnit(t *testing.T) {
	suite.Run(
		t,
		&searchProductsHandlerUnitTests{
			UnitTestSharedFixture: unittest.NewUnitTestSharedFixture(t),
		},
	)
}

func (c *searchProductsHandlerUnitTests) SetupTest() {
	// call base SetupTest hook before running child hook
	c.UnitTestSharedFixture.SetupTest()
	c.handler = searchingproductsv1.NewSearchProductsHandler(
		fxparams.ProductHandlerParams{
			CatalogsDBContext: c.CatalogDBContext,
			Tracer:            c.Tracer,
			RabbitmqProducer:  c.Bus,
			Log:               c.Log,
		})
}

func (c *searchProductsHandlerUnitTests) TearDownTest() {
	// call base TearDownTest hook before running child hook
	c.UnitTestSharedFixture.TearDownTest()
}

func (c *searchProductsHandlerUnitTests) Test_Handle_Should_Return_Products_Successfully() {
	query, err := searchingproductsv1.NewSearchProducts(
		c.Products[0].Name,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	res, err := c.handler.Handle(c.Ctx, query)
	c.Require().NoError(err)
	c.NotNil(res)
	c.NotEmpty(res.Products)
	c.Equal(len(res.Products.Items), 1)
}

func (c *searchProductsHandlerUnitTests) Test_Handle_Should_Return_Error_For_Mapping_List_Result() {
	query, err := searchingproductsv1.NewSearchProducts(
		c.Products[0].Name,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	mapper.ClearMappings()

	res, err := c.handler.Handle(c.Ctx, query)
	c.Require().Error(err)
	c.True(customErrors.IsApplicationError(err, http.StatusInternalServerError))
	c.Nil(res)
}
