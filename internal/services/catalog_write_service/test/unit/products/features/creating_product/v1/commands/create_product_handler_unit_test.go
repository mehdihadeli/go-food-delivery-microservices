//go:build unit
// +build unit

package commands

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	createProductCommandV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creating_product/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/mocks/testData"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/unit_test"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mapper"
)

type createProductHandlerUnitTests struct {
	*unit_test.UnitTestSharedFixture
}

func TestCreateProductHandlerUnit(t *testing.T) {
	suite.Run(
		t,
		&createProductHandlerUnitTests{
			UnitTestSharedFixture: unit_test.NewUnitTestSharedFixture(t),
		},
	)
}

func (c *createProductHandlerUnitTests) Test_Handle_Should_Create_New_Product_With_Valid_Data() {
	ctx := context.Background()
	id := uuid.NewV4()

	createProductHandler := createProductCommandV1.NewCreateProductHandler(
		c.Log,
		c.Uow,
		c.Bus,
		c.Tracer,
	)

	createProduct := &createProductCommandV1.CreateProduct{
		ProductID:   id,
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	product := testData.Products[0]

	c.ProductRepository.On("CreateProduct", mock.Anything, mock.Anything).
		Once().
		Return(product, nil)

	dto, err := createProductHandler.Handle(ctx, createProduct)
	c.Require().NoError(err)

	c.Uow.AssertNumberOfCalls(c.T(), "Do", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "CreateProduct", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
	c.Equal(dto.ProductID, id)
}

func (c *createProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_Duplicate_Item() {
	ctx := context.Background()
	id := uuid.NewV4()

	createProduct := &createProductCommandV1.CreateProduct{
		ProductID:   id,
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	c.ProductRepository.On("CreateProduct", mock.Anything, mock.Anything).
		Once().
		Return(nil, errors.New("error duplicate product"))

	createProductHandler := createProductCommandV1.NewCreateProductHandler(
		c.Log,
		c.Uow,
		c.Bus,
		c.Tracer,
	)

	dto, err := createProductHandler.Handle(ctx, createProduct)

	c.Uow.AssertNumberOfCalls(c.T(), "Do", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "CreateProduct", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 0)
	c.True(customErrors.IsApplicationError(err, http.StatusConflict))
	c.ErrorContains(err, "product already exists")
	c.Nil(dto)
}

func (c *createProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_Error_In_Bus() {
	ctx := context.Background()
	id := uuid.NewV4()

	createProduct := &createProductCommandV1.CreateProduct{
		ProductID:   id,
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	product := testData.Products[0]
	c.ProductRepository.On("CreateProduct", mock.Anything, mock.Anything).
		Once().
		Return(product, nil)

	// override called mock
	// https://github.com/stretchr/testify/issues/558
	c.Bus.Mock.ExpectedCalls = nil
	c.Bus.On("PublishMessage", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return(errors.New("error in the publish message"))

	createProductHandler := createProductCommandV1.NewCreateProductHandler(
		c.Log,
		c.Uow,
		c.Bus,
		c.Tracer,
	)

	dto, err := createProductHandler.Handle(ctx, createProduct)

	c.Uow.AssertNumberOfCalls(c.T(), "Do", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "CreateProduct", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
	c.ErrorContains(err, "error in the publish message")
	c.ErrorContains(err, "error in publishing ProductCreated integration_events event")
	c.Nil(dto)
}

func (c *createProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_Error_In_Mapping() {
	ctx := context.Background()
	id := uuid.NewV4()

	createProduct := &createProductCommandV1.CreateProduct{
		ProductID:   id,
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	product := testData.Products[0]
	c.ProductRepository.On("CreateProduct", mock.Anything, mock.Anything).
		Once().
		Return(product, nil)

	mapper.ClearMappings()

	createProductHandler := createProductCommandV1.NewCreateProductHandler(
		c.Log,
		c.Uow,
		c.Bus,
		c.Tracer,
	)

	dto, err := createProductHandler.Handle(ctx, createProduct)

	c.Uow.AssertNumberOfCalls(c.T(), "Do", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "CreateProduct", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 0)
	c.ErrorContains(err, "error in the mapping ProductDto")
	c.True(customErrors.IsApplicationError(err, http.StatusInternalServerError))
	c.Nil(dto)
}
