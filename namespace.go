package collargo

// import (
//  "fmt"
// )

// Namespace the namespace
type Namespace interface {
	GetNamespace() string
	GetMetadata() map[string]string

	/* static operators */

	// Create a sensor node
	Sensor(comment string, watch SensorCallback, deferWatch bool) Sensor
	// Create a processor node
	Processor(comment string, process ProcessCallback) Processor
	// Create an actuator node
	Actuator(comment string, act ActCallback) Actuator
	// Create an error handling node
	Errors(comment string, errHandler ErrorCallback) ErrorNode
}

type namespaceType struct {
	namespace string
	metadata  map[string]string
}

func (ns *namespaceType) GetNamespace() string {
	return ns.namespace
}

func (ns *namespaceType) GetMetadata() map[string]string {
	return ns.metadata
}

/*
  Namespace Operators
*/
// Sensor create a sensor operator
func (ns *namespaceType) Sensor(comment string, watch SensorCallback, deferWatch bool) Sensor {
	node := CreateNode(comment, ns.GetNamespace(), sensorProcessor{
		watch: watch,
	})

	for k, v := range ns.GetMetadata() {
		node.AddMeta(k, v)
	}
	node.SetType("sensor")

	sensor := Sensor{
		Node: node,
	}

	if !deferWatch {
		sensor.Watch("initiated")
	}

	return sensor
}

// Processor create a processor operator
func (ns *namespaceType) Processor(comment string, process ProcessCallback) Processor {
	node := CreateNode(comment, ns.GetNamespace(), mapProcessor{
		process: process,
	})

	for k, v := range ns.GetMetadata() {
		node.AddMeta(k, v)
	}
	node.SetType("processor")

	processor := Processor{
		Node: node,
	}

	return processor
}

// Actuator create an actuator operator
func (ns *namespaceType) Actuator(comment string, act ActCallback) Actuator {
	node := CreateNode(comment, ns.GetNamespace(), actProcessor{
		act: act,
	})

	for k, v := range ns.GetMetadata() {
		node.AddMeta(k, v)
	}
	node.SetType("actuator")

	actuator := Actuator{
		Node: node,
	}

	return actuator
}

// Errors create an error handler operator
func (ns *namespaceType) Errors(comment string, errorHandler ErrorCallback) ErrorNode {
	node := CreateNode(comment, ns.GetNamespace(), errorProcessor{
		errorHandler: errorHandler,
	})

	for k, v := range ns.GetMetadata() {
		node.AddMeta(k, v)
	}
	node.SetType("errorhandler")

	errors := ErrorNode{
		Node: node,
	}

	return errors
}
