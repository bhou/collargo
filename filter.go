package collargo

/**
 * Filter Operator callback
 */

// FilterCallback the callback function for filter operator
type FilterCallback func(s Signal) (bool, error)

/**
 * Signal Processor for Do (Actuator) operator
 */

type filterProcessor struct {
	accept FilterCallback
}

func (filter filterProcessor) OnError(s Signal, send SendSignalFunc) error {
	send(s)
	return nil
}

func (filter filterProcessor) OnSignal(s Signal, send SendSignalFunc) error {
	pass, err := filter.accept(s)

	if err != nil {
		return err
	}

	if pass {
		send(s)
	}

	return nil
}

// Filter the filter operator type
type Filter struct {
	Node
}
