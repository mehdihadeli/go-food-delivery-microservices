//go:build unit
// +build unit

package commands

import (
	"context"
	"net/http"
	"testing"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/updating_product/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/unit_test"

	"emperror.dev/errors"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type updateProductHandlerUnitTests struct {
	*unit_test.UnitTestSharedFixture
}

func TestUpdateProductHandlerUnit(t *testing.T) {
	suite.Run(
		t,
		&updateProductHandlerUnitTests{
			UnitTestSharedFixture: unit_test.NewUnitTestSharedFixture(t),
		},
	)
}

func (c *updateProductHandlerUnitTests) Test_Handle_Should_Update_Product_With_Valid_Data() {
	ctx := context.Background()
	existing := c.Items[0]

	updateProductCommand, err := commands.NewUpdateProduct(
		existing.ProductId,
		gofakeit.Name(),
		gofakeit.EmojiDescription(),
		existing.Price,
	)
	c.Require().NoError(err)

	updated := &models.Product{
		ProductId:   existing.ProductId,
		Name:        updateProductCommand.Name,
		Description: updateProductCommand.Description,
		UpdatedAt:   updateProductCommand.UpdatedAt,
		CreatedAt:   existing.CreatedAt,
		Price:       existing.Price,
	}

	c.ProductRepository.On("GetProductById", mock.Anything, existing.ProductId).
		Once().
		Return(existing, nil)

	c.ProductRepository.On("UpdateProduct", mock.Anything, mock.Anything).
		Once().
		Return(updated, nil)

	updateProductHandler := commands.NewUpdateProductHandler(c.Log, c.Uow, c.Bus, c.Tracer)

	_, err = updateProductHandler.Handle(ctx, updateProductCommand)
	c.Require().NoError(err)

	c.Uow.AssertNumberOfCalls(c.T(), "Do", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "GetProductById", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "UpdateProduct", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
}

func (c *updateProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_NotFound_Item() {
	ctx := context.Background()
	id := uuid.NewV4()

	command, err := commands.NewUpdateProduct(
		id,
		gofakeit.Name(),
		gofakeit.EmojiDescription(),
		gofakeit.Price(150, 6000),
	)
	c.Require().NoError(err)

	c.ProductRepository.On("GetProductById", mock.Anything, mock.Anything).
		Once().
		Return(nil, errors.New("error notfound product"))

	updateProductHandler := commands.NewUpdateProductHandler(c.Log, c.Uow, c.Bus, c.Tracer)
	dto, err := updateProductHandler.Handle(ctx, command)

	c.Uow.AssertNumberOfCalls(c.T(), "Do", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "GetProductById", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "UpdateProduct", 0)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 0)
	c.True(customErrors.IsApplicationError(err, http.StatusNotFound))
	c.ErrorContains(err, "error notfound product")
	c.NotNil(dto)
}

func (c *updateProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_Error_In_Bus() {
	ctx := context.Background()
	existing := c.Items[0]

	updateProductCommand, err := commands.NewUpdateProduct(
		existing.ProductId,
		gofakeit.Name(),
		gofakeit.EmojiDescription(),
		existing.Price,
	)
	c.Require().NoError(err)

	updated := &models.Product{
		ProductId:   existing.ProductId,
		Name:        updateProductCommand.Name,
		Description: updateProductCommand.Description,
		UpdatedAt:   updateProductCommand.UpdatedAt,
		CreatedAt:   existing.CreatedAt,
		Price:       existing.Price,
	}

	c.ProductRepository.On("GetProductById", mock.Anything, existing.ProductId).
		Once().
		Return(existing, nil)

	c.ProductRepository.On("UpdateProduct", mock.Anything, mock.Anything).
		Once().
		Return(updated, nil)

	// override called mock
	// https://github.com/stretchr/testify/issues/558
	c.Bus.Mock.ExpectedCalls = nil
	c.Bus.On("PublishMessage", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return(errors.New("error in the publish message"))

	updateProductHandler := commands.NewUpdateProductHandler(c.Log, c.Uow, c.Bus, c.Tracer)
	dto, err := updateProductHandler.Handle(ctx, updateProductCommand)

	c.Uow.AssertNumberOfCalls(c.T(), "Do", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "UpdateProduct", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "GetProductById", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
	c.ErrorContains(err, "error in the publish message")
	c.ErrorContains(err, "error in publishing 'ProductUpdated' message")
	c.NotNil(dto)
}
