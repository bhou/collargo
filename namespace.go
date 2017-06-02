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
	// Filter a filter node
	Filter(comment string, filter FilterCallback) Filter
	// When an alias of Filter
	When(comment string, filter FilterCallback) Filter

	// Create a processor node
	Processor(comment string, process ProcessCallback) Processor
	// Alias of Processor
	Map(comment string, process ProcessCallback) Processor

	// Create an actuator node
	Actuator(comment string, act ActCallback) Actuator
	// Alias of Actuator
	Do(comment string, act ActCallback) Actuator

	// Create an error handling node
	Errors(comment string, errHandler ErrorCallback) ErrorNode
	// Create an input endpoint operator
	Input(comment string) Input
	// Create an output endpoint operator
	Output(comment string) Output
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

// Filter create a filter operator
func (ns *namespaceType) Filter(comment string, accept FilterCallback) Filter {
	node := CreateNode(comment, ns.GetNamespace(), filterProcessor{
		accept: accept,
	})

	for k, v := range ns.GetMetadata() {
		node.AddMeta(k, v)
	}
	node.SetType("filter")

	filterOp := Filter{
		Node: node,
	}

	return filterOp
}

// When an alias of Filter
func (ns *namespaceType) When(comment string, filter FilterCallback) Filter {
	return ns.Filter(comment, filter)
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

// Map alias of Processor
func (ns *namespaceType) Map(comment string, process ProcessCallback) Processor {
	return ns.Processor(comment, process)
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

// Map alias of Actuator
func (ns *namespaceType) Do(comment string, act ActCallback) Actuator {
	return ns.Actuator(comment, act)
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

// Input create an input operator
func (ns *namespaceType) Input(comment string) Input {
	node := CreateNode(comment, ns.GetNamespace(), endpointProcessor{})

	for k, v := range ns.GetMetadata() {
		node.AddMeta(k, v)
	}
	node.SetType("endpoint.input")

	input := Input{
		Node: node,
	}

	return input
}

// Output create an output operator
func (ns *namespaceType) Output(comment string) Output {
	node := CreateNode(comment, ns.GetNamespace(), endpointProcessor{})

	for k, v := range ns.GetMetadata() {
		node.AddMeta(k, v)
	}
	node.SetType("endpoint.out")

	output := Output{
		Node: node,
	}

	return output
}
