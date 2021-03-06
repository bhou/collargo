package collargo

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

/* Examples */

func ExampleNamespace() {
	// create a namespace
	ns := Collar.NS("com.collargo.test", map[string]string{})

	// create an input node
	input := ns.Input("input")

	// create a flow
	// input -> x2 -> +1 -> test -> output
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
	}).Output("output")

	// push a signal through input node
	input.Push(10)

	time.Sleep(testDelay * time.Millisecond)

	// Output:
	// 21
}

func TestSensor(t *testing.T) {
	sensor := Collar.Sensor("test sensor", func(options string, send SendDataFunc) {
		time.Sleep(1000 * time.Millisecond)
		send("text message")
	}, false)

	node1 := CreateNode("Example Node 1", "com.collartechs.test", passThroughSignalProcessor{})

	node1.Observe(func(node Node, when string, signal Signal, data ...interface{}) error {
		if when == "onReceive" {
			str, existed := signal.Get(AnonPayload)
			assert.True(t, existed)
			assert.Equal(t, "text message", str.(string))
		}
		return nil
	})

	sensor.To("node1", node1)
	fmt.Println("ready")

	time.Sleep(testDelay * time.Millisecond)
}

func TestSensorWithDeferWatch(t *testing.T) {
	sensor := Collar.Sensor("test sensor", func(options string, send SendDataFunc) {
		time.Sleep(100 * time.Millisecond)
		send("text message")
	}, true)

	node1 := CreateNode("Example Node 1", "com.collartechs.test", passThroughSignalProcessor{})

	node1.Observe(func(node Node, when string, signal Signal, data ...interface{}) error {
		if when == "onReceive" {
			str, existed := signal.Get(AnonPayload)
			assert.True(t, existed)
			assert.Equal(t, "text message", str.(string))
		}
		return nil
	})

	sensor.To("node1", node1)

	sensor.Watch("initiated")

	time.Sleep(testDelay * time.Millisecond)
}

func TestProcessor(t *testing.T) {
	sensor := Collar.Sensor("test sensor", func(options string, send SendDataFunc) {
		time.Sleep(100 * time.Millisecond)
		send(10)
	}, false)

	sensor.Map("x2", func(s Signal) (Signal, error) {
		fmt.Println("step 1", "x2", s.Payload)
		v, _ := s.Get(AnonPayload)
		newS := s.New(v.(int) * 2)
		fmt.Println("step 2", newS.Payload)
		return newS, nil
	}).Map("+1", func(s Signal) (Signal, error) {
		fmt.Println("step 3", "+1", s.Payload)
		v, _ := s.Get(AnonPayload)
		newS := s.New(v.(int) + 1)
		fmt.Println("step 4", newS.Payload)
		return newS, nil
	}).Map("test", func(s Signal) (Signal, error) {
		v, _ := s.Get(AnonPayload)
		fmt.Println("step 5", v)
		assert.Equal(t, 21, v.(int))
		return s, nil
	})

	time.Sleep(testDelay * time.Millisecond)
}

func TestActuator(t *testing.T) {
	sensor := Collar.Sensor("test sensor", func(options string, send SendDataFunc) {
		time.Sleep(100 * time.Millisecond)
		send(10)
	}, false)

	sensor.Map("x2", func(s Signal) (Signal, error) {
		fmt.Println("step 1", "x2", s.Payload)
		v, _ := s.Get(AnonPayload)
		newS := s.New(v.(int) * 2)
		fmt.Println("step 2", newS.Payload)
		return newS, nil
	}).Map("+1", func(s Signal) (Signal, error) {
		fmt.Println("step 3", "+1", s.Payload)
		v, _ := s.Get(AnonPayload)
		newS := s.New(v.(int) + 1)
		fmt.Println("step 4", newS.Payload)
		return newS, nil
	}).Do("greeting", func(s Signal) (interface{}, error) {
		fmt.Println("step 5", "greeting", s.Payload)
		v, _ := s.Get(AnonPayload)
		fmt.Println("step 6")
		return "Hello! " + strconv.Itoa(v.(int)), nil
	}).Do("test", func(s Signal) (interface{}, error) {
		v, _ := s.Get(AnonPayload)
		greeting, _ := s.GetResult()
		fmt.Println("step 7", s.Payload)
		assert.Equal(t, 21, v.(int))
		assert.Equal(t, greeting, "Hello! 21")
		return "", nil
	})

	time.Sleep(testDelay * time.Millisecond)
}

func TestErrors(t *testing.T) {
	sensor := Collar.Sensor("test sensor", func(options string, send SendDataFunc) {
		time.Sleep(100 * time.Millisecond)
		send(10)
	}, false)

	sensor.Map("x2", func(s Signal) (Signal, error) {
		fmt.Println("step 1", "x2", s.Payload)
		v, _ := s.Get(AnonPayload)
		newS := s.New(v.(int) * 2)
		fmt.Println("step 2", newS.Payload)
		return newS, errors.New("Test Error")
	}).Errors("error handler", func(s Signal, rethrow SendSignalFunc) error {
		fmt.Println("step 3", "error handling", s.Payload)
		assert.Equal(t, "Test Error", s.Error.Error())
		return nil
	})

	time.Sleep(testDelay * time.Millisecond)
}

func TestMultipleFlow(t *testing.T) {
	sensor := Collar.Sensor("test sensor", func(options string, send SendDataFunc) {
		time.Sleep(100 * time.Millisecond)
		send(10)
	}, false)

	sensor.Map("x2", func(s Signal) (Signal, error) {
		fmt.Println("step 1.1", "x2", s.Payload)
		v, _ := s.Get(AnonPayload)
		newS := s.New(v.(int) * 2)
		fmt.Println("step 1.2", newS.Payload)
		return newS, nil
	}).Map("+1", func(s Signal) (Signal, error) {
		fmt.Println("step 1.3", "+1", s.Payload)
		v, _ := s.Get(AnonPayload)
		newS := s.New(v.(int) + 1)
		fmt.Println("step 1.4", newS.Payload)
		return newS, nil
	}).Do("test", func(s Signal) (interface{}, error) {
		v, _ := s.Get(AnonPayload)
		fmt.Println("step 1.5", v)
		assert.Equal(t, 21, v.(int))
		return "", nil
	})

	sensor.Map("x3", func(s Signal) (Signal, error) {
		fmt.Println("step 2.1", "x3", s.Payload)
		v, _ := s.Get(AnonPayload)
		newS := s.New(v.(int) * 3)
		fmt.Println("step 2.2", newS.Payload)
		return newS, nil
	}).Map("+1", func(s Signal) (Signal, error) {
		fmt.Println("step 2.3", "+1", s.Payload)
		v, _ := s.Get(AnonPayload)
		newS := s.New(v.(int) + 1)
		fmt.Println("step 2.4", newS.Payload)
		return newS, nil
	}).Do("test", func(s Signal) (interface{}, error) {
		v, _ := s.Get(AnonPayload)
		fmt.Println("step 2.5", v)
		assert.Equal(t, 31, v.(int))
		return "", nil
	})

	time.Sleep(testDelay * time.Millisecond)
}

func TestEndpoint(t *testing.T) {
	// devtoolAddon := CreateDevToolAddon("ws://localhost:7500/app")
	// Collar.Use(devtoolAddon)
	//
	ns := Collar.NS("com.collargo.test", map[string]string{})

	input := ns.Input("input")

	input.Map("x2", func(s Signal) (Signal, error) {
		fmt.Println("step 1.1", "x2", s.Payload)
		v, _ := s.Get(AnonPayload)
		newS := s.New(v.(int) * 2)
		fmt.Println("step 1.2", newS.Payload)
		return newS, nil
	}).Map("+1", func(s Signal) (Signal, error) {
		fmt.Println("step 1.3", "+1", s.Payload)
		v, _ := s.Get(AnonPayload)
		newS := s.New(v.(int) + 1)
		fmt.Println("step 1.4", newS.Payload)
		return newS, nil
	}).Do("test", func(s Signal) (interface{}, error) {
		v, _ := s.Get(AnonPayload)
		fmt.Println("step 1.5", v)
		assert.Equal(t, 21, v.(int))
		return "", nil
	}).Output("output")

	input.Push(10)
	time.Sleep(testDelay * time.Millisecond)
}

func TestNodeConnection(t *testing.T) {
	ns := Collar.NS("com.collargo.test", map[string]string{})

	sensor := ns.Sensor("test sensor", func(options string, send SendDataFunc) {
		time.Sleep(100 * time.Millisecond)
		send(11)
		time.Sleep(100 * time.Millisecond)
		send(10)
	}, false)

	errGen := ns.Do("error generator", func(s Signal) (interface{}, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		if v.Value%2 != 0 {
			return nil, errors.New("error")
		} else {
			return nil, nil
		}
	})

	filter := ns.Filter("even", func(s Signal) (bool, error) {
		v := new(IntPayload)
		s.GetValue("__anon__", v)
		return v.Value%2 == 0, nil
	})

	double := ns.Map("x2", func(s Signal) (Signal, error) {
		fmt.Println("step 1", "x2", s.Payload)
		v, _ := s.Get(AnonPayload)
		newS := s.New(v.(int) * 2)
		fmt.Println("step 2", newS.Payload)
		return newS, nil
	})

	inc := ns.Map("+1", func(s Signal) (Signal, error) {
		fmt.Println("step 3", "+1", s.Payload)
		v, _ := s.Get(AnonPayload)
		newS := s.New(v.(int) + 1)
		fmt.Println("step 4", newS.Payload)
		return newS, nil
	})

	greeting := ns.Do("greeting", func(s Signal) (interface{}, error) {
		fmt.Println("step 5", "greeting", s.Payload)
		v, _ := s.Get(AnonPayload)
		fmt.Println("step 6")
		return "Hello! " + strconv.Itoa(v.(int)), nil
	})

	testGreeting := ns.Do("test", func(s Signal) (interface{}, error) {
		v, _ := s.Get(AnonPayload)
		greeting, _ := s.GetResult()
		if v.(int) == 11 {
			return "", nil
		}
		fmt.Println("step 7", s.Payload)
		assert.Equal(t, 21, v.(int))
		assert.Equal(t, greeting, "Hello! 21")
		return "", nil
	})

	errorHandling := ns.Errors("error handler", func(s Signal, rethrow SendSignalFunc) error {
		rethrow(s)
		return nil
	})

	sensor.
		To("error gen", errGen).
		To("even number", filter).
		To("x2", double).
		To("+1", inc).
		To("greeting", greeting).
		To("error", errorHandling).
		To("test", testGreeting)

	time.Sleep(testDelay * time.Millisecond)
}
