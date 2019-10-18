package util

import (
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func OptionalObjectID(value interface{}) *string {
	if value == nil {
		return nil
	} else if objectID, ok := value.(primitive.ObjectID); ok {
		hex := objectID.Hex()
		return &hex
	}
	return nil
}

func OptionalDateTime(value interface{}) *time.Time {
	if value == nil {
		return nil
	} else if d, ok := value.(primitive.DateTime); ok {
		dt := time.Unix(int64(d)/1000, 0)
		return &dt
	}
	return nil
}

func MustDateTime(value interface{}) time.Time {
	dateTime := OptionalDateTime(value)
	if dateTime == nil {
		panic("date-time value is nil but required")
	}
	return *dateTime
}

func OptionalJsonNumber(value interface{}) *json.Number {
	if value == nil {
		return nil
	} else {
		switch val := value.(type) {
		case primitive.Decimal128:
			jsonNumber := json.Number(val.String())
			return &jsonNumber
		case float32, float64:
			jsonNumber := json.Number(fmt.Sprintf("%f", val))
			return &jsonNumber
		case int, int8, int32, int64, uint, uint8, uint32, uint64:
			jsonNumber := json.Number(fmt.Sprintf("%d", val))
			return &jsonNumber
		default:
			return nil
		}
	}
}

func MustJsonNumber(value interface{}) json.Number {
	ptr := OptionalJsonNumber(value)
	if ptr == nil {
		panic("number value is nil but required")
	} else {
		return *ptr
	}
}
