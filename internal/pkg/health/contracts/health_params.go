package contracts

import (
	"go.uber.org/fx"
)

type HealthParams struct {
	fx.In

	Healths []Health `group:"healths"`
}
