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
