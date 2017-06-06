package collargo

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSensorInputSignal(t *testing.T) {
	ns := Collar.NS("com.collargo.test", map[string]string{})

	input := ns.Input("input")

	sensor := ns.Sensor("sensor", func(options string, send SendDataFunc) {
	}, false)

	sensor.Map("x 2", func(s Signal) (Signal, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		return s.New(v.Value * 2), nil
	}).Do("test", func(s Signal) (interface{}, error) {
		assert.Fail(t, "should not go here")
		return nil, nil
	})

	input.To("sensor flow", sensor)

	input.Push(2)
	assert.True(t, true)

	time.Sleep(testDelay * time.Millisecond)
}

func TestSensorInputError(t *testing.T) {
	ns := Collar.NS("com.collargo.test", map[string]string{})

	input := ns.Input("input")

	errGen := ns.Do("error generator", func(s Signal) (interface{}, error) {
		return nil, errors.New("error")
	})

	sensor := ns.Sensor("sensor", func(options string, send SendDataFunc) {
	}, false)

	input.To("error generator", errGen).To(
		"sensor",
		sensor).Map("x 2", func(s Signal) (Signal, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		return s.New(v.Value * 2), nil
	}).Do("test", func(s Signal) (interface{}, error) {
		assert.Fail(t, "should not go here")
		return nil, nil
	})

	input.To("sensor flow", sensor)

	input.Push(2)
	assert.True(t, true)

	time.Sleep(testDelay * time.Millisecond)
}
