package collargo

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFilter(t *testing.T) {
	// devtoolAddon := CreateDevToolAddon("ws://localhost:7500/app")
	// Collar.Use(devtoolAddon)

	ns := Collar.NS("com.collargo.test", map[string]string{})

	input := ns.Input("input")

	input.When("even", func(s Signal) (bool, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		return v.Value%2 == 0, nil
	}).Map("x 2", func(s Signal) (Signal, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		return s.New(v.Value * 2), nil
	}).Do("test", func(s Signal) (interface{}, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		assert.Equal(t, 20, v.Value)
		return v.Value, nil
	})

	input.When("odd", func(s Signal) (bool, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		return v.Value%2 != 0, nil
	}).Map("+ 1", func(s Signal) (Signal, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		return s.New(v.Value + 1), nil
	}).Do("test", func(s Signal) (interface{}, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		assert.Equal(t, 11, v.Value)
		return v.Value, nil
	})

	input.Push(10)

	time.Sleep(testDelay * time.Millisecond)
}

func TestErrorInFilter(t *testing.T) {
	ns := Collar.NS("com.collargo.test", map[string]string{})

	input := ns.Input("input")

	input.When("even", func(s Signal) (bool, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		if v.Value%2 != 0 {
			return false, errors.New("Input must be an even interger")
		}
		return true, nil
	}).Map("x 2", func(s Signal) (Signal, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		return s.New(v.Value * 2), nil
	}).Do("test", func(s Signal) (interface{}, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		assert.Fail(t, "should not go here")
		return v.Value, nil
	})

	input.Push(11)

	time.Sleep(testDelay * time.Millisecond)
}

func TestFilterErrorSignal(t *testing.T) {
	ns := Collar.NS("com.collargo.test", map[string]string{})

	input := ns.Input("input")

	errorGen := ns.Map("error generator", func(s Signal) (Signal, error) {
		newSig := s.SetError(errors.New("error when processing"))
		return newSig, nil
	})

	// input.Do("generate error", func(s Signal) (interface{}, error) {
	// return nil, errors.New("error when processing")
	// })
	input.When("even", func(s Signal) (bool, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		return v.Value%2 == 0, nil
	}).Map("x 2", func(s Signal) (Signal, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		return s.New(v.Value * 2), nil
	}).Do("test", func(s Signal) (interface{}, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		assert.Fail(t, "should not go here")
		return v.Value, nil
	}).Errors("error handling", func(s Signal, rethrow SendSignalFunc) error {
		assert.Equal(t, "error when processing", s.Error.Error())
		return nil
	})

	// input.Push(11)

	errorGen.To("flow", input)

	errorGen.Push(11)

	time.Sleep(testDelay * time.Millisecond)

}

func TestFilterErrorSignalWithActuator(t *testing.T) {
	ns := Collar.NS("com.collargo.test", map[string]string{})

	input := ns.Input("input")

	input.Do("generate error", func(s Signal) (interface{}, error) {
		return nil, errors.New("error when processing")
	}).When("even", func(s Signal) (bool, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		return v.Value%2 == 0, nil
	}).Map("x 2", func(s Signal) (Signal, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		return s.New(v.Value * 2), nil
	}).Do("test", func(s Signal) (interface{}, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		assert.Fail(t, "should not go here")
		return v.Value, nil
	}).Errors("error handling", func(s Signal, rethrow SendSignalFunc) error {
		assert.Equal(t, "error when processing", s.Error.Error())
		return nil
	})

	input.Push(11)

	time.Sleep(testDelay * time.Millisecond)

}
