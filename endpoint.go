package collargo

/**
 * Signal Processor for endpoint operator (input/output)
 */

type endpointProcessor struct {
}

func (processor endpointProcessor) OnError(s Signal, send SendSignalFunc) error {
	send(s)
	return nil
}

func (processor endpointProcessor) OnSignal(s Signal, send SendSignalFunc) error {
	send(s)
	return nil
}

// Input input endpoint operator
type Input struct {
	Node
}

// Output output endpoint operator
type Output struct {
	Node
}
