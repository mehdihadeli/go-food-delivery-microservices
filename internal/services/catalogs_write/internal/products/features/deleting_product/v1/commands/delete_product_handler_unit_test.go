//go:build.sh unit
// +build.sh unit

package commands

import (
	"net/http"
	"testing"

	"emperror.dev/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/mocks/testData"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/test_fixtures/unit_test"
)

type deleteProductHandlerUnitTests struct {
	*unit_test.UnitTestSharedFixture
	*unit_test.UnitTestMockFixture
	deleteProductHandler *DeleteProductHandler
}

func TestDeleteProductHandlerUnit(t *testing.T) {
	suite.Run(
		t,
		&deleteProductHandlerUnitTests{
			UnitTestSharedFixture: unit_test.NewUnitTestSharedFixture(t),
		},
	)
}

func (c *deleteProductHandlerUnitTests) SetupTest() {
	// create new mocks or clear mocks before executing
	c.UnitTestMockFixture = unit_test.NewUnitTestMockFixture(c.T())
	c.deleteProductHandler = NewDeleteProductHandler(c.Log, c.Cfg, c.Uow, c.Bus)
}

func (c *deleteProductHandlerUnitTests) Test_Handle_Should_Delete_Product_With_Valid_Id() {
	id := testData.Products[0].ProductId

	deleteProduct := &DeleteProduct{
		ProductID: id,
	}

	c.ProductRepository.On("DeleteProductByID", mock.Anything, id).
		Once().
		Return(nil)

	_, err := c.deleteProductHandler.Handle(c.Ctx, deleteProduct)
	c.Require().NoError(err)

	c.Uow.AssertNumberOfCalls(c.T(), "Do", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "DeleteProductByID", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
}

func (c *deleteProductHandlerUnitTests) Test_Handle_Should_Return_NotFound_Error_When_Id_Is_Invalid() {
	id := uuid.NewV4()

	deleteProduct := &DeleteProduct{
		ProductID: id,
	}

	c.ProductRepository.On("DeleteProductByID", mock.Anything, id).
		Once().
		Return(customErrors.NewNotFoundError("error finding product"))

	res, err := c.deleteProductHandler.Handle(c.Ctx, deleteProduct)
	c.Require().Error(err)
	c.True(customErrors.IsNotFoundError(err))
	c.True(customErrors.IsApplicationError(err, http.StatusNotFound))
	c.NotNil(res)

	c.Uow.AssertNumberOfCalls(c.T(), "Do", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "DeleteProductByID", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 0)
}

func (c *deleteProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_Error_In_Bus() {
	id := testData.Products[0].ProductId

	deleteProduct := &DeleteProduct{
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

	dto, err := c.deleteProductHandler.Handle(c.Ctx, deleteProduct)

	c.NotNil(dto)

	c.Uow.AssertNumberOfCalls(c.T(), "Do", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "DeleteProductByID", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
	c.True(customErrors.IsApplicationError(err, http.StatusInternalServerError))
	c.ErrorContains(err, "error in publishing 'ProductDeleted' message")
}
