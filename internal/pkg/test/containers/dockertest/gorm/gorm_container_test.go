package gorm

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config/environment"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/external/fxlog"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/zap"
	gormPostgres "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"
)

func Test_Gorm_Container(t *testing.T) {
	ctx := context.Background()
	var gorm *gorm.DB

	fxtest.New(t,
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		core.Module,
		gormPostgres.Module,
		fx.Decorate(GormDockerTestConatnerOptionsDecorator(t, ctx)),
		fx.Populate(&gorm),
	).RequireStart()

	assert.NotNil(t, gorm)
}
