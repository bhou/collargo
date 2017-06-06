package collargo

import (
	"fmt"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDevToolAddon(t *testing.T) {
	devtoolAddon := CreateDevToolAddon("ws://localhost:7500/app")
	Collar.Use(devtoolAddon)

	ns := Collar.NS("com.collargo.test", map[string]string{
		"module": "test",
	})

	sensor := ns.Sensor("test sensor", func(options string, send SendDataFunc) {
		time.Sleep(1000 * time.Millisecond)
		send(float64(10))
	}, false)

	sensor.Map("@double x2", func(s Signal) (Signal, error) {
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
		fmt.Println("step 1.5", v.Value)
		assert.Equal(t, 21, v.Value)
		return "", nil
	})

	sensor.
		Map("x3", func(s Signal) (Signal, error) {
			fmt.Println("step 2.1", "x3", s.Payload)
			v := new(IntPayload)
			s.GetValue(AnonPayload, v)
			newS := s.New(v.Value * 3)
			fmt.Println("step 2.2", newS.Payload)
			return newS, nil
		}).
		Map("+1", func(s Signal) (Signal, error) {
			fmt.Println("step 2.3", "+1", s.Payload)
			v := new(IntPayload)
			s.GetValue(AnonPayload, v)
			newS := s.New(v.Value + 1)
			fmt.Println("step 2.4", newS.Payload)
			return newS, nil
		}).
		Do("test", func(s Signal) (interface{}, error) {
			v := new(IntPayload)
			s.GetValue(AnonPayload, v)
			fmt.Println("step 2.5", v.Value)
			assert.Equal(t, 31, v.Value)
			return "", nil
		})

	time.Sleep(testDelay * time.Millisecond)

	// devtoolAddon.Stop()
}

func TestPrintStackTrace(t *testing.T) {
	// debug.PrintStack()
	stack := debug.Stack()
	stackSlice := strings.Split(string(stack), "\n")

	for i := range stackSlice {
		// lines := strings.Split(stackSlice[i], " ")
		// fmt.Println(lines[0])
		fmt.Println(stackSlice[i])
	}

	assert.True(t, true)
}

func TestHandleNodeAndEdge(t *testing.T) {
	ns := Collar.NS("com.collargo.test", map[string]string{
		"module": "test",
	})

	n1 := ns.Map("@passthrough processor test", func(s Signal) (Signal, error) {
		return s, nil
	})

	n2 := ns.Do("@print print", func(s Signal) (interface{}, error) {
		fmt.Println(s.Payload)
		return nil, nil
	})

	elem := handleNode(n1)
	edge := handleEdge(n1, n2)

	assert.Equal(t, "processor", elem.Data.Model)
	assert.Equal(t, "processor test", elem.Data.Label)
	assert.Equal(t, "com.collargo.test.passthrough", elem.Data.FullName)
	assert.Equal(t, "", elem.Classes)

	assert.Equal(t, n1.ID(), edge.Data.Source)
	assert.Equal(t, n2.ID(), edge.Data.Target)
}
