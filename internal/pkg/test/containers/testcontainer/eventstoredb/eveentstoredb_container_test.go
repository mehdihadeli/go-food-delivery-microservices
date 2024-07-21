package eventstoredb

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config/environment"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/eventstroredb"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/external/fxlog"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/zap"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func Test_Custom_EventStoreDB_Container(t *testing.T) {
	var esdbClient *esdb.Client
	ctx := context.Background()

	fxtest.New(t,
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		core.Module,
		eventstroredb.ModuleFunc(func() {
		}),
		fx.Decorate(EventstoreDBContainerOptionsDecorator(t, ctx)),
		fx.Populate(&esdbClient),
	).RequireStart()
}
