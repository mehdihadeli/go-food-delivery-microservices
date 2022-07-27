package streamName

import (
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts"
	uuid "github.com/satori/go.uuid"
	"reflect"
	"strings"
)

// For gets stream name for Aggregate
func For[T contracts.IEventSourcedAggregateRoot](aggregate T) string {
	var aggregateName string
	if t := reflect.TypeOf(aggregate); t.Kind() == reflect.Ptr {
		aggregateName = reflect.TypeOf(aggregate).Elem().Name()
	} else {
		aggregateName = reflect.TypeOf(aggregate).Name()
	}

	return fmt.Sprintf("%s-%s", strings.ToLower(aggregateName), aggregate.GetID().String())
}

// ForID gets stream name for AggregateID
func ForID[T contracts.IEventSourcedAggregateRoot](aggregateID uuid.UUID) string {
	var aggregate T
	var aggregateName string
	if t := reflect.TypeOf(aggregate); t.Kind() == reflect.Ptr {
		aggregateName = reflect.TypeOf(aggregate).Elem().Name()
	} else {
		aggregateName = reflect.TypeOf(aggregate).Name()
	}

	return fmt.Sprintf("%s-%s", strings.ToLower(aggregateName), aggregateID.String())
}
