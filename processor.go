package collargo

/**
 * Processor Operator callback
 */

// ProcessCallback the callback function for processor operator
type ProcessCallback func(s Signal) (Signal, error)

/**
 * Signal Processor for map operator
 */

type mapProcessor struct {
	process ProcessCallback
}

func (processor mapProcessor) OnError(s Signal, send SendSignalFunc) error {
	send(s)
	return nil
}

func (processor mapProcessor) OnSignal(s Signal, send SendSignalFunc) error {
	newSignal, err := processor.process(s)

	if err != nil {
		return err
	}

	send(newSignal)
	return nil
}

// Processor the processor operator
type Processor struct {
	Node
}
