package cqrs

import (
	"fmt"

	"go.uber.org/fx"
)

// when we register multiple handlers with output type `mediatr.RequestHandler` we get exception `type already provided`, so we should use tags annotation

// AsHandler annotates the given constructor to state that
// it provides a handler to the "handlers" group.
func AsHandler(handler interface{}, handlerGroupName string) interface{} {
	return fx.Annotate(
		handler,
		fx.As(new(HandlerRegisterer)),
		fx.ResultTags(fmt.Sprintf(
			`group:"%s"`,
			handlerGroupName,
		)),
	)
}
