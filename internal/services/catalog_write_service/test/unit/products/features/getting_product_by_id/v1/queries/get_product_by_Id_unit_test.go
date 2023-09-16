//go:build unit
// +build unit

package queries

import (
	"testing"

	getProductByIdQuery "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_product_by_id/v1/queries"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/unit_test"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type getProductByIdUnitTests struct {
	*unit_test.UnitTestSharedFixture
}

func TestGetProductByIdUnit(t *testing.T) {
	suite.Run(
		t,
		&getProductByIdUnitTests{UnitTestSharedFixture: unit_test.NewUnitTestSharedFixture(t)},
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
