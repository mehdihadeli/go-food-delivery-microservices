package eventstroredb

import (
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors/custom_errors"
)

var ErrStreamNotFound = func(err error) error {
	return customErrors.NewNotFoundErrorWrap(err, "stream  doesn't not exist")
}

var ErrAggregateNotFound = func(err error, id string) error {
	return customErrors.NewNotFoundErrorWrap(err, fmt.Sprintf("aggregate with id '%s' doesn't not exist", id))
}

var ErrAppendToStream = func(err error, stream string) error {
	return customErrors.NewBadRequestErrorWrap(err, fmt.Sprintf("unable to append events to stream %s", stream))
}

var ErrReadFromStream = func(err error) error {
	return customErrors.NewInternalServerErrorWrap(err, "unable to read events from stream")
}

var ErrDeleteStream = func(err error, stream string) error {
	return customErrors.NewInternalServerErrorWrap(err, fmt.Sprintf("unable to delete stream %s", stream))
}

var ErrTruncateStream = func(err error, stream string) error {
	return customErrors.NewInternalServerErrorWrap(err, fmt.Sprintf("unable to truncate stream %s", stream))
}
