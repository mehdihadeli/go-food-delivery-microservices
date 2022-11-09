package v1

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"

	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/mocks/testData"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/unit_test"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type getProductByIdHandlerTest struct {
	*unit_test.UnitTestSharedFixture
	*unit_test.UnitTestMockFixture
	getProductByIdHandler *GetProductByIdHandler
}

func TestGetProductByIdUnit(t *testing.T) {
	suite.Run(t, &getProductByIdHandlerTest{UnitTestSharedFixture: unit_test.NewUnitTestSharedFixture(t)})
}

func (c *getProductByIdHandlerTest) SetupTest() {
}

func (c *getProductByIdHandlerTest) BeforeTest(suiteName, testName string) {
}

func (c *getProductByIdHandlerTest) Test_Get_Product_By_Id() {
	product := testData.Products[0]
	id := uuid.NewV4()

	testCases := []struct {
		Name                          string
		id                            uuid.UUID
		HandlerError                  error
		ProductRepositoryNumberOfCall int
		ExpectedName                  string
		ExpectedId                    uuid.UUID
		RepositoryReturnProduct       *models.Product
		RepositoryReturnError         error
		fn                            func()
	}{
		{
			Name:                          "Handle_Should_Get_Product_Successfully",
			id:                            product.ProductId,
			HandlerError:                  nil,
			ProductRepositoryNumberOfCall: 1,
			ExpectedId:                    product.ProductId,
			ExpectedName:                  product.Name,
			RepositoryReturnProduct:       product,
			RepositoryReturnError:         nil,
		},
		{
			Name:                          "Handle_Should_Return_NotFound_Error_For_NotFound_Item",
			id:                            id,
			HandlerError:                  customErrors.NewApplicationErrorWithCode(fmt.Sprintf("error in getting product with id %s in the repository", id.String()), http.StatusNotFound),
			ProductRepositoryNumberOfCall: 1,
			ExpectedId:                    *new(uuid.UUID),
			ExpectedName:                  "",
			RepositoryReturnProduct:       nil,
			RepositoryReturnError:         customErrors.NewNotFoundError("product not found"),
		},
		{
			Name:                          "Handle_Should_Return_Error_For_Error_In_Mapping",
			id:                            product.ProductId,
			HandlerError:                  customErrors.NewApplicationErrorWithCode("error in the mapping product", http.StatusInternalServerError),
			ProductRepositoryNumberOfCall: 1,
			ExpectedId:                    *new(uuid.UUID),
			ExpectedName:                  "",
			RepositoryReturnProduct:       product,
			RepositoryReturnError:         nil,
			fn: func() {
				mapper.ClearMappings()
			},
		},
	}

	for _, testCase := range testCases {
		c.Run(testCase.Name, func() {
			// arrange
			// create new mocks or clear mocks before executing
			c.UnitTestMockFixture = unit_test.NewUnitTestMockFixture(c.T())

			c.getProductByIdHandler = NewGetProductByIdHandler(c.Log, c.Cfg, c.ProductRepository)

			c.ProductRepository.On("GetProductById", mock.Anything, testCase.id).
				Once().
				Return(testCase.RepositoryReturnProduct, testCase.RepositoryReturnError)

			if testCase.fn != nil {
				testCase.fn()
			}

			query := NewGetProductById(testCase.id)

			// act
			dto, err := c.getProductByIdHandler.Handle(c.Ctx, query)

			// assert
			c.ProductRepository.AssertNumberOfCalls(c.T(), "GetProductById", testCase.ProductRepositoryNumberOfCall)
			if testCase.HandlerError == nil {
				// success path with a valid dto
				c.Require().NoError(err)
				c.NotNil(dto.Product)
				c.Equal(testCase.ExpectedId, dto.Product.ProductId)
				c.Equal(testCase.ExpectedName, dto.Product.Name)
			} else {
				// handler error path
				c.Nil(dto)
				c.ErrorContains(err, testCase.HandlerError.Error())
				if customErrors.IsApplicationError(testCase.HandlerError, http.StatusNotFound) {
					// not found error
					c.True(customErrors.IsNotFoundError(err))
					c.True(customErrors.IsApplicationError(err, http.StatusNotFound))
					c.ErrorContains(err, testCase.HandlerError.Error())
				} else {
					// mapping error
					c.ErrorContains(err, testCase.HandlerError.Error())
					c.True(customErrors.IsApplicationError(err, http.StatusInternalServerError))
				}
			}
		})
	}

	//c.Run("Handle_Should_Get_Product_Successfully", func() {
	//	// create new mocks or clear mocks before executing
	//	c.UnitTestMockFixture = unit_test.NewUnitTestMockFixture(c.T())
	//	c.getProductByIdHandler = NewGetProductByIdHandler(c.Log, c.Cfg, c.ProductRepository)
	//
	//	c.ProductRepository.On("GetProductById", mock.Anything, product.ProductId).
	//		Once().
	//		Return(product, nil)
	//
	//	query := NewGetProductById(product.ProductId)
	//
	//	dto, err := c.getProductByIdHandler.Handle(c.Ctx, query)
	//	c.Require().NoError(err)
	//
	//	c.ProductRepository.AssertNumberOfCalls(c.T(), "GetProductById", 1)
	//	c.Equal(product.ProductId, dto.Product.ProductId)
	//	c.Equal(product.Name, dto.Product.Name)
	//})
	//
	//c.Run("Handle_Should_Return_NotFound_Error_For_NotFound_Item", func() {
	//	// create new mocks or clear mocks before executing
	//	c.UnitTestMockFixture = unit_test.NewUnitTestMockFixture(c.T())
	//	c.getProductByIdHandler = NewGetProductByIdHandler(c.Log, c.Cfg, c.ProductRepository)
	//
	//	c.ProductRepository.On("GetProductById", mock.Anything, id).
	//		Once().
	//		Return(nil, customErrors.NewNotFoundError("product not found"))
	//
	//	query := NewGetProductById(id)
	//
	//	dto, err := c.getProductByIdHandler.Handle(c.Ctx, query)
	//
	//	c.Require().Error(err)
	//	c.True(customErrors.IsApplicationError(err, http.StatusNotFound))
	//	c.True(customErrors.IsNotFoundError(err))
	//	c.ErrorContains(err, fmt.Sprintf("error in getting product with id %s in the repository", id.String()))
	//	c.Nil(dto)
	//
	//	c.ProductRepository.AssertNumberOfCalls(c.T(), "GetProductById", 1)
	//})
	//
	//c.Run("Handle_Should_Return_Error_For_Error_In_Mapping", func() {
	//	// create new mocks or clear mocks before executing
	//	c.UnitTestMockFixture = unit_test.NewUnitTestMockFixture(c.T())
	//	c.getProductByIdHandler = NewGetProductByIdHandler(c.Log, c.Cfg, c.ProductRepository)
	//
	//	product := testData.Products[0]
	//	c.ProductRepository.On("GetProductById", mock.Anything, product.ProductId).
	//		Once().
	//		Return(product, nil)
	//
	//	mapper.ClearMappings()
	//
	//	query := NewGetProductById(product.ProductId)
	//
	//	dto, err := c.getProductByIdHandler.Handle(c.Ctx, query)
	//
	//	c.ProductRepository.AssertNumberOfCalls(c.T(), "GetProductById", 1)
	//	c.Nil(dto)
	//	c.Require().Error(err)
	//	c.True(customErrors.IsApplicationError(err, http.StatusInternalServerError))
	//	c.ErrorContains(err, "error in the mapping product")
	//})
}
