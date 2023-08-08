package hypothesis

import (
	"context"
	"time"

	"github.com/goccy/go-reflect"
	"github.com/stretchr/testify/assert"

	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"
)

type Hypothesis[T any] interface {
	Validate(ctx context.Context, message string, time time.Duration)
	Test(ctx context.Context, item T)
}

type hypothesis[T any] struct {
	data      T
	condition func(item T) bool
	t         assert.TestingT
}

func (h *hypothesis[T]) Validate(ctx context.Context, message string, time time.Duration) {
	err := testUtils.WaitUntilConditionMet(func() bool {
		if reflect.ValueOf(h.data).IsZero() {
			return false
		}
		return true
	}, time)

	if err != nil {
		assert.FailNowf(h.t, "hypothesis validation failed, %s", message)
	}
}

func (h *hypothesis[T]) Test(ctx context.Context, item T) {
	h.data = item
}
