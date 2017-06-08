package collargo

import (
	"errors"
	// "fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRethrowError(t *testing.T) {
	ns := Collar.NS("com.collargo.test", map[string]string{})

	input := ns.Input("input")

	input.
		Do("error generator", func(s Signal) (interface{}, error) {
			return nil, errors.New("this is an error")
		}).
		Errors("handle error", func(s Signal, rethrow SendSignalFunc) error {
			newSignal := s.New(map[string]interface{}{
				"value": 100,
			})
			rethrow(newSignal)
			return nil
		}).
		Do("handle signal", func(s Signal) (interface{}, error) {
			v := new(IntPayload)
			s.GetValue("value", v)
			assert.Equal(t, 100, v.Value)
			return nil, nil
		})

	input.Push(1)

	time.Sleep(testDelay * time.Millisecond)
}

func TestProcessingError(t *testing.T) {
	ns := Collar.NS("com.collargo.test", map[string]string{})

	input := ns.Input("input")

	input.
		Errors("handle error", func(s Signal, rethrow SendSignalFunc) error {
			return nil
		}).
		Do("handle signal", func(s Signal) (interface{}, error) {
			v := new(IntPayload)
			s.GetValue("__anon__", v)
			assert.Equal(t, 100, v.Value)
			return nil, nil
		})

	input.Push(100)

	time.Sleep(testDelay * time.Millisecond)
}
