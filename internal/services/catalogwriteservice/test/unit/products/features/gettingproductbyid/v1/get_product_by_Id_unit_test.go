//go:build unit
// +build unit

package v1

import (
	"testing"

	getProductByIdQuery "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/gettingproductbyid/v1"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type getProductByIdUnitTests struct {
	*unittest.UnitTestSharedFixture
}

func TestGetProductByIdUnit(t *testing.T) {
	suite.Run(
		t,
		&getProductByIdUnitTests{UnitTestSharedFixture: unittest.NewUnitTestSharedFixture(t)},
	)
}

func (c *getProductByIdUnitTests) Test_New_Get_Product_By_Id_Should_Return_No_Error_For_Valid_Input() {
	id := uuid.NewV4()

	query, err := getProductByIdQuery.NewGetProductById(id)

	c.Assert().NotNil(query)
	c.Assert().Equal(query.ProductID, id)
	c.Require().NoError(err)
}

func (c *getProductByIdUnitTests) Test_New_Get_Product_By_Id_Should_Return_Error_For_Invalid_Id() {
	query, err := getProductByIdQuery.NewGetProductById(uuid.UUID{})

	c.Assert().Nil(query)
	c.Require().Error(err)
}
