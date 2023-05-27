package getProductByIdQuery

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/test_fixtures/unit_test"
)

type getProductByIdUnitTests struct {
	*unit_test.UnitTestSharedFixture
	*unit_test.UnitTestMockFixture
}

func TestGetProductByIdUnit(t *testing.T) {
	suite.Run(t, &getProductByIdUnitTests{UnitTestSharedFixture: unit_test.NewUnitTestSharedFixture(t)})
}

func (c *getProductByIdUnitTests) SetupTest() {
	// create new mocks or clear mocks before executing
	c.UnitTestMockFixture = unit_test.NewUnitTestMockFixture(c.T())
}

func (c *getProductByIdUnitTests) Test_New_Get_Product_By_Id_Should_Return_No_Error_For_Valid_Input() {
	id := uuid.NewV4()

	query, err := NewGetProductById(id)

	c.Assert().NotNil(query)
	c.Assert().Equal(query.ProductID, id)
	c.Require().NoError(err)
}

func (c *getProductByIdUnitTests) Test_New_Get_Product_By_Id_Should_Return_Error_For_Invalid_Id() {
	query, err := NewGetProductById(uuid.UUID{})

	c.Assert().Nil(query)
	c.Require().Error(err)
}
