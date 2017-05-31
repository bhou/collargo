package collargo

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type Message struct {
	Type string      `json:"type"`
	ID   string      `json:"id"`
	Data interface{} `json:"data"`
}

type MessageHandler func(interface{}) error

// WebsocketClient the client to connect to websocket
type WebsocketClient struct {
	url          string
	clientID     string
	clientSecret string
	conn         *websocket.Conn
	handlers     map[string]([]MessageHandler)
}

// CreateWebsocketClient Create a websocket client
func CreateWebsocketClient(url string, clientID string, clientSecret string) WebsocketClient {
	return WebsocketClient{
		url:          url,
		clientID:     clientID,
		clientSecret: clientSecret,
		handlers:     map[string]([]MessageHandler){},
	}
}

// Connect  connect to the server
func (client *WebsocketClient) Connect() error {
	var err error
	conn, _, err := websocket.DefaultDialer.Dial(client.url, nil)
	if err != nil {
		return err
	}

	client.conn = conn

	go func() {
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("receive err", err.Error())
				return
			}

			log.Println("receive message", string(message))

			var receivedMsg Message

			err = json.Unmarshal(message, &receivedMsg)

			log.Println("received message", receivedMsg)
			if err != nil {
				log.Println("Failed to unmarshal received message", err)
				continue
			}

			// handle received message
			switch receivedMsg.Type {
			case "authorized":
				log.Println("authorized!")
			case "unauthorized":
				log.Println("unauthorized!")
				return
			default:
				if handlers, ok := client.handlers[receivedMsg.Type]; ok {
					for i := range handlers {
						err := handlers[i](receivedMsg.Data)
						if err != nil {
							log.Println("Websocket message handler: ", err)
						}
					}
				}
			}
		}
	}()

	err = client.Emit("authentication", map[string]string{
		"clientId":     client.clientID,
		"clientSecret": client.clientSecret,
	})
	return err
}

// On add a message handler
func (client *WebsocketClient) On(msg string, handler MessageHandler) {
	_, ok := client.handlers[msg]
	if !ok {
		client.handlers[msg] = []MessageHandler{}
	}
	handlers := client.handlers[msg]
	client.handlers[msg] = append(handlers, handler)
}

// Emit emit a message to server
func (client *WebsocketClient) Emit(msg string, data interface{}) error {
	var err error
	m := Message{
		Type: msg,
		ID:   client.clientID,
		Data: data,
	}

	byteMsg, err := json.Marshal(m)
	if err != nil {
		log.Println("Emit:", "error when marshaling dev tool message", err)
		return err
	}

	// marshaled, _ := json.MarshalIndent(m, "", " ")
	// log.Println(string(marshaled))

	err = client.conn.WriteMessage(websocket.TextMessage, byteMsg)
	if err != nil {
		log.Println("Emit:", "error when sending message to collar dev server", err)
		return err
	}

	return nil
}
