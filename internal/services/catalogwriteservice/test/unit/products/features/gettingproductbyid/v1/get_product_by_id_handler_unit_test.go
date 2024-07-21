//go:build unit
// +build unit

package v1

import (
	"fmt"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/cqrs"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	gettingproductbyidv1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/gettingproductbyid/v1"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/gettingproductbyid/v1/dtos"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type getProductByIdHandlerTest struct {
	*unittest.UnitTestSharedFixture
	handler cqrs.RequestHandlerWithRegisterer[*gettingproductbyidv1.GetProductById, *dtos.GetProductByIdResponseDto]
}

func TestGetProductByIdHandlerUnit(t *testing.T) {
	suite.Run(t, &getProductByIdHandlerTest{
		UnitTestSharedFixture: unittest.NewUnitTestSharedFixture(t),
	})
}

func (c *getProductByIdHandlerTest) SetupTest() {
	// call base SetupTest hook before running child hook
	c.UnitTestSharedFixture.SetupTest()
	c.handler = gettingproductbyidv1.NewGetProductByIDHandler(
		fxparams.ProductHandlerParams{
			CatalogsDBContext: c.CatalogDBContext,
			Tracer:            c.Tracer,
			RabbitmqProducer:  c.Bus,
			Log:               c.Log,
		})
}

func (c *getProductByIdHandlerTest) TearDownTest() {
	// call base TearDownTest hook before running child hook
	c.UnitTestSharedFixture.TearDownTest()
}

func (c *getProductByIdHandlerTest) Test_Handle_Should_Return_Correct_Product_By_ID() {
	product := c.Products[0]

	query, err := gettingproductbyidv1.NewGetProductById(product.Id)
	c.Require().NoError(err)

	dto, err := c.handler.Handle(c.Ctx, query)
	c.Require().NoError(err)
	c.Assert().NotNil(dto)
	c.Assert().NotNil(dto.Product)
	c.Assert().Equal(dto.Product.Id, product.Id)
	c.Assert().Equal(dto.Product.Name, product.Name)
}

func (c *getProductByIdHandlerTest) Test_Handle_Should_Return_NotFound_Error_For_NotFound_Item() {
	id := uuid.NewV4()

	query, err := gettingproductbyidv1.NewGetProductById(id)
	c.Require().NoError(err)

	dto, err := c.handler.Handle(c.Ctx, query)
	c.Require().Error(err)
	c.True(customErrors.IsNotFoundError(err))
	c.ErrorContains(
		err,
		fmt.Sprintf(
			"product_data_model with id `%s` not found in the database",
			id.String(),
		),
	)
	c.Nil(dto)
}

func (c *getProductByIdHandlerTest) Test_Handle_Should_Return_Error_For_Error_In_Mapping() {
	mapper.ClearMappings()

	product := c.Products[0]

	query, err := gettingproductbyidv1.NewGetProductById(product.Id)
	c.Require().NoError(err)

	dto, err := c.handler.Handle(c.Ctx, query)

	c.Nil(dto)
	c.Require().Error(err)
	c.True(customErrors.IsInternalServerError(err))
	c.ErrorContains(err, "error in the mapping product")
}
