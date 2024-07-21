//go:build unit
// +build unit

package v1

import (
	"fmt"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/cqrs"
	createProductCommand "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"

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

	createProduct, err := createProductCommand.NewCreateProductWithValidation(
		name,
		description,
		price,
	)
	var g interface{} = createProduct
	d, ok := g.(cqrs.Command)
	if ok {
		fmt.Println(d)
	}

	c.Assert().NotNil(createProduct)
	c.Assert().Equal(name, createProduct.Name)
	c.Assert().Equal(price, createProduct.Price)

	c.Require().NoError(err)
}

func (c *createProductUnitTests) Test_New_Create_Product_Should_Return_Error_For_Invalid_Price() {
	command, err := createProductCommand.NewCreateProductWithValidation(
		gofakeit.Name(),
		gofakeit.EmojiDescription(),
		0,
	)

	c.Require().Error(err)
	c.Assert().Nil(command)
}

func (c *createProductUnitTests) Test_New_Create_Product_Should_Return_Error_For_Empty_Name() {
	command, err := createProductCommand.NewCreateProductWithValidation(
		"",
		gofakeit.EmojiDescription(),
		120,
	)

	c.Require().Error(err)
	c.Assert().Nil(command)
}

func (c *createProductUnitTests) Test_New_Create_Product_Should_Return_Error_For_Empty_Description() {
	command, err := createProductCommand.NewCreateProductWithValidation(
		gofakeit.Name(),
		"",
		120,
	)

	c.Require().Error(err)
	c.Assert().Nil(command)
}
