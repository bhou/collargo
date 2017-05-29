package collargo

/**
 * Errors Operator callback
 */

// ErrorCallback the callback function for Errors operator
type ErrorCallback func(s Signal, rethrow SendSignalFunc) error

/**
 * Signal Processor for error
 */

type errorProcessor struct {
	errorHandler ErrorCallback
}

func (ep errorProcessor) OnError(s Signal, send SendSignalFunc) error {
	err := ep.errorHandler(s, send)
	return err
}

func (ep errorProcessor) OnSignal(s Signal, send SendSignalFunc) error {
	send(s)
	return nil
}

// ErrorNode the error handler operator
type ErrorNode struct {
	Node
}
