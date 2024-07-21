//go:build unit
// +build unit

package v1

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/cqrs"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm/gormdbcontext"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/data/datamodels"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	updatingoroductsv1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/updatingproduct/v1"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"

	"emperror.dev/errors"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type updateProductHandlerUnitTests struct {
	*unittest.UnitTestSharedFixture
	handler cqrs.RequestHandlerWithRegisterer[*updatingoroductsv1.UpdateProduct, *mediatr.Unit]
}

func TestUpdateProductHandlerUnit(t *testing.T) {
	suite.Run(
		t,
		&updateProductHandlerUnitTests{
			UnitTestSharedFixture: unittest.NewUnitTestSharedFixture(t),
		},
	)
}

func (c *updateProductHandlerUnitTests) SetupTest() {
	// call base `SetupTest hook` before running child hook
	c.UnitTestSharedFixture.SetupTest()
	c.handler = updatingoroductsv1.NewUpdateProductHandler(
		fxparams.ProductHandlerParams{
			CatalogsDBContext: c.CatalogDBContext,
			Tracer:            c.Tracer,
			RabbitmqProducer:  c.Bus,
			Log:               c.Log,
		},
	)
}

func (c *updateProductHandlerUnitTests) TearDownTest() {
	// call base `TearDownTest hook` before running child hook
	c.UnitTestSharedFixture.TearDownTest()
}

func (c *updateProductHandlerUnitTests) Test_Handle_Should_Update_Product_With_Valid_Data() {
	existing := c.Products[0]

	updateProductCommand, err := updatingoroductsv1.NewUpdateProduct(
		existing.Id,
		gofakeit.Name(),
		gofakeit.EmojiDescription(),
		existing.Price,
	)
	c.Require().NoError(err)

	c.BeginTx()
	_, err = c.handler.Handle(c.Ctx, updateProductCommand)
	c.CommitTx()

	c.Require().NoError(err)

	updatedProduct, err := gormdbcontext.FindDataModelByID[*datamodels.ProductDataModel](
		c.Ctx,
		c.CatalogDBContext,
		updateProductCommand.ProductID,
	)
	c.Require().NoError(err)

	c.Assert().Equal(updatedProduct.Id, updateProductCommand.ProductID)
	c.Assert().Equal(updatedProduct.Name, updateProductCommand.Name)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
}

func (c *updateProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_NotFound_Item() {
	id := uuid.NewV4()

	command, err := updatingoroductsv1.NewUpdateProduct(
		id,
		gofakeit.Name(),
		gofakeit.EmojiDescription(),
		gofakeit.Price(150, 6000),
	)
	c.Require().NoError(err)

	c.BeginTx()
	_, err = c.handler.Handle(c.Ctx, command)
	c.CommitTx()

	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 0)
	c.True(customErrors.IsApplicationError(err, http.StatusNotFound))
	c.ErrorContains(
		err,
		fmt.Sprintf("product with id `%s` not found", id.String()),
	)
}

func (c *updateProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_Error_In_Bus() {
	existing := c.Products[0]

	updateProductCommand, err := updatingoroductsv1.NewUpdateProduct(
		existing.Id,
		gofakeit.Name(),
		gofakeit.EmojiDescription(),
		existing.Price,
	)
	c.Require().NoError(err)

	// override called mock
	// https://github.com/stretchr/testify/issues/558
	c.Bus.Mock.ExpectedCalls = nil
	c.Bus.On("PublishMessage", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return(errors.New("error in the publish message"))

	c.BeginTx()
	_, err = c.handler.Handle(c.Ctx, updateProductCommand)
	c.CommitTx()

	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
	c.ErrorContains(err, "error in the publish message")
	c.ErrorContains(err, "error in publishing 'ProductUpdated' message")
}
