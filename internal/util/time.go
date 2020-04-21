package util

import "time"

func OptionalDateTime(value interface{}) *time.Time {
	if value == nil {
		return nil
	} else if d, ok := value.(time.Time); ok {
		return &d
	} else if d, ok := value.(*time.Time); ok {
		return d
	} else {
		return nil
	}
}

func MustParseTimeRFC3339(value string) time.Time {
	if v, err := time.Parse(time.RFC3339, value); err != nil {
		panic(err)
	} else {
		return v
	}
}
