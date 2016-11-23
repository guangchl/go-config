package goconf

import (
	"fmt"
)

// Return this error when the searched key doesn't present.
type NoSuchKeyError struct {
	key string
}

func (err NoSuchKeyError) Error() string {
	return fmt.Sprintf("Cannot find key: %s", err.key)
}

// Return this error when the value type doesn't match.
type InvalidValueTypeError struct {
	key   string
	value interface{}
}

func (err InvalidValueTypeError) Error() string {
	return fmt.Sprintf("Invalid value type, key: %s, value: %v (%T)",
		err.key, err.value, err.value)
}
