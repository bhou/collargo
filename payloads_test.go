package collargo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIntPayload(t *testing.T) {
	// devtoolAddon := CreateDevToolAddon("ws://localhost:7500/app")
	// Collar.Use(devtoolAddon)

	ns := Collar.NS("com.collargo.test", map[string]string{})

	input := ns.Input("input")

	input.Map("x2", func(s Signal) (Signal, error) {
		fmt.Println("step 1.1", "x2", s.Payload)
		v := new(IntPayload)
		s.GetValue(AnonPayload, v)
		newS := s.New(v.Value * 2)
		fmt.Println("step 1.2", newS.Payload)
		return newS, nil
	}).Map("+1", func(s Signal) (Signal, error) {
		fmt.Println("step 1.3", "+1", s.Payload)
		v := new(IntPayload)
		s.GetValue(AnonPayload, v)
		newS := s.New(v.Value + 1)
		fmt.Println("step 1.4", newS.Payload)
		return newS, nil
	}).Do("test", func(s Signal) (interface{}, error) {
		v := new(IntPayload)
		s.GetValue(AnonPayload, v)
		fmt.Println("step 1.5", v)
		assert.Equal(t, 21, v.Value)
		return "", nil
	}).Output("output")

	input.Push(10)
	time.Sleep(3000 * time.Millisecond)
}
