package collargo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	// "time"
)

func TestToFlowFunc(t *testing.T) {
	ns := Collar.NS("com.collargo.test", map[string]string{})

	input := ns.Input("input")
	output := ns.Output("output")

	input.Map("x2", func(s Signal) (Signal, error) {
		v := new(IntPayload)
		s.GetValue(AnonPayload, v)
		newS := s.New(v.Value * 2)
		return newS, nil
	}).Map("+1", func(s Signal) (Signal, error) {
		v := new(IntPayload)
		s.GetValue(AnonPayload, v)
		newS := s.New(v.Value + 1)
		return newS, nil
	}).Do("test", func(s Signal) (interface{}, error) {
		v := new(IntPayload)
		s.GetValue(AnonPayload, v)
		fmt.Println(v.Value)
		return "", nil
	}).To("output", output)

	flowFunc := Collar.ToFlowFunc(input, output)

	r, err := flowFunc(10)
	v, _ := r["__anon__"]
	assert.Nil(t, err)
	assert.Equal(t, 21, v.(int))
	fmt.Println("assert", v)
}
