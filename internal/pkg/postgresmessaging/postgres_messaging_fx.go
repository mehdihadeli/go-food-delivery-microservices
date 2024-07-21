package postgresmessaging

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/persistmessage"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresmessaging/messagepersistence"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Module(
	"postgresmessagingfx",
	fx.Provide(
		messagepersistence.NewPostgresMessagePersistenceDBContext,
		messagepersistence.NewPostgresMessageService,
	),
	fx.Invoke(migrateMessaging),
)

func migrateMessaging(db *gorm.DB) error {
	err := db.Migrator().AutoMigrate(&persistmessage.StoreMessage{})

	return err
}
