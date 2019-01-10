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
	}

	rootHash := computeHash(append(leaves[0].Hash, leaves[1].Hash...))

	assert.NoError(t, tree.BuildTree(leaves))
	assert.Equal(t, tree.Root.Hash, rootHash)
}

func TestBuildTreeOddNode(t *testing.T) {
	tree := Tree{}
	leaves := []*Node{
		NewNode([]byte("1")),
		NewNode([]byte("2")),
		NewNode([]byte("3")),
	}

	firstRootHash := computeHash(append(leaves[0].Hash, leaves[1].Hash...))
	rootHash := computeHash(append(firstRootHash, leaves[2].Hash...))

	assert.NoError(t, tree.BuildTree(leaves))
	assert.Equal(t, tree.Root.Hash, rootHash)
}

func TestAppendLeaf(t *testing.T) {
	tree := Tree{}
	tree.AppendLeaf(NewNode([]byte("1")))
	tree.AppendLeaf(NewNode([]byte("2")))

	for i, leaf := range tree.Leaves {
		assert.Equal(t, tree.Leaves[i], leaf)
	}
}

func TestFindLeaf(t *testing.T) {
	target := NewNode([]byte("2"))

	tree := Tree{}
	tree.AppendLeaf(NewNode([]byte("1")))
	tree.AppendLeaf(target)
	leaf := tree.FindLeaf(target.Hash)

	assert.Equal(t, leaf.Hash, target.Hash)
}

func TestFindLeafOnBuildTree(t *testing.T) {
	target := NewNode([]byte("2"))
	leaves := []*Node{
		NewNode([]byte("1")),
		target,
	}

	tree := Tree{}
	tree.BuildTree(leaves)
	leaf := tree.FindLeaf(target.Hash)

	assert.Equal(t, leaf.Hash, target.Hash)
}

func TestFindLeafNotFound(t *testing.T) {
	target := computeHash([]byte("3"))

	tree := Tree{}
	tree.AppendLeaf(NewNode([]byte("1")))
	tree.AppendLeaf(NewNode([]byte("2")))

	leaf := tree.FindLeaf(target)

	assert.Nil(t, leaf)
}

func TestBuildAuditTrailOddLeaves(t *testing.T) {
	auditTrail := []*ProofHash{}

	tree := Tree{}
	leaves := []*Node{
		NewNode([]byte("1")),
		NewNode([]byte("2")),
		NewNode([]byte("3")),
	}
	assert.NoError(t, tree.BuildTree(leaves))

	target := leaves[1]
	auditTrail, err := tree.BuildAuditTrail(auditTrail, target.Parent, target)
	assert.NoError(t, err)

	assert.Equal(t, leaves[0].Hash, auditTrail[0].Hash)
	assert.Equal(t, LeftBranch, auditTrail[0].Direction)
	assert.Equal(t, leaves[2].Hash, auditTrail[1].Hash)
	assert.Equal(t, RightBranch, auditTrail[1].Direction)
}

func TestBuildAuditTrailEventLeaves(t *testing.T) {
	auditTrail := []*ProofHash{}

	tree := Tree{}
	leaves := []*Node{
		NewNode([]byte("1")),
		NewNode([]byte("2")),
		NewNode([]byte("3")),
		NewNode([]byte("4")),
	}
	assert.NoError(t, tree.BuildTree(leaves))

	target := leaves[2]
	auditTrail, err := tree.BuildAuditTrail(auditTrail, target.Parent, target)
	assert.NoError(t, err)

	leaf12Hash := append(leaves[0].Hash, leaves[1].Hash...)[:]
	leaf12 := NewNode(leaf12Hash)

	assert.Equal(t, leaves[3].Hash, auditTrail[0].Hash)
	assert.Equal(t, RightBranch, auditTrail[0].Direction)
	assert.Equal(t, leaf12.Hash, auditTrail[1].Hash)
	assert.Equal(t, LeftBranch, auditTrail[1].Direction)
}

func TestAuditProof(t *testing.T) {
	tree := Tree{}
	leaves := []*Node{
		NewNode([]byte("1")),
		NewNode([]byte("2")),
		NewNode([]byte("3")),
		NewNode([]byte("4")),
	}
	assert.NoError(t, tree.BuildTree(leaves))

	auditTrail, err := tree.AuditProof(leaves[2].Hash)
	assert.NoError(t, err)

	leaf12Hash := append(leaves[0].Hash, leaves[1].Hash...)[:]
	leaf12 := NewNode(leaf12Hash)

	assert.Equal(t, leaves[3].Hash, auditTrail[0].Hash)
	assert.Equal(t, RightBranch, auditTrail[0].Direction)
	assert.Equal(t, leaf12.Hash, auditTrail[1].Hash)
	assert.Equal(t, LeftBranch, auditTrail[1].Direction)
}
