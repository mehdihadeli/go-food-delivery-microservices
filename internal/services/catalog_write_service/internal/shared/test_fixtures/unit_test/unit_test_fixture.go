package unit_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	defaultLogger "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/default_logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mapper"
	mocks3 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/mocks"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts/data"
	dto "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dto/v1"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/mocks/testData"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/mocks"
)

type UnitTestSharedFixture struct {
	Cfg *config.AppOptions
	Log logger.Logger
	suite.Suite
	Items             []*models.Product
	Uow               *mocks.CatalogUnitOfWork
	ProductRepository *mocks.ProductRepository
	Bus               *mocks3.Bus
	Tracer            trace.Tracer
}

func NewUnitTestSharedFixture(t *testing.T) *UnitTestSharedFixture {
	// we could use EmptyLogger if we don't want to log anything
	defaultLogger.SetupDefaultLogger()
	log := defaultLogger.Logger
	cfg := &config.AppOptions{}

	err := configMapper()
	require.NoError(t, err)

	// empty tracer, just for testing
	nopetracer := trace.NewNoopTracerProvider()
	testTracer := nopetracer.Tracer("test_tracer")

	unit := &UnitTestSharedFixture{
		Cfg:    cfg,
		Log:    log,
		Items:  testData.Products,
		Tracer: testTracer,
	}

	return unit
}

func configMapper() error {
	err := mapper.CreateMap[*models.Product, *dto.ProductDto]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*dto.ProductDto, *models.Product]()
	if err != nil {
		return err
	}

	return nil
}

// //////////////Shared Hooks////////////////
func (c *UnitTestSharedFixture) SetupTest() {
	// create new mocks
	productRepository := &mocks.ProductRepository{}
	bus := &mocks3.Bus{}
	uow := &mocks.CatalogUnitOfWork{}
	catalogContext := &mocks.CatalogContext{}

	//// or just clear the mocks
	//c.Bus.ExpectedCalls = nil
	//c.Bus.Calls = nil
	//c.Uow.ExpectedCalls = nil
	//c.Uow.Calls = nil
	//c.ProductRepository.ExpectedCalls = nil
	//c.ProductRepository.Calls = nil

	uow.On("Products").Return(productRepository)
	catalogContext.On("Products").Return(productRepository)

	var mockUOW *mock.Call
	mockUOW = uow.On("Do", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			fn, ok := args.Get(1).(data.CatalogUnitOfWorkActionFunc)
			if !ok {
				panic("argument mismatch")
			}
			fmt.Println(fn)

			mockUOW.Return(fn(catalogContext))
		})

	mockUOW.Times(1)
	bus.On("PublishMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	c.Uow = uow
	c.ProductRepository = productRepository
	c.Bus = bus
}

func (c *UnitTestSharedFixture) CleanupMocks() {
	c.SetupTest()
}

func (c *UnitTestSharedFixture) TearDownSuite() {
	mapper.ClearMappings()
}
