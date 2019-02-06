package merkletree

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildEmptyTree(t *testing.T) {
	tree := Tree{
		Leaves: []*Node{},
	}

	assert.Error(t, tree.BuildTree())
}

func TestBuildTreeSingleNode(t *testing.T) {
	leaves := []*Node{
		NewNode([]byte("1")),
	}
	tree := Tree{
		Leaves: leaves,
	}

	assert.NoError(t, tree.BuildTree())
	assert.NotNil(t, tree.Root)
	assert.Equal(t, tree.Root.Hash, leaves[0].Hash)
	assert.Equal(t, 1, len(tree.Leaves))
}

func TestBuildTree(t *testing.T) {
	leaves := []*Node{
		NewNode([]byte("1")),
		NewNode([]byte("2")),
	}
	tree := Tree{
		Leaves: leaves,
	}

	rootHash := computeHash(append(leaves[0].Hash, leaves[1].Hash...))

	assert.NoError(t, tree.BuildTree())
	assert.Equal(t, tree.Root.Hash, rootHash)
	assert.Equal(t, 2, len(tree.Leaves))
}

func TestBuildTreeOddNode(t *testing.T) {
	leaves := []*Node{
		NewNode([]byte("1")),
		NewNode([]byte("2")),
		NewNode([]byte("3")),
	}
	tree := Tree{
		Leaves: leaves,
	}

	firstRootHash := computeHash(append(leaves[0].Hash, leaves[1].Hash...))
	rootHash := computeHash(append(firstRootHash, leaves[2].Hash...))

	assert.NoError(t, tree.BuildTree())
	assert.Equal(t, tree.Root.Hash, rootHash)
	assert.Equal(t, 3, len(tree.Leaves))
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

	tree := Tree{
		Leaves: leaves,
	}
	tree.BuildTree()
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

	leaves := []*Node{
		NewNode([]byte("1")),
		NewNode([]byte("2")),
		NewNode([]byte("3")),
	}
	tree := Tree{
		Leaves: leaves,
	}

	assert.NoError(t, tree.BuildTree())

	target := leaves[1]
	auditTrail, err := tree.BuildAuditTrail(auditTrail, target.parent, target)

	assert.NoError(t, err)
	assert.Equal(t, leaves[0].Hash, auditTrail[0].Hash)
	assert.Equal(t, LeftBranch, auditTrail[0].Direction)
	assert.Equal(t, leaves[2].Hash, auditTrail[1].Hash)
	assert.Equal(t, RightBranch, auditTrail[1].Direction)
}

func TestBuildAuditTrailEventLeaves(t *testing.T) {
	auditTrail := []*ProofHash{}

	leaves := []*Node{
		NewNode([]byte("1")),
		NewNode([]byte("2")),
		NewNode([]byte("3")),
		NewNode([]byte("4")),
	}
	tree := Tree{
		Leaves: leaves,
	}

	assert.NoError(t, tree.BuildTree())

	target := leaves[2]
	auditTrail, err := tree.BuildAuditTrail(auditTrail, target.parent, target)

	assert.NoError(t, err)

	leaf12Hash := append(leaves[0].Hash, leaves[1].Hash...)[:]
	leaf12 := NewNode(leaf12Hash)

	assert.Equal(t, leaves[3].Hash, auditTrail[0].Hash)
	assert.Equal(t, RightBranch, auditTrail[0].Direction)
	assert.Equal(t, leaf12.Hash, auditTrail[1].Hash)
	assert.Equal(t, LeftBranch, auditTrail[1].Direction)
}

func TestAuditProof(t *testing.T) {
	leaves := []*Node{
		NewNode([]byte("1")),
		NewNode([]byte("2")),
		NewNode([]byte("3")),
		NewNode([]byte("4")),
	}
	tree := Tree{
		Leaves: leaves,
	}

	assert.NoError(t, tree.BuildTree())

	auditTrail, err := tree.AuditProof(leaves[2].Hash)

	assert.NoError(t, err)

	leaf12Hash := append(leaves[0].Hash, leaves[1].Hash...)[:]
	leaf12 := NewNode(leaf12Hash)

	assert.Equal(t, leaves[3].Hash, auditTrail[0].Hash)
	assert.Equal(t, RightBranch, auditTrail[0].Direction)
	assert.Equal(t, leaf12.Hash, auditTrail[1].Hash)
	assert.Equal(t, LeftBranch, auditTrail[1].Direction)
}

func TestAuditProofErrorIfRootSupplied(t *testing.T) {
	leaves := []*Node{
		NewNode([]byte("1")),
	}
	tree := Tree{
		Leaves: leaves,
	}

	assert.NoError(t, tree.BuildTree())

	auditTrail, err := tree.AuditProof(tree.Root.Hash)

	assert.Error(t, err)
	assert.Nil(t, auditTrail)
}

func TestAuditProofReturnNilIfHashNotFound(t *testing.T) {
	leaves := []*Node{
		NewNode([]byte("1")),
	}
	tree := Tree{
		Leaves: leaves,
	}

	assert.NoError(t, tree.BuildTree())

	auditTrail, err := tree.AuditProof(NewNode([]byte("2")).Hash)

	assert.NoError(t, err)
	assert.Nil(t, auditTrail)
}

func TestVerifyAudit(t *testing.T) {
	target := NewNode([]byte("2"))
	leaves := []*Node{
		NewNode([]byte("1")),
		target,
		NewNode([]byte("3")),
		NewNode([]byte("4")),
	}
	tree := Tree{
		Leaves: leaves,
	}

	assert.NoError(t, tree.BuildTree())

	rootHash := tree.Root.Hash
	auditTrail, err := tree.AuditProof(target.Hash)

	assert.NoError(t, err)
	assert.True(t, tree.VerifyAudit(rootHash, target.Hash, auditTrail))
}

func TestVerifyAuditFalseIfTargetIsNotInTheTree(t *testing.T) {
	target := NewNode([]byte("5"))
	leaves := []*Node{
		NewNode([]byte("1")),
		NewNode([]byte("2")),
		NewNode([]byte("3")),
		NewNode([]byte("4")),
	}
	tree := Tree{
		Leaves: leaves,
	}

	assert.NoError(t, tree.BuildTree())

	rootHash := tree.Root.Hash
	auditTrail, err := tree.AuditProof(target.Hash)

	assert.NoError(t, err)
	assert.False(t, tree.VerifyAudit(rootHash, target.Hash, auditTrail))
}

func TestVerify(t *testing.T) {
	target := NewNode([]byte("2"))
	leaves := []*Node{
		NewNode([]byte("1")),
		NewNode([]byte("2")),
		NewNode([]byte("3")),
		NewNode([]byte("4")),
	}
	tree := Tree{
		Leaves: leaves,
	}

	assert.NoError(t, tree.BuildTree())
	assert.True(t, tree.Verify(tree.Root.Hash, target.Hash))
}

func TestVerifyReturnTrueIfNotvalid(t *testing.T) {
	target := NewNode([]byte("5"))
	leaves := []*Node{
		NewNode([]byte("1")),
		NewNode([]byte("2")),
		NewNode([]byte("3")),
		NewNode([]byte("4")),
	}
	tree := Tree{
		Leaves: leaves,
	}

	assert.NoError(t, tree.BuildTree())
	assert.False(t, tree.Verify(tree.Root.Hash, target.Hash))
}

func TestMarshalTree(t *testing.T) {
	leaves := []*Node{
		NewNode([]byte("1")),
		NewNode([]byte("2")),
		NewNode([]byte("3")),
		NewNode([]byte("4")),
	}
	tree := Tree{
		Leaves: leaves,
	}

	assert.NoError(t, tree.BuildTree())

	marshaled, err := json.Marshal(tree)

	assert.NoError(t, err)

	unmarshaled := Tree{}
	assert.NoError(t, unmarshaled.FromJSON(marshaled))
	assert.Equal(t, tree, unmarshaled)
}
