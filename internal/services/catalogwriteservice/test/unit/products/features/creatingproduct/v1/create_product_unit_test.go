//go:build unit
// +build unit

package v1

import (
	"testing"

	createProductCommand "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/suite"
)

type createProductUnitTests struct {
	*unittest.UnitTestSharedFixture
}

func TestCreateProductUnit(t *testing.T) {
	suite.Run(
		t,
		&createProductUnitTests{
			UnitTestSharedFixture: unittest.NewUnitTestSharedFixture(t),
		},
	)
}

func (c *createProductUnitTests) Test_New_Create_Product_Should_Return_No_Error_For_Valid_Input() {
	name := gofakeit.Name()
	description := gofakeit.EmojiDescription()
	price := gofakeit.Price(150, 6000)

	updateProduct, err := createProductCommand.NewCreateProduct(
		name,
		description,
		price,
	)

	c.Assert().NotNil(updateProduct)
	c.Assert().Equal(name, updateProduct.Name)
	c.Assert().Equal(price, updateProduct.Price)

	c.Require().NoError(err)
}

func (c *createProductUnitTests) Test_New_Create_Product_Should_Return_Error_For_Invalid_Price() {
	command, err := createProductCommand.NewCreateProduct(
		gofakeit.Name(),
		gofakeit.EmojiDescription(),
		0,
	)

	c.Require().Error(err)
	c.Assert().Nil(command)
}

func (c *createProductUnitTests) Test_New_Create_Product_Should_Return_Error_For_Empty_Name() {
	command, err := createProductCommand.NewCreateProduct(
		"",
		gofakeit.EmojiDescription(),
		120,
	)

	c.Require().Error(err)
	c.Assert().Nil(command)
}

func (c *createProductUnitTests) Test_New_Create_Product_Should_Return_Error_For_Empty_Description() {
	command, err := createProductCommand.NewCreateProduct(
		gofakeit.Name(),
		"",
		120,
	)

	c.Require().Error(err)
	c.Assert().Nil(command)
}
