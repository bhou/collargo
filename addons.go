package collargo

import (
	"github.com/satori/go.uuid"
	"log"
	// "reflect"
	"sync"
	"time"
)

// Addon interface
type Addon interface {
	Observers() []Observer // get an array of global observers

	Run() // run the addon

	Stop() // stop the addon
}

/**
 * dev addon
 */
type elemData struct {
	ID       string            `json:"id"`
	Model    string            `json:"model"`
	FullName string            `json:"fullName"`
	Label    string            `json:"label"`
	Inputs   map[string]string `json:"inputs"`
	Outputs  map[string]string `json:"outputs"`
	Stack    map[string]string `json:"stack"`
	Meta     map[string]string `json:"meta"`
	Tags     []string          `json:"tags,omitempty"`
	Source   string            `json:"source"`
	Target   string            `json:"target"`
}

type elemType struct {
	Group   string            `json:"group"`
	Data    elemData          `json:"data"`
	Style   map[string]string `json:"style"`
	Classes string            `json:"classes"`
}

type signalData struct {
	ID string `json:"id"`
}

type signalType struct {
	When    string                 `json:"when"`
	Time    int64                  `json:"time"`
	NodeId  string                 `json:"nodeId"`
	Seq     string                 `json:"seq"`
	Payload map[string]interface{} `json:"payload"`
	Error   error                  `json:"error"`
	End     bool                   `json:"end"`
}

// func getStructType(myvar interface{}) string {
//  if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
//    return "*" + t.Elem().Name()
//  } else {
//    return t.Name()
//  }
// }

/**
 * Observer
 */

func handleNode(node Node) elemType {
	return elemType{
		Group:   "nodes",
		Style:   map[string]string{},
		Classes: "",
		Data: elemData{
			ID:       node.ID(),
			Model:    node.Type(),
			FullName: node.FullName(),
			Label:    node.Comment(),
			Inputs:   map[string]string{},
			Outputs:  map[string]string{},
			Stack:    map[string]string{},
			Meta:     node.GetAllMeta(),
			Tags:     node.Tags(),
		},
	}
}

func handleEdge(upstream Node, downstream Node) elemType {
	return elemType{
		Group: "edges",
		Data: elemData{
			ID:     uuid.NewV4().String(),
			Source: upstream.ID(),
			Target: downstream.ID(),
		},
	}
}

/**
 * Addon
 */

// DevToolAddon the devtool addon
type DevToolAddon struct {
	sync.RWMutex
	observers []Observer
	elements  []elemType
	signals   []signalType

	nodes map[string]Node

	client *WebsocketClient

	command chan string
}

// Observers get the observers of the devtool addon
func (addon DevToolAddon) Observers() []Observer {
	return addon.observers
}

// Stop stop the addon
func (addon *DevToolAddon) Stop() {
	addon.command <- "quit"
}

// Run run the dev tool addon
func (addon *DevToolAddon) Run() {
	var err error

	err = addon.client.Connect()
	if err != nil {
		log.Println("Failed to connect to collar dev server", err)
	}
	addon.client.Emit("new model", map[string]string{
		"process": "__anonymous__",
	})

	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case _ = <-ticker.C:
				addon.Lock()
				addon.pushBufferedElements()
				addon.pushBufferedSignals()
				addon.Unlock()

			case c := <-addon.command:
				switch c {
				case "quit":
					break
				}
			}
		}
	}()
}

func (addon *DevToolAddon) staticTopologyObserver(node Node, when string, s Signal, data ...interface{}) error {
	if when != "to" {
		return nil
	}

	addon.Lock()
	downstream := data[0].(Node)

	// add nodes to elements list
	addon.elements = append(addon.elements, handleNode(node))
	addon.elements = append(addon.elements, handleNode(downstream))

	// mark nodes as processed
	if _, ok := addon.nodes[node.ID()]; !ok {
		addon.nodes[node.ID()] = node
	}
	if _, ok := addon.nodes[downstream.ID()]; !ok {
		addon.nodes[downstream.ID()] = downstream
	}

	// handle edges
	addon.elements = append(addon.elements, handleEdge(node, downstream))
	addon.Unlock()

	return nil
}

func (addon *DevToolAddon) signalFlowObserver(node Node, when string, s Signal, data ...interface{}) error {
	if when != "onReceive" && when != "send" {
		return nil
	}

	signalElem := signalType{
		When:    when,
		Time:    time.Now().UnixNano() / int64(time.Millisecond),
		NodeId:  node.ID(),
		Seq:     s.ID,
		Payload: s.Payload,
		Error:   nil,
		End:     s.End,
	}

	addon.Lock()
	addon.signals = append(addon.signals, signalElem)
	addon.Unlock()

	return nil
}

func (addon *DevToolAddon) pushBufferedElements() {
	if len(addon.elements) <= 0 {
		return
	}

	elemToBeSent := []elemType{}

	for i := range addon.elements {
		elemToBeSent = append(elemToBeSent, addon.elements[i])
	}
	addon.elements = []elemType{}

	addon.client.Emit("append elements", map[string]interface{}{
		"elements": elemToBeSent,
	})
}

func (addon *DevToolAddon) pushBufferedSignals() {
	if len(addon.signals) <= 0 {
		return
	}

	signalToBeSent := []signalType{}

	for i := range addon.signals {
		signalToBeSent = append(signalToBeSent, addon.signals[i])
	}

	addon.signals = []signalType{}

	addon.client.Emit("append signals", map[string]interface{}{
		"signals": signalToBeSent,
	})
}

// CreateDevToolAddon create a new development addon
func CreateDevToolAddon(url string) Addon {
	client := CreateWebsocketClient(url, "", "")

	addon := DevToolAddon{
		observers: []Observer{},
		elements:  []elemType{},
		signals:   []signalType{},
		client:    &client,
		nodes:     map[string]Node{},
		command:   make(chan string),
	}
	addon.observers = append(addon.observers, addon.staticTopologyObserver)
	addon.observers = append(addon.observers, addon.signalFlowObserver)

	client.On("push", func(data interface{}) error {
		mapData := data.(map[string]interface{})
		nodeId, ok := mapData["nodeId"]

		payload, _ := mapData["signal"].(map[string]interface{})["payload"]

		id := nodeId.(string)

		if !ok {
			log.Println("Failed to push data: data don't have nodeId property")
			return nil
		}

		if node, ok := addon.nodes[id]; ok {
			log.Println(node.Downstreams())
			node.Push(payload)
		} else {
			log.Println("Failed to push data: couldn't find node with id:", id)
		}

		return nil
	})

	client.On("send", func(data interface{}) error {
		mapData := data.(map[string]interface{})
		nodeId, ok := mapData["nodeId"]

		payload, _ := mapData["signal"].(map[string]interface{})["payload"]

		id := nodeId.(string)

		if !ok {
			log.Println("Failed to send data: data don't have nodeId property")
			return nil
		}

		if node, ok := addon.nodes[id]; ok {
			node.Send(payload)
		} else {
			log.Println("Failed to send data: couldn't  don't find node with id:", id)
		}

		return nil
	})

	return &addon
}
