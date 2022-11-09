package v1

import (
	"net/http"
	"testing"

	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"

	"emperror.dev/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/mocks/testData"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/unit_test"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type deleteProductHandlerUnitTests struct {
	*unit_test.UnitTestSharedFixture
	*unit_test.UnitTestMockFixture
	deleteProductHandler *DeleteProductHandler
}

func TestDeleteProductUnit(t *testing.T) {
	suite.Run(t, &deleteProductHandlerUnitTests{UnitTestSharedFixture: unit_test.NewUnitTestSharedFixture(t)})
}

func (c *deleteProductHandlerUnitTests) SetupTest() {
	// create new mocks or clear mocks before executing
	c.UnitTestMockFixture = unit_test.NewUnitTestMockFixture(c.T())
	c.deleteProductHandler = NewDeleteProductHandler(c.Log, c.Cfg, c.ProductRepository, c.Bus)
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
	c.Assert().True(customErrors.IsNotFoundError(err))
	c.Assert().True(customErrors.IsApplicationError(err, http.StatusNotFound))
	c.Assert().Nil(res)

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
	c.Bus.On("PublishMessage", mock.Anything, mock.Anything, mock.Anything).Once().Return(errors.New("error in the publish message"))

	dto, err := c.deleteProductHandler.Handle(c.Ctx, deleteProduct)

	c.Nil(dto)

	c.ProductRepository.AssertNumberOfCalls(c.T(), "DeleteProductByID", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
	c.True(customErrors.IsApplicationError(err, http.StatusInternalServerError))
	c.ErrorContains(err, "error in publishing 'ProductDeleted' message")
}
