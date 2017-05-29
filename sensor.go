package collargo

/**
 * Sensor operator callback
 */

// SensorCallback the callback function for sensor operator
type SensorCallback func(options string, send SendDataFunc)

/**
 * Signal Processor for sensor
 */

type sensorProcessor struct {
	watch SensorCallback
}

func (processor sensorProcessor) OnError(s Signal, send SendSignalFunc) error {
	send(s)
	return nil
}

func (processor sensorProcessor) OnSignal(s Signal, send SendSignalFunc) error {
	// sensor will not process incoming signal
	return nil
}

/**
 * Sensor Node
 */

// Sensor the sensor operator
type Sensor struct {
	Node
}

// Watch start to watch the external world
func (sensor *Sensor) Watch(options string) {
	go sensor.SignalProcessor().(sensorProcessor).watch(options, func(data interface{}) {
		sensor.Send(data)
	})
}
