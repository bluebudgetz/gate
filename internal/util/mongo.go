package util

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoReturnDocAfter = options.After
)

func OptionalObjectID(value interface{}) *primitive.ObjectID {
	if value == nil {
		return nil
	} else if objectID, ok := value.(primitive.ObjectID); ok {
		return &objectID
	} else if stringPtrValue, ok := value.(*string); ok {
		if stringPtrValue == nil {
			return nil
		} else if objectID, err := primitive.ObjectIDFromHex(*stringPtrValue); err != nil {
			return nil
		} else {
			return &objectID
		}
	} else if stringValue, ok := value.(string); ok {
		if objectID, err := primitive.ObjectIDFromHex(stringValue); err != nil {
			return nil
		} else {
			return &objectID
		}
	} else {
		return nil
	}
}

func OptionalObjectIDHex(value interface{}) *string {
	if value == nil {
		return nil
	} else if objectID, ok := value.(primitive.ObjectID); ok {
		hex := objectID.Hex()
		return &hex
	} else if stringPtrValue, ok := value.(*string); ok {
		return stringPtrValue
	} else if stringValue, ok := value.(string); ok {
		return &stringValue
	} else {
		return nil
	}
}

func OptionalDateTime(value interface{}) *time.Time {
	if value == nil {
		return nil
	} else if d, ok := value.(primitive.DateTime); ok {
		dt := time.Unix(int64(d)/1000, 0)
		return &dt
	} else if d, ok := value.(time.Time); ok {
		return &d
	} else if d, ok := value.(*time.Time); ok {
		return d
	} else {
		return nil
	}
}

func MustDateTime(value interface{}) time.Time {
	dateTime := OptionalDateTime(value)
	if dateTime == nil {
		panic("date-time value is nil but required")
	}
	return *dateTime
}
