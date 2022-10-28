package v1

import (
	"emperror.dev/errors"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/mocks/test_data"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/unit_test"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

// https://github.com/MarkNijhof/Fohjin/
// https://jeremydmiller.com/2022/10/24/using-context-specification-to-better-express-complicated-tests/
type createProductHandlerTest struct {
	*unit_test.UnitTestSharedFixture
	*unit_test.UnitTestMockFixture
	createProductHandler *CreateProductHandler
}

func TestCreateProductUnit(t *testing.T) {
	suite.Run(t, &createProductHandlerTest{UnitTestSharedFixture: unit_test.NewUnitTestSharedFixture(t)})
}

func (c *createProductHandlerTest) SetupTest() {
	// create new mocks or clear mocks before executing
	c.UnitTestMockFixture = unit_test.NewUnitTestMockFixture(c.T())
	c.createProductHandler = NewCreateProductHandler(c.Log, c.Cfg, c.ProductRepository, c.Bus)
}

func (c *createProductHandlerTest) Test_Handle_Should_Create_New_Product_With_Valid_Data() {
	id := uuid.NewV4()

	createProduct := &CreateProduct{
		ProductID:   id,
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	product := test_data.Products[0]

	c.ProductRepository.On("CreateProduct", mock.Anything, mock.Anything).
		Once().
		Return(product, nil)

	dto, err := c.createProductHandler.Handle(c.Ctx, createProduct)
	c.Require().NoError(err)

	c.ProductRepository.AssertNumberOfCalls(c.T(), "CreateProduct", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
	c.Equal(dto.ProductID, id)
}

func (c *createProductHandlerTest) Test_Handle_Should_Return_Error_For_Error_In_Repository() {
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
		Return(nil, errors.New("error creating product"))

	dto, err := c.createProductHandler.Handle(c.Ctx, createProduct)

	c.ProductRepository.AssertNumberOfCalls(c.T(), "CreateProduct", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 0)
	c.ErrorContains(err, "error creating product")
	c.ErrorContains(err, "error in creating product in the repository")
	c.Nil(dto)
}

func (c *createProductHandlerTest) Test_Handle_Should_Return_Error_For_Error_In_Bus() {
	id := uuid.NewV4()

	createProduct := &CreateProduct{
		ProductID:   id,
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	product := test_data.Products[0]
	c.ProductRepository.On("CreateProduct", mock.Anything, mock.Anything).
		Once().
		Return(product, nil)

	// override called mock
	//https://github.com/stretchr/testify/issues/558
	c.Bus.Mock.ExpectedCalls = nil
	c.Bus.On("PublishMessage", mock.Anything, mock.Anything, mock.Anything).Once().Return(errors.New("error in the publish message"))

	dto, err := c.createProductHandler.Handle(c.Ctx, createProduct)

	c.ProductRepository.AssertNumberOfCalls(c.T(), "CreateProduct", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
	c.ErrorContains(err, "error in the publish message")
	c.ErrorContains(err, "error in publishing ProductCreated integration event")
	c.Nil(dto)
}

func (c *createProductHandlerTest) Test_Handle_Should_Return_Error_For_Error_In_Mapping() {
	id := uuid.NewV4()

	createProduct := &CreateProduct{
		ProductID:   id,
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	product := test_data.Products[0]
	c.ProductRepository.On("CreateProduct", mock.Anything, mock.Anything).
		Once().
		Return(product, nil)

	mapper.ClearMappings()

	dto, err := c.createProductHandler.Handle(c.Ctx, createProduct)

	c.ProductRepository.AssertNumberOfCalls(c.T(), "CreateProduct", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 0)
	c.ErrorContains(err, "error in the mapping ProductDto")
	c.Nil(dto)
}
