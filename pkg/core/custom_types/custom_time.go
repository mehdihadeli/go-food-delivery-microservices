package customTypes

//https://gist.github.com/lummie/2cd6240378372079a8be7df782b380fc
//https://stackoverflow.com/questions/25845172/parsing-rfc-3339-iso-8601-date-time-string-in-go

import (
	"fmt"
	"github.com/araddon/dateparse"
	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	"strings"
	"time"
)

// CustomTime provides an example of how to declare a new time Type with a custom formatter.
// Note that time.Time methods are not available, if needed you can add and cast like the String method does
// Otherwise, only use in the json struct at marshal/unmarshal time.
type CustomTime time.Time

// UnmarshalJSON Parses the json string in the custom format
func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	//nt, err := time.Parse(time.RFC3339, s)
	nt, err := dateparse.ParseLocal(s)
	if err != nil {
		return httpErrors.NewBadRequestErrorWrap(err, fmt.Sprintf("invalid time format: %s", s))
	}
	*ct = CustomTime(nt)
	return
}

// MarshalJSON writes a quoted string in the custom format
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(ct.String()), nil
}

// String returns the time in the custom format
func (ct *CustomTime) String() string {
	t := time.Time(*ct)
	return fmt.Sprintf("%q", t)
}
