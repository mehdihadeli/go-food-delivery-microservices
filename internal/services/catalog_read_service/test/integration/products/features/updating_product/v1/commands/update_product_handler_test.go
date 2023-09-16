//go:build integration
// +build integration

package commands

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/updating_products/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/shared/test_fixture/integration"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUpdateProduct(t *testing.T) {
	integrationTestSharedFixture := integration.NewIntegrationTestSharedFixture(t)

	Convey("Updating Product Feature", t, func() {
		ctx := context.Background()
		integrationTestSharedFixture.InitializeTest()

		// https://specflow.org/learn/gherkin/#learn-gherkin
		// scenario
		Convey("Updating an existing product in the database", func() {
			Convey("Given an existing product in the system", func() {
				productId, err := uuid.FromString(integrationTestSharedFixture.Items[0].ProductId)
				So(err, ShouldBeNil)

				updateProduct, err := commands.NewUpdateProduct(
					productId,
					gofakeit.Name(),
					gofakeit.AdjectiveDescriptive(),
					gofakeit.Price(150, 6000),
				)
				So(err, ShouldBeNil)

				Convey("When a UpdateProduct command executed for a existing product", func() {
					result, err := mediatr.Send[*commands.UpdateProduct, *mediatr.Unit](ctx, updateProduct)

					Convey("Then the product should be updated successfully", func() {
						// Assert that the error is nil (indicating success).
						So(err, ShouldBeNil)
						So(result, ShouldNotBeNil)

						Convey("And the updated product details should be reflected in the system", func() {
							// Fetch the updated product from the database.
							updatedProduct, _ := integrationTestSharedFixture.ProductRepository.GetProductByProductId(
								ctx,
								productId.String(),
							)

							Convey("And the product's properties should match the updated data", func() {
								// Assert that the product properties match the updated data.
								So(updatedProduct.Name, ShouldEqual, updatedProduct.Name)
								So(updatedProduct.Price, ShouldEqual, updatedProduct.Price)
								// Add more assertions as needed for other properties.
							})
						})
					})
				})
			})
		})

		integrationTestSharedFixture.DisposeTest()
	})
}
