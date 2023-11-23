//go:build unit
// +build unit

package commands

import (
	"testing"

	v1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/deleting_product/v1"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/unit_test"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type deleteProductUnitTests struct {
	*unit_test.UnitTestSharedFixture
}

func TestDeleteProductByIdUnit(t *testing.T) {
	suite.Run(
		t,
		&deleteProductUnitTests{UnitTestSharedFixture: unit_test.NewUnitTestSharedFixture(t)},
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
