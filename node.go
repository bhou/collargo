package collargo

import (
	"github.com/satori/go.uuid"
	"regexp"
	"sync"
)

// Observer function: observe signal processing
type Observer func(Node, string, Signal, ...interface{}) error

// SignalProcessor the basic execution unit inside of node
//
// a processor must provide 2 processing functions to handle data signal and error signal
type SignalProcessor interface {
	// handle the signal when signal represents an error
	OnError(s Signal, send SendSignalFunc) error
	// handle the signal when signal represents an end signal
	// OnEnd(s Signal, send func(Signal)) error
	// handle the signal
	OnSignal(s Signal, send SendSignalFunc) error
}

// Node the node interface
type Node interface {
	ID() string        // Get the id of the node
	Seq() string       // Get the sequence of the node (same as id)
	Name() string      // Get the name of the node, unique name in namespace
	Namespace() string // Get the namespace of the node
	FullName() string  // Get the full name of the name: namespace + name
	Comment() string   // Get the node comment
	Type() string      // Get the node type

	SetType(string) // Set the node type

	Upstreams() map[string]Node   // Get the upstreams
	Downstreams() map[string]Node // Get the downstreams

	SignalProcessor() SignalProcessor // Get the signal processor

	Tags() []string         // Get all tags
	HasTag(tag string) bool // Check if the node has tag or not

	AddMeta(name string, metadata string) Node // Add a metadata to the node
	GetAllMeta() map[string]string             // Get all metadatas
	GetMeta(name string) (string, bool)        // Get a metadata with name

	Push(data interface{}) Node // Push data to the node
	Send(data interface{}) Node // Send data to the downstream nodes

	To(comment string, next Node) Node // Connect the current node to a downstream node

	Observe(observer Observer) // Add an observer
	Observers() []Observer     // Get All observers of this node

	/* Flow related API */
	GetFlowOutputObserver() (Observer, bool)         // get flow output observer
	SetFlowOutputObserver(Observer)                  // set flow output observer
	AddFlowFunc(outID string, flowFunc FlowFunc)     // attach a flow function to the node
	GetFlowFunc(outID string) (FlowFunc, bool)       // get a flow function of the node
	AddSignalCallback(sigID string, cb Callback)     // add signal processing callback
	GetSignalCallback(sigID string) (Callback, bool) // get signal processing callback
	DelSignalCallback(sigID string)                  // delete signal processing callback

	// operators
	Do(comment string, act ActCallback) Actuator
	When(comment string, accept FilterCallback) Filter
	Map(comment string, process ProcessCallback) Processor
	Errors(comment string, errHandler ErrorCallback) ErrorNode
	Input(comment string) Input
	Output(comment string) Output
}

// parseNameFromComment   In the node comment you can put a unique (unique in namespace) name with @ sign
func parseNameFromComment(comment string) string {
	re := regexp.MustCompile(`@(\w|-)+`)
	found := re.FindString(comment)
	if found == "" {
		return found
	}
	return found[1:]
}

// parseTagsFromComment   You can also add tags in comment, each tag starts with # sign
func parseTagsFromComment(comment string) []string {
	re := regexp.MustCompile(`#(\w|-)+`)
	tags := re.FindAllString(comment, -1)
	if tags == nil {
		return []string{}
	}

	var ret []string
	for _, tag := range tags {
		if tag != "#" {
			ret = append(ret, tag[1:])
		}
	}

	return ret
}

// parseInfoFromComment   get the node name, tags and real comment from the comment text
func parseInfoFromComment(comment string) (name string, tags []string, newComment string) {
	name = parseNameFromComment(comment)
	tags = parseTagsFromComment(comment)

	re := regexp.MustCompile(`(@(\w|-)+|#(\w|-)+)`)
	whitePrefix := regexp.MustCompile(`^(\s)+`)

	newComment = re.ReplaceAllLiteralString(comment, "")

	newComment = whitePrefix.ReplaceAllLiteralString(newComment, "")

	return name, tags, newComment
}

// node the node struct
type node struct {
	sync.RWMutex
	id        string
	seq       string
	comment   string
	name      string
	namespace string
	nodeType  string

	upstreams   map[string]Node
	downstreams map[string]Node

	tags []string
	meta map[string]string

	observers []Observer
	processor SignalProcessor

	// property used for flow function
	flowOutputObserver Observer
	flowFuncs          map[string]FlowFunc
	signalCallbacks    map[string]Callback
}

// CreateNode create a node
func CreateNode(
	comment string, // comment text of the node
	namespace string, // namespace of the node
	processor SignalProcessor, // processor used to handle signals
) Node {
	n := new(node)
	name, tags, commentText := parseInfoFromComment(comment)

	n.name = name
	n.comment = commentText
	n.tags = tags

	n.id = uuid.NewV1().String()
	n.seq = n.id
	n.namespace = namespace
	n.observers = []Observer{}
	n.upstreams = map[string]Node{}
	n.downstreams = map[string]Node{}
	n.processor = processor
	n.meta = map[string]string{
		"namespace": namespace,
	}

	n.flowOutputObserver = nil
	n.flowFuncs = map[string]FlowFunc{}
	n.signalCallbacks = map[string]Callback{}

	return n
}

// ID Get node id
func (n *node) ID() string {
	return n.id
}

// Seq get the node unique id
func (n *node) Seq() string {
	return n.id
}

// Name Get node name
func (n *node) Name() string {
	if n.name == "" {
		return n.id
	}
	return n.name
}

// Type Get node type
func (n *node) Type() string {
	if n.nodeType == "" {
		return "node"
	}
	return n.nodeType
}

// SetType set the node type
func (n *node) SetType(t string) {
	n.nodeType = t
}

// Namespace Get node namespace
func (n *node) Namespace() string {
	return n.namespace
}

// FullName Get node full name
func (n *node) FullName() string {
	return n.namespace + "." + n.Name()
}

// Comment Get the comment of the node
func (n *node) Comment() string {
	return n.comment
}

// Upstreams Get the upstreams of this node
func (n *node) Upstreams() map[string]Node {
	return n.upstreams
}

// Downstreams Get the downstreams of this node
func (n *node) Downstreams() map[string]Node {
	return n.downstreams
}

// SignalProcessor get the signal processor of this node
func (n *node) SignalProcessor() SignalProcessor {
	return n.processor
}

// Tags Get node tags
func (n *node) Tags() []string {
	return n.tags
}

// HasTag Check if node has tag
func (n *node) HasTag(name string) bool {
	for _, v := range n.tags {
		if v == name {
			return true
		}
	}
	return false
}

// AddMeta add the meta data to node
func (n *node) AddMeta(name string, data string) Node {
	n.meta[name] = data
	return n
}

// GetAllMeta get all metadatas attatched to this node
func (n *node) GetAllMeta() map[string]string {
	return n.meta
}

// GetMeta get the metadata according to name
func (n *node) GetMeta(name string) (string, bool) {
	v, ok := n.meta[name]
	return v, ok
}

// Push push a signal to the node for processing
func (n *node) Push(data interface{}) Node {
	s := CreateSignal(data)
	// fmt.Println("push", s.Payload)
	n.onReceive(s)
	return n
}

// handle received signal
func (n *node) onReceive(s Signal) Node {
	err := n.invokeOnReceiveObservers(s)

	if err != nil {
		panic(err)
	}

	// fmt.Println("onReceive", s.Payload)
	if s.Error != nil {
		Collar.GetExecutor().Schedule(n.processor.OnError, n, s)
	} else {
		Collar.GetExecutor().Schedule(n.processor.OnSignal, n, s)
	}

	return n
}

// Send send a signal to the downstream nodes
func (n *node) Send(data interface{}) Node {
	s := CreateSignal(data)

	// fmt.Println("send signal", s.Payload)

	err := n.invokeSendObservers(s)

	if err != nil {
		panic(err)
	}

	// Each downstream node handles the signal in a goroutine
	for _, stream := range n.downstreams {
		go stream.Push(s)
	}

	return n
}

// To Connect the current node To the next node
func (n *node) To(comment string, next Node) Node {
	err := n.invokeToObservers(next)

	if err != nil {
		panic(err)
	}

	n.Downstreams()[next.ID()] = next
	next.Upstreams()[n.ID()] = n

	return next
}

// Observe observe the node with an observer
func (n *node) Observe(observer Observer) {
	n.observers = append(n.observers, observer)
}

// Observers get all observers of this node
func (n *node) Observers() []Observer {
	return n.observers
}

func (n *node) GetFlowOutputObserver() (Observer, bool) {
	// n.RLock()
	if n.flowOutputObserver == nil {
		return nil, false
	}
	observer := n.flowOutputObserver
	// n.RUnlock()
	return observer, true
}

func (n *node) SetFlowOutputObserver(observer Observer) {
	n.Lock()
	n.flowOutputObserver = observer
	n.Unlock()
}

func (n *node) AddFlowFunc(outID string, flowFunc FlowFunc) {
	n.Lock()
	n.flowFuncs[outID] = flowFunc
	n.Unlock()
}

func (n *node) GetFlowFunc(outID string) (FlowFunc, bool) {
	// n.RLock()
	flowFunc, existed := n.flowFuncs[outID]
	// n.RUnlock()
	return flowFunc, existed
}

func (n *node) AddSignalCallback(sigID string, cb Callback) {
	n.Lock()
	n.signalCallbacks[sigID] = cb
	n.Unlock()
}

func (n *node) GetSignalCallback(sigID string) (Callback, bool) {
	n.RLock()
	cb, existed := n.signalCallbacks[sigID]
	n.RUnlock()
	return cb, existed
}

func (n *node) DelSignalCallback(sigID string) {
	_, existed := n.GetSignalCallback(sigID)
	if existed {
		n.Lock()
		delete(n.signalCallbacks, sigID)
		n.Unlock()
	}
}

/*
 Operators
*/

func (n *node) When(comment string, accept FilterCallback) Filter {
	filterNode := CreateNode(comment, n.Namespace(), filterProcessor{
		accept: accept,
	})

	filterNode.SetType("filter")

	filter := Filter{
		Node: filterNode,
	}

	n.To(comment, filter)

	return filter
}

func (n *node) Map(comment string, process ProcessCallback) Processor {
	mapNode := CreateNode(comment, n.Namespace(), mapProcessor{
		process: process,
	})

	mapNode.SetType("processor")

	processor := Processor{
		Node: mapNode,
	}

	n.To(comment, processor)

	return processor
}

func (n *node) Do(comment string, act ActCallback) Actuator {
	actNode := CreateNode(comment, n.Namespace(), actProcessor{
		act: act,
	})

	actNode.SetType("actuator")

	actuator := Actuator{
		Node: actNode,
	}

	n.To(comment, actuator)

	return actuator
}

func (n *node) Errors(comment string, errorHandler ErrorCallback) ErrorNode {
	errNode := CreateNode(comment, n.Namespace(), errorProcessor{
		errorHandler: errorHandler,
	})

	errNode.SetType("errorhandler")

	errors := ErrorNode{
		Node: errNode,
	}

	n.To(comment, errors)

	return errors
}

func (n *node) Input(comment string) Input {
	inputNode := CreateNode(comment, n.Namespace(), endpointProcessor{})

	inputNode.SetType("endpoint.input")

	input := Input{
		Node: inputNode,
	}

	n.To(comment, input)

	return input
}

func (n *node) Output(comment string) Output {
	outputNode := CreateNode(comment, n.Namespace(), endpointProcessor{})

	outputNode.SetType("endpoint.output")

	output := Output{
		Node: outputNode,
	}

	n.To(comment, output)

	return output
}

/*
 private
*/

// invoke Global observers
func (n *node) invokeGlobalObservers(when string, signal Signal, data ...interface{}) error {
	var err error
	for _, observer := range Collar.observers {
		err = observer(n, when, signal, data...)
		if err != nil {
			// fmt.Println("global observers error", err)
			return err
		}
	}
	return nil
}

// invoke OnReceive observers
func (n *node) invokeOnReceiveObservers(signal Signal) error {
	// fmt.Println("invoke onReceive", signal.Payload)
	err := n.invokeGlobalObservers("onReceive", signal)
	if err != nil {
		return err
	}

	for _, observer := range n.observers {
		err = observer(n, "onReceive", signal)
		if err != nil {
			return err
		}
	}

	return nil
}

// invoke Send observers
func (n *node) invokeSendObservers(signal Signal) error {
	err := n.invokeGlobalObservers("send", signal)
	if err != nil {
		return err
	}

	for _, observer := range n.observers {
		err = observer(n, "send", signal)
		if err != nil {
			return err
		}
	}
	return nil
}

// invoke To observers
func (n *node) invokeToObservers(downstream Node) error {
	err := n.invokeGlobalObservers("to", Signal{}, downstream)
	if err != nil {
		return err
	}

	for _, observer := range n.observers {
		err = observer(n, "to", Signal{}, downstream)
		if err != nil {
			return err
		}
	}
	return nil
}
