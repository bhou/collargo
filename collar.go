package collargo

import (
// "log"
)

// CollarType the top level collar type
type collarType struct {
	Namespace
	// Observers the global observers for all nodes
	observers []Observer
	// Executor the executor
	executor Executor
}

type callbackResult struct {
	err    error
	result interface{}
}

// SetExecutor set the executor
func (collar *collarType) SetExecutor(executor Executor) {
	collar.executor = executor
}

// GetExecutor get the executor
func (collar collarType) GetExecutor() Executor {
	return collar.executor
}

func (collar collarType) NS(ns string, meta map[string]string) Namespace {
	return &namespaceType{
		namespace: ns,
		metadata:  meta,
	}
}

func (collar *collarType) Use(addon Addon) {
	obs := addon.Observers()

	for i := range obs {
		collar.observers = append(collar.observers, obs[i])
	}

	addon.Run()
}

func (collar *collarType) ToFlowFunc(input Node, output Node) FlowFunc {
	_, existed := output.GetFlowOutputObserver()
	if !existed {
		observer := func(node Node, when string, signal Signal, data ...interface{}) error {
			if when != "send" {
				return nil
			}

			destTag, ok := signal.GetTag("__to_node_dest__")
			if !ok || destTag != output.ID() {
				return nil
			}

			cb, existed := output.GetSignalCallback(signal.ID)

			if !existed {
				return nil
			}

			output.DelSignalCallback(signal.ID)

			if signal.Error != nil {
				cb(signal.Error, nil)
			} else {
				cb(nil, signal.Payload)
			}

			return nil
		}

		output.SetFlowOutputObserver(observer)

		output.Observe(observer)
	}

	flowFunc, existed := input.GetFlowFunc(output.ID())

	if existed {
		return flowFunc
	}

	flowFunc = func(data interface{}) (interface{}, error) {
		signal := CreateSignal(data)
		signal = signal.SetTag("__to_node_dest__", output.ID())

		ch := make(chan callbackResult)
		output.AddSignalCallback(signal.ID, func(err error, result interface{}) {
			if err != nil {
				ch <- callbackResult{
					err: err,
				}
				return
			}
			ch <- callbackResult{
				err:    nil,
				result: result,
			}
		})

		input.Push(signal)

		var result callbackResult
		result = <-ch
		return result.result, result.err
	}

	input.AddFlowFunc(output.ID(), flowFunc)

	return flowFunc
}

var (
	defaultNS = namespaceType{
		namespace: "",
		metadata: map[string]string{
			"namespace": "",
		},
	}

	// Collar the top level collar object
	Collar = collarType{
		Namespace: &defaultNS,
		observers: []Observer{},
		executor:  defaultExecutor{},
	}
)
