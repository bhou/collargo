package collargo

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type passThroughSignalProcessor struct {
}

func (processor passThroughSignalProcessor) OnError(s Signal, send SendSignalFunc) error {
	send(s)
	return nil
}

// func (processor passThroughSignalProcessor) OnEnd(s Signal, send func(signal Signal)) error {
//  send(s)
//  return nil
// }

func (processor passThroughSignalProcessor) OnSignal(s Signal, send SendSignalFunc) error {
	send(s)
	return nil
}

func TestCreateNode(t *testing.T) {
	node1 := CreateNode("Example Node 1", "com.collartechs.test", passThroughSignalProcessor{})
	node2 := CreateNode("Example Node 2", "com.collartechs.test", passThroughSignalProcessor{})

	assert.Equal(t, "com.collartechs.test", node1.Namespace())
	assert.Equal(t, "com.collartechs.test", node2.Namespace())

	assert.Equal(t, "Example Node 1", node1.Comment())

	assert.NotNil(t, node1.ID())
	assert.Equal(t, node1.Name(), node1.ID())
	assert.Equal(t, "com.collartechs.test."+node1.Name(), node1.FullName())
	assert.NotEqual(t, node1.ID(), node2.ID())

	// test creating node with name and tags
	node3 := CreateNode("@uniqueName, #golang #collartechs", "com.collartechs.test", passThroughSignalProcessor{})

	assert.Equal(t, "uniqueName", node3.Name())
	assert.Equal(t, "golang", node3.Tags()[0])
	assert.Equal(t, "collartechs", node3.Tags()[1])
	assert.Equal(t, "com.collartechs.test.uniqueName", node3.FullName())
	assert.True(t, node3.HasTag("golang"))
	assert.True(t, node3.HasTag("collartechs"))
	assert.False(t, node3.HasTag("anotherTag"))

	assert.True(t, true)
}

func TestObserver(t *testing.T) {
	node := CreateNode("test node", "com.collartechs.test", passThroughSignalProcessor{})

	node.Observe(func(node Node, when string, signal Signal, data ...interface{}) error {
		if when == "onReceive" {
			str, existed := signal.Get(AnonPayload)
			assert.Equal(t, "text message", str.(string))
			assert.True(t, existed)
		}
		return nil
	})

	node.Push("text message")
}

func TestProcessingObserver(t *testing.T) {
	node := CreateNode("test node", "com.collartechs.test", passThroughSignalProcessor{})

	node.Observe(func(node Node, when string, signal Signal, data ...interface{}) error {
		if when == "onReceive" {
			err := signal.Error
			assert.NotNil(t, err)
			assert.Equal(t, "test error", err.Error())
		}
		return nil
	})

	node.Push(errors.New("test error"))
}

func TestChainedNode(t *testing.T) {
	node1 := CreateNode("test node 1", "com.collartechs.test", passThroughSignalProcessor{})
	node2 := CreateNode("test node 2", "com.collartechs.test", passThroughSignalProcessor{})

	node1.
		To("node2", node2)

	node2.Observe(func(node Node, when string, signal Signal, data ...interface{}) error {
		fmt.Println("observed onReceive", signal.Payload)
		if when == "onReceive" {
			str, existed := signal.Get(AnonPayload)
			assert.Equal(t, "test message", str.(string))
			assert.True(t, existed)
		}
		return nil
	})

	node1.Push("test message")
}

/* private method tests */

func TestParseNameFromComment(t *testing.T) {
	comment := "@This_is_name This is comment"
	name := parseNameFromComment(comment)
	assert.Equal(t, "This_is_name", name)

	// test dash in name
	comment = "@this-is-name this is comment"
	name = parseNameFromComment(comment)
	assert.Equal(t, "this-is-name", name)

	// test number in name
	comment = "@this-is-name-123 this is comment"
	name = parseNameFromComment(comment)
	assert.Equal(t, "this-is-name-123", name)
}

func TestParseTagsFromComment(t *testing.T) {
	comment := "@This_is_name #tag1 #tag_2 #tag-3 this is comment"
	tags := parseTagsFromComment(comment)
	assert.Equal(t, "tag1", tags[0])
	assert.Equal(t, "tag_2", tags[1])
	assert.Equal(t, "tag-3", tags[2])
}

func TestParseInfoFromComment(t *testing.T) {
	comment := "@this-is-name #tag1 #tag2 this is comment"
	name, tags, newComment := parseInfoFromComment(comment)

	assert.Equal(t, "this-is-name", name)
	assert.Equal(t, "tag1", tags[0])
	assert.Equal(t, "tag2", tags[1])
	assert.Equal(t, "this is comment", newComment)
}
