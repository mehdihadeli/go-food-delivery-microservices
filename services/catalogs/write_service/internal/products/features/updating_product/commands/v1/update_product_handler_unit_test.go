package v1

import (
	"testing"

	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/mocks/testData"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/unit_test"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type createProductHandlerUnitTests struct {
	*unit_test.UnitTestSharedFixture
	*unit_test.UnitTestMockFixture
	updateProductHandler *UpdateProductHandler
}

func TestCreateProductUnit(t *testing.T) {
	suite.Run(t, &createProductHandlerUnitTests{UnitTestSharedFixture: unit_test.NewUnitTestSharedFixture(t)})
}

func (c *createProductHandlerUnitTests) SetupTest() {
	// create new mocks or clear mocks before executing
	c.UnitTestMockFixture = unit_test.NewUnitTestMockFixture(c.T())
	c.updateProductHandler = NewUpdateProductHandler(c.Log, c.Cfg, c.ProductRepository, c.Bus)
}

func (c *createProductHandlerUnitTests) Test_Handle_Should_Update_Product_With_Valid_Data() {
	existing := testData.Products[0]

	updateProductCommand := NewUpdateProduct(existing.ProductId, gofakeit.Name(), gofakeit.EmojiDescription(), existing.Price)

	updated := &models.Product{
		ProductId:   existing.ProductId,
		Name:        updateProductCommand.Name,
		Description: updateProductCommand.Description,
		UpdatedAt:   updateProductCommand.UpdatedAt,
		CreatedAt:   existing.CreatedAt,
		Price:       existing.Price,
	}

	c.ProductRepository.On("UpdateProduct", mock.Anything, mock.Anything).
		Once().
		Return(updated, nil)

	_, err := c.updateProductHandler.Handle(c.Ctx, updateProductCommand)
	c.Require().NoError(err)

	c.ProductRepository.AssertNumberOfCalls(c.T(), "UpdateProduct", 1)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
}

//func (c *createProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_Duplicate_Item() {
//	id := uuid.NewV4()
//
//	createProduct := &CreateProduct{
//		ProductID:   id,
//		Name:        gofakeit.Name(),
//		CreatedAt:   time.Now(),
//		Description: gofakeit.EmojiDescription(),
//		Price:       gofakeit.Price(100, 1000),
//	}
//
//	c.ProductRepository.On("CreateProduct", mock.Anything, mock.Anything).
//		Once().
//		Return(nil, errors.New("error duplicate product"))
//
//	dto, err := c.createProductHandler.Handle(c.Ctx, createProduct)
//
//	c.ProductRepository.AssertNumberOfCalls(c.T(), "CreateProduct", 1)
//	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 0)
//	c.True(customErrors.IsApplicationError(err, http.StatusConflict))
//	c.ErrorContains(err, "product already exists")
//	c.Nil(dto)
//}
//
//func (c *createProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_Error_In_Bus() {
//	id := uuid.NewV4()
//
//	createProduct := &CreateProduct{
//		ProductID:   id,
//		Name:        gofakeit.Name(),
//		CreatedAt:   time.Now(),
//		Description: gofakeit.EmojiDescription(),
//		Price:       gofakeit.Price(100, 1000),
//	}
//
//	product := testData.Products[0]
//	c.ProductRepository.On("CreateProduct", mock.Anything, mock.Anything).
//		Once().
//		Return(product, nil)
//
//	// override called mock
//	// https://github.com/stretchr/testify/issues/558
//	c.Bus.Mock.ExpectedCalls = nil
//	c.Bus.On("PublishMessage", mock.Anything, mock.Anything, mock.Anything).Once().Return(errors.New("error in the publish message"))
//
//	dto, err := c.createProductHandler.Handle(c.Ctx, createProduct)
//
//	c.ProductRepository.AssertNumberOfCalls(c.T(), "CreateProduct", 1)
//	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
//	c.ErrorContains(err, "error in the publish message")
//	c.ErrorContains(err, "error in publishing ProductCreated integration event")
//	c.Nil(dto)
//}
//
//func (c *createProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_Error_In_Mapping() {
//	id := uuid.NewV4()
//
//	createProduct := &CreateProduct{
//		ProductID:   id,
//		Name:        gofakeit.Name(),
//		CreatedAt:   time.Now(),
//		Description: gofakeit.EmojiDescription(),
//		Price:       gofakeit.Price(100, 1000),
//	}
//
//	product := testData.Products[0]
//	c.ProductRepository.On("CreateProduct", mock.Anything, mock.Anything).
//		Once().
//		Return(product, nil)
//
//	mapper.ClearMappings()
//
//	dto, err := c.createProductHandler.Handle(c.Ctx, createProduct)
//
//	c.ProductRepository.AssertNumberOfCalls(c.T(), "CreateProduct", 1)
//	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 0)
//	c.ErrorContains(err, "error in the mapping ProductDto")
//	c.True(customErrors.IsApplicationError(err, http.StatusInternalServerError))
//	c.Nil(dto)
//}
