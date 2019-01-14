package merkletree

import (
	"crypto/sha256"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	data := []byte("hello, world")
	expectation := sha256.Sum256(data)

	node := NewNode(data)

	assert.Equal(t, expectation[:], node.Hash)
}

func TestNewParentNode(t *testing.T) {
	left := Node{
		Hash: []byte("1"),
	}
	right := Node{
		Hash: []byte("2"),
	}

	n := NewParentNode(&left, &right)
	assert.Equal(t, n.Left.Hash, left.Hash)
	assert.Equal(t, n.Right.Hash, right.Hash)

	concatedHash := sha256.Sum256(append(left.Hash, right.Hash...))
	assert.Equal(t, n.Hash, concatedHash[:])
}

func TestNewParentWithOddNode(t *testing.T) {
	left := Node{
		Hash: []byte("1"),
	}

	n := NewParentNode(&left, nil)
	assert.NotNil(t, n.Left)
	assert.Nil(t, n.Right)
	assert.Equal(t, n.Hash, left.Hash)
}

func TestMarshalNode(t *testing.T) {
	n := Node{
		Hash: []byte("1"),
	}

	marshaled, err := json.Marshal(n)
	assert.NoError(t, err)

	unmarshaledNode := Node{}
	unmarshaledNode.FromJSON(marshaled)
	assert.Equal(t, n, unmarshaledNode)
}

func TestMarshalNestedNode(t *testing.T) {
	left := Node{
		Hash: []byte("1"),
	}
	right := Node{
		Hash: []byte("2"),
	}

	n := NewParentNode(&left, &right)

	marshaled, err := json.Marshal(n)
	assert.NoError(t, err)

	unmarshaledNode := Node{}
	assert.NoError(t, unmarshaledNode.FromJSON(marshaled))
	assert.Equal(t, n, &unmarshaledNode)
}

func TestMarshal2LevelNode(t *testing.T) {
	left1 := Node{
		Hash: []byte("1"),
	}
	right1 := Node{
		Hash: []byte("2"),
	}
	parentLeft := NewParentNode(&left1, &right1)

	left2 := Node{
		Hash: []byte("3"),
	}
	right2 := Node{
		Hash: []byte("4"),
	}
	parentRight := NewParentNode(&left2, &right2)

	n := NewParentNode(parentLeft, parentRight)

	marshaled, err := json.Marshal(n)
	assert.NoError(t, err)

	unmarshaledNode := Node{}
	assert.NoError(t, unmarshaledNode.FromJSON(marshaled))
	assert.Equal(t, n, &unmarshaledNode)
}
