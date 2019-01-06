package merkletree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildEmptyTree(t *testing.T) {
	tree := Tree{}
	leaves := []*Node{}

	assert.Error(t, tree.BuildTree(leaves))
}

func TestBuildTreeSingleNode(t *testing.T) {
	tree := Tree{}
	leaves := []*Node{
		NewNode([]byte("1")),
	}

	assert.NoError(t, tree.BuildTree(leaves))

	assert.NotNil(t, tree.Root)
	assert.Equal(t, tree.Root.Hash, leaves[0].Hash)
}

func TestBuildTree(t *testing.T) {
	tree := Tree{}
	leaves := []*Node{
		NewNode([]byte("1")),
		NewNode([]byte("2")),
		NewNode([]byte("3")),
	}

	assert.NoError(t, tree.BuildTree(leaves))

	assert.NotNil(t, tree.Root)
}

func TestAppendLeaf(t *testing.T) {
	tree := Tree{}
	tree.AppendLeaf(NewNode([]byte("1")))
	tree.AppendLeaf(NewNode([]byte("2")))

	for i, leafe := range tree.Leaves {
		assert.Equal(t, tree.Leaves[i], leafe)
	}
}
