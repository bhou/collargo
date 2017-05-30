package collargo

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

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

	time.Sleep(3000 * time.Millisecond)
}

func TestSensorWithDeferWatch(t *testing.T) {
	sensor := Collar.Sensor("test sensor", func(options string, send SendDataFunc) {
		time.Sleep(1000 * time.Millisecond)
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

	time.Sleep(3000 * time.Millisecond)
}

func TestProcessor(t *testing.T) {
	sensor := Collar.Sensor("test sensor", func(options string, send SendDataFunc) {
		time.Sleep(1000 * time.Millisecond)
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

	time.Sleep(3000 * time.Millisecond)
}

func TestActuator(t *testing.T) {
	sensor := Collar.Sensor("test sensor", func(options string, send SendDataFunc) {
		time.Sleep(1000 * time.Millisecond)
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

	time.Sleep(3000 * time.Millisecond)
}

func TestErrors(t *testing.T) {
	sensor := Collar.Sensor("test sensor", func(options string, send SendDataFunc) {
		time.Sleep(1000 * time.Millisecond)
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

	time.Sleep(3000 * time.Millisecond)
}

func TestMultipleFlow(t *testing.T) {
	sensor := Collar.Sensor("test sensor", func(options string, send SendDataFunc) {
		time.Sleep(1000 * time.Millisecond)
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

	time.Sleep(3000 * time.Millisecond)
}

func TestEndpoint(t *testing.T) {
	devtoolAddon := CreateDevToolAddon("ws://localhost:7500/app")

	Collar.Use(devtoolAddon)

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
	time.Sleep(3000 * time.Millisecond)
}
