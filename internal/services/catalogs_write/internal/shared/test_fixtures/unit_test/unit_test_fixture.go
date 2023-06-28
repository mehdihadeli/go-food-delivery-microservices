package unit_test

import (
	"context"
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
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/contracts/data"
	dto "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/dto/v1"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/mocks"
)

type UnitTestSharedFixture struct {
	Cfg *config.AppOptions
	Log logger.Logger
	suite.Suite
}

type UnitTestMockFixture struct {
	Uow               *mocks.CatalogUnitOfWork
	ProductRepository *mocks.ProductRepository
	Bus               *mocks3.Bus
	Tracer            trace.Tracer
	Ctx               context.Context
}

func NewUnitTestSharedFixture(t *testing.T) *UnitTestSharedFixture {
	// we could use EmptyLogger if we don't want to log anything
	defaultLogger.SetupDefaultLogger()
	log := defaultLogger.Logger
	cfg := &config.AppOptions{}

	err := configMapper()
	require.NoError(t, err)

	unit := &UnitTestSharedFixture{
		Cfg: cfg,
		Log: log,
	}

	return unit
}

func NewUnitTestMockFixture(t *testing.T) *UnitTestMockFixture {
	ctx, cancel := context.WithCancel(context.Background())

	t.Cleanup(func() {
		// https://dev.to/mcaci/how-to-use-the-context-done-method-in-go-22me
		// https://www.digitalocean.com/community/tutorials/how-to-use-contexts-in-go
		cancel()
	})

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

	// empty tracer, just for testing
	nopetracer := trace.NewNoopTracerProvider()
	testTracer := nopetracer.Tracer("test_tracer")

	return &UnitTestMockFixture{
		Ctx:               ctx,
		Bus:               bus,
		ProductRepository: productRepository,
		Uow:               uow,
		Tracer:            testTracer,
	}
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
