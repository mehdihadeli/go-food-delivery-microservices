//go:build unit
// +build unit

package v1

import (
	"testing"
	"time"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/cqrs"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm/gormdbcontext"
	datamodels "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/data/datamodels"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	creatingproductv1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1"
	creatingproductdtosv1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1/dtos"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/models"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"

	"emperror.dev/errors"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type createProductHandlerUnitTests struct {
	*unittest.UnitTestSharedFixture
	handler cqrs.RequestHandlerWithRegisterer[*creatingproductv1.CreateProduct, *creatingproductdtosv1.CreateProductResponseDto]
}

func TestCreateProductHandlerUnit(t *testing.T) {
	suite.Run(t, &createProductHandlerUnitTests{
		UnitTestSharedFixture: unittest.NewUnitTestSharedFixture(t),
	},
	)
}

func (c *createProductHandlerUnitTests) SetupTest() {
	// call base SetupTest hook before running child hook
	c.UnitTestSharedFixture.SetupTest()
	c.handler = creatingproductv1.NewCreateProductHandler(
		fxparams.ProductHandlerParams{
			CatalogsDBContext: c.CatalogDBContext,
			Tracer:            c.Tracer,
			RabbitmqProducer:  c.Bus,
			Log:               c.Log,
		},
	)
}

func (c *createProductHandlerUnitTests) TearDownTest() {
	// call base TearDownTest hook before running child hook
	c.UnitTestSharedFixture.TearDownTest()
}

func (c *createProductHandlerUnitTests) Test_Handle_Should_Create_New_Product_With_Valid_Data() {
	id := uuid.NewV4()

	createProduct := &creatingproductv1.CreateProduct{
		ProductID:   id,
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	c.BeginTx()
	_, err := c.handler.Handle(c.Ctx, createProduct)
	c.CommitTx()

	c.Require().NoError(err)

	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)

	res, err := gormdbcontext.FindModelByID[*datamodels.ProductDataModel, *models.Product](
		c.Ctx,
		c.CatalogDBContext,
		id,
	)
	c.Require().NoError(err)

	c.Assert().Equal(res.Id, id)
}

func (c *createProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_Duplicate_Item() {
	id := uuid.NewV4()

	createProduct := &creatingproductv1.CreateProduct{
		ProductID:   id,
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	c.BeginTx()
	dto, err := c.handler.Handle(c.Ctx, createProduct)
	c.Require().NoError(err)
	c.Require().NotNil(dto)
	c.CommitTx()

	c.BeginTx()
	dto, err = c.handler.Handle(c.Ctx, createProduct)
	c.CommitTx()

	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
	c.True(customErrors.IsConflictError(err))
	c.ErrorContains(err, "product already exists")
	c.Nil(dto)
}

func (c *createProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_Error_In_Bus() {
	id := uuid.NewV4()

	createProduct := &creatingproductv1.CreateProduct{
		ProductID:   id,
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	// override called mock
	// https://github.com/stretchr/testify/issues/558
	c.Bus.Mock.ExpectedCalls = nil
	c.Bus.On("PublishMessage", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return(errors.New("error in the publish message"))

	c.BeginTx()

	dto, err := c.handler.Handle(c.Ctx, createProduct)

	c.CommitTx()

	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
	c.ErrorContains(err, "error in the publish message")
	c.ErrorContains(
		err,
		"error in publishing ProductCreated integration_events event",
	)
	c.Nil(dto)
}

func (c *createProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_Error_In_Mapping() {
	id := uuid.NewV4()

	createProduct := &creatingproductv1.CreateProduct{
		ProductID:   id,
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	mapper.ClearMappings()

	c.BeginTx()

	dto, err := c.handler.Handle(c.Ctx, createProduct)

	c.CommitTx()

	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 0)
	c.ErrorContains(err, "error in the mapping")
	c.True(customErrors.IsInternalServerError(err))
	c.Nil(dto)
}
