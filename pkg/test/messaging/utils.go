package messaging

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/utils"
	typeMapper "github.com/mehdihadeli/store-golang-microservice-sample/pkg/reflection/type_mappper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/hypothesis"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/messaging/consumer"
)

func ShouldProduced[T types.IMessage](ctx context.Context, bus bus.Bus, condition func(T) bool) hypothesis.Hypothesis[T] {
	hypo := hypothesis.ForT[T](condition)

	bus.AddMessageProducedHandler(func(message types.IMessage) {
		typ := utils.GetMessageBaseReflectType(typeMapper.GenericInstanceByT[T]())
		if utils.GetMessageBaseReflectType(message) == typ {
			m, ok := message.(T)
			if !ok {
				hypo.Test(ctx, *new(T))
			}
			hypo.Test(ctx, m)
		}
	})

	return hypo
}

func ShouldConsume[T types.IMessage](ctx context.Context, bus bus.Bus, condition func(T) bool) hypothesis.Hypothesis[T] {
	hypo := hypothesis.ForT[T](condition)

	bus.AddMessageConsumedHandler(func(message types.IMessage) {
		typ := utils.GetMessageBaseReflectType(typeMapper.GenericInstanceByT[T]())
		if utils.GetMessageBaseReflectType(message) == typ {
			m, ok := message.(T)
			if !ok {
				hypo.Test(ctx, *new(T))
			}
			hypo.Test(ctx, m)
		}
	})

	return hypo
}

func ShouldConsumeNewConsumer[T types.IMessage](ctx context.Context, bus bus.Bus) (hypothesis.Hypothesis[T], error) {
	hypo := hypothesis.ForT[T](nil)
	testConsumer := consumer.NewRabbitMQFakeTestConsumerHandler(hypo)
	err := bus.ConnectConsumerHandler(typeMapper.GenericInstanceByT[T](), testConsumer)
	if err != nil {
		return nil, err
	}

	return hypo, nil
}
