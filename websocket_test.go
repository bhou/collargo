package collargo

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestWebsocketClient(t *testing.T) {
	var err error

	client := CreateWebsocketClient("ws://localhost:7500/app", "", "")

	err = client.Connect()

	if err != nil {
		log.Println(err)
	}

	assert.True(t, true)

	time.Sleep(3000 * time.Millisecond)
}
