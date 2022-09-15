package streamName

import (
	"fmt"
	"github.com/goccy/go-reflect"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/models"
	uuid "github.com/satori/go.uuid"
	"strings"
)

type StreamName string

func (n StreamName) GetId() uuid.UUID {
	name := n.String()
	index := strings.Index(name, "-") + 1
	id := name[index:]

	return uuid.FromStringOrNil(id)
}

func (n StreamName) String() string {
	return string(n)
}

// For gets stream name for Aggregate
func For[T models.IHaveEventSourcedAggregate](aggregate T) StreamName {
	var aggregateName string
	if t := reflect.TypeOf(aggregate); t.Kind() == reflect.Ptr {
		aggregateName = reflect.TypeOf(aggregate).Elem().Name()
	} else {
		aggregateName = reflect.TypeOf(aggregate).Name()
	}

	return StreamName(fmt.Sprintf("%s-%s", strings.ToLower(aggregateName), aggregate.Id().String()))
}

// ForID gets stream name for AggregateID
func ForID[T models.IHaveEventSourcedAggregate](aggregateID uuid.UUID) StreamName {
	var aggregate T
	var aggregateName string
	if t := reflect.TypeOf(aggregate); t.Kind() == reflect.Ptr {
		aggregateName = reflect.TypeOf(aggregate).Elem().Name()
	} else {
		aggregateName = reflect.TypeOf(aggregate).Name()
	}

	return StreamName(fmt.Sprintf("%s-%s", strings.ToLower(aggregateName), aggregateID.String()))
}
