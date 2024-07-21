//go:build unit
// +build unit

package v1

import (
	"testing"

	v1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/deletingproduct/v1"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type deleteProductUnitTests struct {
	*unittest.UnitTestSharedFixture
}

func TestDeleteProductByIdUnit(t *testing.T) {
	suite.Run(
		t,
		&deleteProductUnitTests{UnitTestSharedFixture: unittest.NewUnitTestSharedFixture(t)},
	)
}

func (c *deleteProductUnitTests) Test_New_Delete_Product_Should_Return_No_Error_For_Valid_Input() {
	id := uuid.NewV4()

	query, err := v1.NewDeleteProduct(id)

	c.Assert().NotNil(query)
	c.Assert().Equal(query.ProductID, id)
	c.Require().NoError(err)
}

func (c *deleteProductUnitTests) Test_New_Delete_Product_Should_Return_Error_For_Invalid_Id() {
	query, err := v1.NewDeleteProduct(uuid.UUID{})

	c.Assert().Nil(query)
	c.Require().Error(err)
}
