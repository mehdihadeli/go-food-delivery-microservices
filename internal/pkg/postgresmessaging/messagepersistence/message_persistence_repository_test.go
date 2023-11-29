package messagepersistence

import (
	"context"
	"testing"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/messaging/persistmessage"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	defaultLogger "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/defaultlogger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgresgorm"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgresgorm/helpers"
	gorm2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/testcontainer/gorm"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type messagePersistenceRepositoryTest struct {
	suite.Suite
	DB                  *gorm.DB
	logger              logger.Logger
	messagingRepository persistmessage.MessagePersistenceRepository
	dbContext           *PostgresMessagePersistenceDBContext
	storeMessages       []*persistmessage.StoreMessage
	ctx                 context.Context
}

func TestMessagePersistenceRepository(t *testing.T) {
	suite.Run(
		t,
		&messagePersistenceRepositoryTest{logger: defaultLogger.GetLogger()},
	)
}

func (c *messagePersistenceRepositoryTest) SetupSuite() {
	opts, err := gorm2.NewGormTestContainers(defaultLogger.GetLogger()).
		PopulateContainerOptions(context.Background(), c.T())
	c.Require().NoError(err)

	gormDB, err := postgresgorm.NewGorm(opts)
	c.Require().NoError(err)
	c.DB = gormDB

	err = migrationDatabase(gormDB)
	c.Require().NoError(err)

	c.dbContext = NewPostgresMessagePersistenceDBContext(gormDB)
	c.messagingRepository = NewMessagePersistenceRepository(
		c.dbContext,
		defaultLogger.GetLogger(),
	)
}

func (c *messagePersistenceRepositoryTest) SetupTest() {
	ctx := context.Background()
	c.ctx = ctx
	p, err := seedData(context.Background(), c.DB)
	c.Require().NoError(err)
	c.storeMessages = p
}

func (c *messagePersistenceRepositoryTest) TearDownTest() {
	err := c.cleanupPostgresData()
	c.Require().NoError(err)
}

func (c *messagePersistenceRepositoryTest) BeginTx() {
	c.logger.Info("starting transaction")
	tx := c.dbContext.Begin()
	gormContext := helpers.SetTxToContext(c.ctx, tx)
	c.ctx = gormContext
}

func (c *messagePersistenceRepositoryTest) CommitTx() {
	tx := helpers.GetTxFromContextIfExists(c.ctx)
	if tx != nil {
		c.logger.Info("committing transaction")
		tx.Commit()
	}
}

func (c *messagePersistenceRepositoryTest) Test_Add() {
	message := &persistmessage.StoreMessage{
		ID:            uuid.NewV4(),
		MessageStatus: persistmessage.Processed,
		Data:          "test data 3",
		DataType:      "string",
		Created:       time.Now(),
		DeliveryType:  persistmessage.Outbox,
	}

	c.BeginTx()
	err := c.messagingRepository.Add(c.ctx, message)
	c.CommitTx()

	c.Require().NoError(err)

	m, err := c.messagingRepository.GetById(c.ctx, message.ID)
	if err != nil {
		return
	}

	c.Assert().NotNil(m)
	c.Assert().Equal(message.ID, m.ID)
}

func migrationDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(&persistmessage.StoreMessage{})
	if err != nil {
		return err
	}

	return nil
}

func seedData(
	ctx context.Context,
	db *gorm.DB,
) ([]*persistmessage.StoreMessage, error) {
	messages := []*persistmessage.StoreMessage{
		{
			ID:            uuid.NewV4(),
			MessageStatus: persistmessage.Processed,
			Data:          "test data",
			DataType:      "string",
			Created:       time.Now(),
			DeliveryType:  persistmessage.Outbox,
		},
		{
			ID:            uuid.NewV4(),
			MessageStatus: persistmessage.Processed,
			Data:          "test data 2",
			DataType:      "string",
			Created:       time.Now(),
			DeliveryType:  persistmessage.Outbox,
		},
	}

	err := db.WithContext(ctx).Create(messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (c *messagePersistenceRepositoryTest) cleanupPostgresData() error {
	tables := []string{"store_messages"}
	// Iterate over the tables and delete all records
	for _, table := range tables {
		err := c.DB.Exec("DELETE FROM " + table).Error

		return err
	}

	return nil
}
