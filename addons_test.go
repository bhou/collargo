package collargo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDevToolAddon(t *testing.T) {
	devtoolAddon := CreateDevToolAddon("ws://localhost:7500/app")

	ns := Collar.NS("com.collargo.test", map[string]string{
		"module": "test",
	})

	Collar.Use(devtoolAddon)

	sensor := ns.Sensor("test sensor", func(options string, send SendDataFunc) {
		time.Sleep(1000 * time.Millisecond)
		send(float64(10))
		time.Sleep(2000 * time.Millisecond)
		send(float64(30))
	}, false)

	sensor.Map("@double x2", func(s Signal) (Signal, error) {
		fmt.Println("step 1.1", "x2", s.Payload)
		v, _ := s.Get(AnonPayload)
		newS := s.New(v.(float64) * 2)
		fmt.Println("step 1.2", newS.Payload)
		return newS, nil
	}).Map("+1", func(s Signal) (Signal, error) {
		fmt.Println("step 1.3", "+1", s.Payload)
		v, _ := s.Get(AnonPayload)
		newS := s.New(v.(float64) + 1)
		fmt.Println("step 1.4", newS.Payload)
		return newS, nil
	}).Do("test", func(s Signal) (interface{}, error) {
		v, _ := s.Get(AnonPayload)
		fmt.Println("step 1.5", v)
		assert.Equal(t, 21, v.(float64))
		return "", nil
	})

	sensor.Map("x3", func(s Signal) (Signal, error) {
		fmt.Println("step 2.1", "x3", s.Payload)
		v, _ := s.Get(AnonPayload)
		newS := s.New(v.(float64) * 3)
		fmt.Println("step 2.2", newS.Payload)
		return newS, nil
	}).Map("+1", func(s Signal) (Signal, error) {
		fmt.Println("step 2.3", "+1", s.Payload)
		v, _ := s.Get(AnonPayload)
		newS := s.New(v.(float64) + 1)
		fmt.Println("step 2.4", newS.Payload)
		return newS, nil
	}).Do("test", func(s Signal) (interface{}, error) {
		v, _ := s.Get(AnonPayload)
		fmt.Println("step 2.5", v)
		assert.Equal(t, 31, v.(float64))
		return "", nil
	})

	assert.True(t, true)

	time.Sleep(100000 * time.Millisecond)
}
