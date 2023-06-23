package integration

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/stretchr/testify/suite"
	"gopkg.in/khaiql/dbcleaner.v2"
	"gopkg.in/khaiql/dbcleaner.v2/engine"
	"gorm.io/gorm"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"
	config2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/contracts/data"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/app/test"
)

type IntegrationTestSharedFixture struct {
	Cfg                *config.AppOptions
	Log                logger.Logger
	Bus                bus.Bus
	CatalogUnitOfWorks data.CatalogUnitOfWork
	ProductRepository  data.ProductRepository
	suite.Suite
	Container       contracts.Container
	DbCleaner       dbcleaner.DbCleaner
	RabbitmqCleaner *rabbithole.Client
	rabbitmqOptions *config2.RabbitmqOptions
	Gorm            *gorm.DB
}

func NewIntegrationTestSharedFixture(t *testing.T) *IntegrationTestSharedFixture {
	result := test.NewTestApp().Run(t)

	// https://github.com/michaelklishin/rabbit-hole
	rmqc, _ := rabbithole.NewClient(
		fmt.Sprintf(result.RabbitmqOptions.RabbitmqHostOptions.HttpEndPoint()),
		result.RabbitmqOptions.RabbitmqHostOptions.UserName,
		result.RabbitmqOptions.RabbitmqHostOptions.Password)

	// https://github.com/khaiql/dbcleaner
	postgresEngine := engine.NewPostgresEngine(result.GormOptions.Dns())
	postgresCleaner := dbcleaner.New()
	postgresCleaner.SetEngine(postgresEngine)

	shared := &IntegrationTestSharedFixture{
		Log:                result.Logger,
		Container:          result.Container,
		Cfg:                result.Cfg,
		RabbitmqCleaner:    rmqc,
		DbCleaner:          postgresCleaner,
		ProductRepository:  result.ProductRepository,
		CatalogUnitOfWorks: result.CatalogUnitOfWorks,
		Bus:                result.Bus,
		rabbitmqOptions:    result.RabbitmqOptions,
		Gorm:               result.Gorm,
	}

	return shared
}

func (i *IntegrationTestSharedFixture) CleanupRabbitmqData() error {
	// https://github.com/michaelklishin/rabbit-hole
	// Get all queues
	queues, err := i.RabbitmqCleaner.ListQueuesIn(i.rabbitmqOptions.RabbitmqHostOptions.VirtualHost)
	if err != nil {
		return err
	}

	// clear each queue
	for _, queue := range queues {
		_, err = i.RabbitmqCleaner.PurgeQueue(
			i.rabbitmqOptions.RabbitmqHostOptions.VirtualHost,
			queue.Name,
		)
		i.Require().NoError(err)
	}

	return nil
}

func (i *IntegrationTestSharedFixture) CleanupPostgresData() {
	tables := []string{"products"}
	err := cleanupTables(i.Gorm, tables)
	i.Require().NoError(err)
}

func cleanupTables(db *gorm.DB, tables []string) error {
	// Iterate over the tables and delete all records
	for _, table := range tables {
		err := db.Exec("DELETE FROM " + table).Error
		if err != nil {
			return err
		}
	}
	return nil
}
