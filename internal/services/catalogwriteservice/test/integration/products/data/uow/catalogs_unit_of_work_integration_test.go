//go:build integration
// +build integration

package uow

import (
	"context"
	"testing"
	"time"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"
	data2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/contracts"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/models"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"

	"emperror.dev/errors"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var integrationFixture *integration.IntegrationTestSharedFixture

func TestUnitOfWork(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewIntegrationTestSharedFixture(t)
	RunSpecs(t, "CatalogsUnitOfWork Integration Tests")
}

var _ = Describe("CatalogsUnitOfWork Feature", func() {
	// Define variables to hold repository and product data
	var (
		ctx      context.Context
		err      error
		products *utils.ListResult[*models.Product]
	)

	_ = BeforeEach(func() {
		ctx = context.Background()
		By("Seeding the required data")
		integrationFixture.SetupTest()
	})

	_ = AfterEach(func() {
		By("Cleanup test data")
		integrationFixture.TearDownTest()
	})

	// "Scenario" step for testing a UnitOfWork action that should roll back on error
	Describe("Rollback on error", func() {
		// "When" step
		When("The UnitOfWork Do executed and there is an error in the execution", func() {
			It("Should roll back the changes and not affect the database", func() {
				err = integrationFixture.CatalogUnitOfWorks.Do(ctx, func(catalogContext data2.CatalogContext) error {
					_, err := catalogContext.Products().CreateProduct(ctx,
						&models.Product{
							Name:        gofakeit.Name(),
							Description: gofakeit.AdjectiveDescriptive(),
							Id:          uuid.NewV4(),
							Price:       gofakeit.Price(100, 1000),
							CreatedAt:   time.Now(),
						})
					Expect(err).NotTo(HaveOccurred()) // Successful product creation

					return errors.New("error rollback")
				})
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring("error rollback")))

				products, err := integrationFixture.ProductRepository.GetAllProducts(ctx, utils.NewListQuery(10, 1))
				Expect(err).To(BeNil())

				Expect(len(products.Items)).To(Equal(2)) // Ensure no changes in the database
			})
		})
	})

	// "Scenario" step for testing a UnitOfWork action that should rollback on panic
	Describe("Rollback on panic", func() {
		// "When" step
		When("The UnitOfWork Do executed and there is an panic in the execution", func() {
			It("Should roll back the changes and not affect the database", func() {
				err = integrationFixture.CatalogUnitOfWorks.Do(ctx, func(catalogContext data2.CatalogContext) error {
					_, err := catalogContext.Products().CreateProduct(ctx,
						&models.Product{
							Name:        gofakeit.Name(),
							Description: gofakeit.AdjectiveDescriptive(),
							Id:          uuid.NewV4(),
							Price:       gofakeit.Price(100, 1000),
							CreatedAt:   time.Now(),
						})
					Expect(err).To(BeNil()) // Successful product creation

					panic(errors.New("panic rollback"))
				})
				Expect(err).To(HaveOccurred())

				products, err = integrationFixture.ProductRepository.GetAllProducts(ctx, utils.NewListQuery(10, 1))
				Expect(err).To(BeNil())

				Expect(len(products.Items)).To(Equal(2)) // Ensure no changes in the database
			})
		})
	})

	// "Scenario" step for testing a UnitOfWork action that should rollback when the context is canceled
	Describe("Cancelling the context", func() {
		// "When" step
		When("the UnitOfWork Do executed and cancel the context", func() {
			It("Should roll back the changes and not affect the database", func() {
				cancelCtx, cancel := context.WithCancel(ctx)

				err := integrationFixture.CatalogUnitOfWorks.Do(
					cancelCtx,
					func(catalogContext data2.CatalogContext) error {
						_, err := catalogContext.Products().CreateProduct(ctx,
							&models.Product{
								Name:        gofakeit.Name(),
								Description: gofakeit.AdjectiveDescriptive(),
								Id:          uuid.NewV4(),
								Price:       gofakeit.Price(100, 1000),
								CreatedAt:   time.Now(),
							})
						Expect(err).To(BeNil()) // Successful product creation

						_, err = catalogContext.Products().CreateProduct(ctx,
							&models.Product{
								Name:        gofakeit.Name(),
								Description: gofakeit.AdjectiveDescriptive(),
								Id:          uuid.NewV4(),
								Price:       gofakeit.Price(100, 1000),
								CreatedAt:   time.Now(),
							})
						Expect(err).To(BeNil()) // Successful product creation

						cancel() // Cancel the context

						return err
					},
				)
				Expect(err).To(HaveOccurred())

				// Validate that changes are rolled back in the database
				products, err := integrationFixture.ProductRepository.GetAllProducts(ctx, utils.NewListQuery(10, 1))
				Expect(err).To(BeNil())
				Expect(len(products.Items)).To(Equal(2)) // Ensure no changes in the database
			})
		})
	})

	// "Scenario" step for testing a UnitOfWork action that should commit on success
	Describe("Commit on success", func() {
		// "When" step
		When("the UnitOfWork Do executed and operation was successfull", func() {
			It("Should commit the changes to the database", func() {
				err := integrationFixture.CatalogUnitOfWorks.Do(ctx, func(catalogContext data2.CatalogContext) error {
					_, err := catalogContext.Products().CreateProduct(ctx,
						&models.Product{
							Name:        gofakeit.Name(),
							Description: gofakeit.AdjectiveDescriptive(),
							Id:          uuid.NewV4(),
							Price:       gofakeit.Price(100, 1000),
							CreatedAt:   time.Now(),
						})
					Expect(err).To(BeNil()) // Successful product creation

					return err
				})
				Expect(err).To(BeNil()) // No error indicates success

				// Validate that changes are committed in the database
				products, err := integrationFixture.ProductRepository.GetAllProducts(ctx, utils.NewListQuery(10, 1))
				Expect(err).To(BeNil())
				Expect(len(products.Items)).To(Equal(3)) // Ensure changes in the database
			})
		})
	})
})
