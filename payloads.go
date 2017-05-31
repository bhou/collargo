package collargo

import (
	"errors"
	"strconv"
)

// IntPayload int type payload
type IntPayload struct {
	Value int
}

// Convert convert interface{} to int payload
func (payload *IntPayload) Convert(v interface{}) error {
	switch value := v.(type) {
	case int:
		payload.Value = value
	case float64:
		payload.Value = int(value)
	case string:
		// v is a string here, so e.g. v + " Yeah!" is possible.
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		payload.Value = intValue
	default:
		// And here I'm feeling dumb. ;)
		if intValue, ok := v.(int); ok {
			payload.Value = intValue
		} else {
			return errors.New("Failed to convert payload to int type")
		}
	}
	return nil
}

// Float64Payload float type payload
type Float64Payload struct {
	Value float64
}

// Convert convert interface{} to int payload
func (payload *Float64Payload) Convert(v interface{}) error {
	switch value := v.(type) {
	case int:
		payload.Value = float64(value)
	case float64:
		payload.Value = value
	case string:
		// v is a string here, so e.g. v + " Yeah!" is possible.
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		payload.Value = floatValue
	default:
		if floatValue, ok := v.(float64); ok {
			payload.Value = floatValue
		} else {
			return errors.New("Failed to convert payload to float64 type")
		}
	}
	return nil
}
