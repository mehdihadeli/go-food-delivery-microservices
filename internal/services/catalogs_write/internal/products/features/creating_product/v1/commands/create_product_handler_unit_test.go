package createProductCommand

import (
	"net/http"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mapper"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/mocks/testData"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/test_fixtures/unit_test"
)

type createProductHandlerUnitTests struct {
	*unit_test.UnitTestSharedFixture
	*unit_test.UnitTestMockFixture
	createProductHandler *CreateProductHandler
}

func TestCreateProductHandlerUnit(t *testing.T) {
	suite.Run(t, &createProductHandlerUnitTests{UnitTestSharedFixture: unit_test.NewUnitTestSharedFixture(t)})
}

func (c *createProductHandlerUnitTests) SetupTest() {
	// create new mocks or clear mocks before executing
	c.UnitTestMockFixture = unit_test.NewUnitTestMockFixture(c.T())
	c.createProductHandler = NewCreateProductHandler(c.Log, c.Cfg, c.Uow, c.Bus)
}

func (c *createProductHandlerUnitTests) Test_Handle_Should_Create_New_Product_With_Valid_Data() {
	id := uuid.NewV4()

	createProduct := &CreateProduct{
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

	dto, err := c.createProductHandler.Handle(c.Ctx, createProduct)
	c.Require().NoError(err)

	c.Uow.AssertNumberOfCalls(c.T(), "Do", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "CreateProduct", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
	c.Equal(dto.ProductID, id)
}

func (c *createProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_Duplicate_Item() {
	id := uuid.NewV4()

	createProduct := &CreateProduct{
		ProductID:   id,
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	c.ProductRepository.On("CreateProduct", mock.Anything, mock.Anything).
		Once().
		Return(nil, errors.New("error duplicate product"))

	dto, err := c.createProductHandler.Handle(c.Ctx, createProduct)

	c.Uow.AssertNumberOfCalls(c.T(), "Do", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "CreateProduct", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 0)
	c.True(customErrors.IsApplicationError(err, http.StatusConflict))
	c.ErrorContains(err, "product already exists")
	c.Nil(dto)
}

func (c *createProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_Error_In_Bus() {
	id := uuid.NewV4()

	createProduct := &CreateProduct{
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
	c.Bus.On("PublishMessage", mock.Anything, mock.Anything, mock.Anything).Once().Return(errors.New("error in the publish message"))

	dto, err := c.createProductHandler.Handle(c.Ctx, createProduct)

	c.Uow.AssertNumberOfCalls(c.T(), "Do", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "CreateProduct", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
	c.ErrorContains(err, "error in the publish message")
	c.ErrorContains(err, "error in publishing ProductCreated integration_events event")
	c.Nil(dto)
}

func (c *createProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_Error_In_Mapping() {
	id := uuid.NewV4()

	createProduct := &CreateProduct{
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

	dto, err := c.createProductHandler.Handle(c.Ctx, createProduct)

	c.Uow.AssertNumberOfCalls(c.T(), "Do", 1)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "CreateProduct", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 0)
	c.ErrorContains(err, "error in the mapping ProductDto")
	c.True(customErrors.IsApplicationError(err, http.StatusInternalServerError))
	c.Nil(dto)
}
