package collargo

import (
	"github.com/satori/go.uuid"
	"log"
	"reflect"
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

type elementType map[string]interface{}

type signalType map[string]interface{}

func getStructType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

/**
 * Observer
 */

func handleNode(node Node) map[string]interface{} {
	ret := map[string]interface{}{}

	ret["group"] = "nodes"
	ret["data"] = map[string]interface{}{}
	ret["style"] = map[string]string{}
	ret["classses"] = ""

	data := ret["data"].(map[string]interface{})
	data["id"] = node.ID()
	data["model"] = node.Type() // getStructType(node)
	data["fullName"] = node.FullName()
	data["label"] = node.Comment()
	data["inputs"] = map[string]string{}
	data["outputs"] = map[string]string{}
	data["stack"] = map[string]string{}
	data["meta"] = node.GetAllMeta()
	data["tags"] = node.Tags()

	return ret
}

func handleEdge(upstream Node, downstream Node) map[string]interface{} {
	ret := map[string]interface{}{}

	ret["group"] = "edges"
	ret["data"] = map[string]interface{}{}

	data := ret["data"].(map[string]interface{})
	data["id"] = uuid.NewV4().String()
	data["source"] = upstream.ID()
	data["target"] = downstream.ID()
	return ret
}

/**
 * Addon
 */

// DevToolAddon the devtool addon
type DevToolAddon struct {
	observers []Observer
	elements  []elementType
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
				addon.pushBufferedElements()
				addon.pushBufferedSignals()

			case c := <-addon.command:
				switch c {
				case "quit":
					break
				}
			}
		}

		// for range ticker.C {
		//  addon.pushBufferedElements()
		//  addon.pushBufferedSignals()
		// }
	}()
}

func (addon *DevToolAddon) staticTopologyObserver(node Node, when string, s Signal, data ...interface{}) error {
	if when != "to" {
		return nil
	}

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

	return nil
}

func (addon *DevToolAddon) signalFlowObserver(node Node, when string, s Signal, data ...interface{}) error {
	if when != "onReceive" && when != "send" {
		return nil
	}

	signalElem := map[string]interface{}{
		"when":    when,
		"time":    time.Now().UnixNano() / int64(time.Millisecond),
		"nodeId":  node.ID(),
		"seq":     s.ID,
		"payload": s.Payload,
		"error":   nil,
		"end":     s.End,
	}

	addon.signals = append(addon.signals, signalElem)

	return nil
}

func (addon *DevToolAddon) pushBufferedElements() {
	if len(addon.elements) <= 0 {
		return
	}

	elemToBeSent := []elementType{}

	for i := range addon.elements {
		elemToBeSent = append(elemToBeSent, addon.elements[i])
	}

	addon.elements = []elementType{}

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
		elements:  []elementType{},
		signals:   []signalType{},
		client:    &client,
		nodes:     map[string]Node{},
		command:   make(chan string),
	}
	addon.observers = append(addon.observers, addon.staticTopologyObserver)
	addon.observers = append(addon.observers, addon.signalFlowObserver)

	client.On("push", func(data interface{}) error {
		log.Println("dev addon push", data)
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
			log.Println("Failed to push data: couldn't  don't find node with id:", id)
		}

		return nil
	})

	client.On("send", func(data interface{}) error {
		log.Println("dev addon send", data)
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
