package eventstroredb

import (
	"fmt"
	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
)

var ErrStreamNotFound = func(err error) httpErrors.ProblemDetailErr {
	return httpErrors.NewNotFoundError(err, "stream  doesn't not exist")
}

var ErrAggregateNotFound = func(err error, id string) httpErrors.ProblemDetailErr {
	return httpErrors.NewNotFoundError(err, fmt.Sprintf("aggregate with id '%s' doesn't not exist", id))
}

var ErrStreamAlreadyExists = func(err error, stream string) httpErrors.ProblemDetailErr {
	return httpErrors.NewConflictError(err, fmt.Sprintf("stream %s already exists", stream))
}

var ErrEventStore = func(err error, messgae string) httpErrors.ProblemDetailErr {
	return httpErrors.NewInternalServerError(err, messgae)
}

var ErrAppendToStream = func(err error, stream string) httpErrors.ProblemDetailErr {
	return httpErrors.NewBadRequestError(err, fmt.Sprintf("unable to append events to stream %s", stream))
}

var ErrReadFromStream = func(err error) httpErrors.ProblemDetailErr {
	return httpErrors.NewInternalServerError(err, "unable to read events from stream")
}

var ErrDeleteStream = func(err error, stream string) httpErrors.ProblemDetailErr {
	return httpErrors.NewInternalServerError(err, fmt.Sprintf("unable to delete stream %s", stream))
}

var ErrTruncateStream = func(err error, stream string) httpErrors.ProblemDetailErr {
	return httpErrors.NewInternalServerError(err, fmt.Sprintf("unable to truncate stream %s", stream))
}
