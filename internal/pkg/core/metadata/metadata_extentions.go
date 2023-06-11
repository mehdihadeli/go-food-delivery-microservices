package metadata

import (
	"time"

	"github.com/goccy/go-json"
)

func (m Metadata) GetString(key string) string {
	val, ok := m.Get(key).(string)
	if ok {
		return val
	}

	return ""
}

func (m Metadata) GetTime(key string) time.Time {
	val, ok := m.Get(key).(time.Time)
	if ok {
		return val
	}

	return *new(time.Time)
}

func (m Metadata) ToJson() string {
	marshal, err := json.Marshal(m)
	if err != nil {
		return ""
	}

	return string(marshal)
}
