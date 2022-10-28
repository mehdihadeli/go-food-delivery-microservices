package unit_test

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/data"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	defaultLogger "github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/default_logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	mocks3 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/mocks"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto"
	mocks2 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/mocks"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UnitTestSharedFixture struct {
	Cfg *config.Config
	Log logger.Logger
	suite.Suite
}

type UnitTestMockFixture struct {
	Uow               *mocks2.CatalogsUnitOfWorks
	ProductRepository *mocks2.ProductRepository
	Bus               *mocks3.Bus
	Ctx               context.Context
}

func NewUnitTestSharedFixture(t *testing.T) *UnitTestSharedFixture {
	// we could use EmptyLogger if we don't want to log anything
	log := defaultLogger.Logger
	cfg := &config.Config{}

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
		//https://dev.to/mcaci/how-to-use-the-context-done-method-in-go-22me
		//https://www.digitalocean.com/community/tutorials/how-to-use-contexts-in-go
		cancel()
	})

	// create new mocks
	productRepository := &mocks2.ProductRepository{}
	bus := &mocks3.Bus{}
	uow := &mocks2.CatalogsUnitOfWorks{}

	//// or just clear the mocks
	//c.Bus.ExpectedCalls = nil
	//c.Bus.Calls = nil
	//c.Uow.ExpectedCalls = nil
	//c.Uow.Calls = nil
	//c.ProductRepository.ExpectedCalls = nil
	//c.ProductRepository.Calls = nil

	uow.On("Products").Return(productRepository)

	var mockUOW *mock.Call
	mockUOW = uow.On("SaveWithTx", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			fn, ok := args.Get(1).(data.UnitOfWorkActionFunc)
			if !ok {
				panic("argument mismatch")
			}

			mockUOW.Return(fn())
		})

	mockUOW.Times(1)

	bus.On("PublishMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	return &UnitTestMockFixture{
		Ctx:               ctx,
		Bus:               bus,
		ProductRepository: productRepository,
		Uow:               uow,
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
