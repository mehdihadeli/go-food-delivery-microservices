package redis

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config/environment"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/external/fxlog"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/zap"
	redis2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/redis"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func Test_Custom_Redis_Container(t *testing.T) {
	ctx := context.Background()
	var redisClient redis.UniversalClient

	fxtest.New(t,
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		core.Module,
		redis2.Module,
		fx.Decorate(RedisContainerOptionsDecorator(t, ctx)),
		fx.Populate(&redisClient),
	).RequireStart()

	assert.NotNil(t, redisClient)
}
