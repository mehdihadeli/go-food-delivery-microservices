//go:build integration
// +build integration

package commands

import (
	"context"
	"testing"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/creating_product/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/creating_product/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/shared/test_fixture/integration"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateProduct(t *testing.T) {
	integrationTestSharedFixture := integration.NewIntegrationTestSharedFixture(t)

	Convey("Creating Product Feature", t, func() {
		ctx := context.Background()
		integrationTestSharedFixture.InitializeTest()

		// https://specflow.org/learn/gherkin/#learn-gherkin
		// scenario
		Convey(
			"Creating a new product and saving it to the database for a none-existing product",
			func() {
				Convey("Given new product doesn't exists in the system", func() {
					command, err := commands.NewCreateProduct(
						uuid.NewV4().String(),
						gofakeit.Name(),
						gofakeit.AdjectiveDescriptive(),
						gofakeit.Price(150, 6000),
						time.Now(),
					)
					So(err, ShouldBeNil)

					Convey(
						"When the CreateProduct command is executed and product doesn't exists",
						func() {
							result, err := mediatr.Send[*commands.CreateProduct, *dtos.CreateProductResponseDto](
								ctx,
								command,
							)

							Convey("Then the product should be created successfully", func() {
								So(err, ShouldBeNil)
								So(result, ShouldNotBeNil)

								Convey(
									"And the product ID should not be empty and same as commandId",
									func() {
										So(result.Id, ShouldEqual, command.Id)

										Convey(
											"And product detail should be retrievable from the database",
											func() {
												createdProduct, err := integrationTestSharedFixture.ProductRepository.GetProductById(
													ctx,
													result.Id,
												)
												So(err, ShouldBeNil)
												So(createdProduct, ShouldNotBeNil)
											},
										)
									},
								)
							})
						},
					)
				})
			},
		)

		integrationTestSharedFixture.DisposeTest()
	})
}
