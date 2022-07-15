package infrastructure

import (
	"context"
	"github.com/go-redis/redis/v8"
	redis_client "github.com/mehdihadeli/store-golang-microservice-sample/pkg/redis"
)

func (ic *infrastructureConfigurator) configRedis(ctx context.Context) (redis.UniversalClient, error, func()) {
	rd := redis_client.NewUniversalRedisClient(ic.cfg.Redis)
	ic.log.Infof("Redis connected: %+v", rd.PoolStats())
	return rd, nil, func() {
		defer rd.Close() // nolint: errcheck
	}
}
