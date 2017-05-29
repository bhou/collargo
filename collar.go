package collargo

// CollarType the top level collar type
type collarType struct {
	Namespace
	// Observers the global observers for all nodes
	observers []Observer
	// Executor the executor
	executor Executor
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
