package merkletree

import (
	"crypto/sha256"
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
