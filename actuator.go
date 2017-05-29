package collargo

/**
 * Actuator Operator callback
 */

// ActCallback the callback function for actuator operator
type ActCallback func(s Signal) (interface{}, error)

/**
 * Signal Processor for Do (Actuator) operator
 */

type actProcessor struct {
	act ActCallback
}

func (actuator actProcessor) OnError(s Signal, send SendSignalFunc) error {
	send(s)
	return nil
}

func (actuator actProcessor) OnSignal(s Signal, send SendSignalFunc) error {
	result, err := actuator.act(s)

	if err != nil {
		return err
	}

	newSignal := s.SetResult(result)
	send(newSignal)
	return nil
}

// Actuator the actuator operator type
type Actuator struct {
	Node
}
