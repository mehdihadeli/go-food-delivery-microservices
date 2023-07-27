//go:build.sh unit
// +build.sh unit

package commands

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/test_fixtures/unit_test"
)

type deleteProductUnitTests struct {
	*unit_test.UnitTestSharedFixture
	*unit_test.UnitTestMockFixture
}

func TestDeleteProductByIdUnit(t *testing.T) {
	suite.Run(
		t,
		&deleteProductUnitTests{UnitTestSharedFixture: unit_test.NewUnitTestSharedFixture(t)},
	)
}

func (c *deleteProductUnitTests) SetupTest() {
	// create new mocks or clear mocks before executing
	c.UnitTestMockFixture = unit_test.NewUnitTestMockFixture(c.T())
}

func (c *deleteProductUnitTests) Test_New_Delete_Product_Should_Return_No_Error_For_Valid_Input() {
	id := uuid.NewV4()

	query, err := NewDeleteProduct(id)

	c.Assert().NotNil(query)
	c.Assert().Equal(query.ProductID, id)
	c.Require().NoError(err)
}

func (c *deleteProductUnitTests) Test_New_Delete_Product_Should_Return_Error_For_Invalid_Id() {
	query, err := NewDeleteProduct(uuid.UUID{})

	c.Assert().Nil(query)
	c.Require().Error(err)
}
