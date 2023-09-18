//go:build unit
// +build unit

package commands

import (
	"context"
	"net/http"
	"testing"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/deleting_product/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/unit_test"

	"emperror.dev/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type deleteProductHandlerUnitTests struct {
	*unit_test.UnitTestSharedFixture
}

func TestDeleteProductHandlerUnit(t *testing.T) {
	suite.Run(
		t,
		&deleteProductHandlerUnitTests{
			UnitTestSharedFixture: unit_test.NewUnitTestSharedFixture(t),
		},
	)
}

func (c *deleteProductHandlerUnitTests) Test_Handle_Should_Delete_Product_With_Valid_Id() {
	ctx := context.Background()
	id := c.Items[0].ProductId

	deleteProduct := &commands.DeleteProduct{
		ProductID: id,
	}

	c.ProductRepository.On("DeleteProductByID", mock.Anything, id).
		Once().
		Return(nil)

	deleteProductHandler := commands.NewDeleteProductHandler(c.Log, c.Uow, c.Bus, c.Tracer)

	_, err := deleteProductHandler.Handle(ctx, deleteProduct)
	c.Require().NoError(err)

	c.Uow.AssertNumberOfCalls(c.T(), "Do", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "DeleteProductByID", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
}

func (c *deleteProductHandlerUnitTests) Test_Handle_Should_Return_NotFound_Error_When_Id_Is_Invalid() {
	ctx := context.Background()
	id := uuid.NewV4()

	deleteProduct := &commands.DeleteProduct{
		ProductID: id,
	}

	c.ProductRepository.On("DeleteProductByID", mock.Anything, id).
		Once().
		Return(customErrors.NewNotFoundError("error finding product"))

	deleteProductHandler := commands.NewDeleteProductHandler(c.Log, c.Uow, c.Bus, c.Tracer)

	res, err := deleteProductHandler.Handle(ctx, deleteProduct)
	c.Require().Error(err)
	c.True(customErrors.IsNotFoundError(err))
	c.True(customErrors.IsApplicationError(err, http.StatusNotFound))
	c.NotNil(res)

	c.Uow.AssertNumberOfCalls(c.T(), "Do", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "DeleteProductByID", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 0)
}

func (c *deleteProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_Error_In_Bus() {
	ctx := context.Background()
	id := c.Items[0].ProductId

	deleteProduct := &commands.DeleteProduct{
		ProductID: id,
	}

	c.ProductRepository.On("DeleteProductByID", mock.Anything, id).
		Once().
		Return(nil)

	// override called mock
	// https://github.com/stretchr/testify/issues/558
	c.Bus.Mock.ExpectedCalls = nil
	c.Bus.On("PublishMessage", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return(errors.New("error in the publish message"))

	deleteProductHandler := commands.NewDeleteProductHandler(c.Log, c.Uow, c.Bus, c.Tracer)
	dto, err := deleteProductHandler.Handle(ctx, deleteProduct)

	c.NotNil(dto)

	c.Uow.AssertNumberOfCalls(c.T(), "Do", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "DeleteProductByID", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
	c.True(customErrors.IsApplicationError(err, http.StatusInternalServerError))
	c.ErrorContains(err, "error in publishing 'ProductDeleted' message")
}
