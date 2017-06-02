package collargo

// Executable the executable function type
type Executable func(s Signal, send SendSignalFunc) error

// Executor the executor interface
type Executor interface {
	// Schedule an executable with a given signal and send function
	Schedule(executable Executable, node Node, s Signal)
	// start execute
	Execute()
}

type defaultExecutor struct {
}

func (executor defaultExecutor) Schedule(executable Executable, node Node, s Signal) {
	send := func(signal Signal) {
		node.Send(signal)
	}

	run := func() {
		err := executable(s, send)

		if err != nil {
			errSignal := s.SetError(err)
			send(errSignal)
		}
	}
	go run()
}

func (executor defaultExecutor) Execute() {
	// do nothing
}
