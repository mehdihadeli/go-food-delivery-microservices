package contracts

import (
	"context"

	"gorm.io/gorm"
)

type GormContext struct {
	Tx *gorm.DB
	context.Context
}
