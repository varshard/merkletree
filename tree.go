package merkletree

import (
	"fmt"
	"reflect"
)

// TODO: verification

// Tree a merkle tree structure
type Tree struct {
	Root   *Node
	Leaves []*Node
}

// BuildTree build a tree out of a slice of leaves
func (t *Tree) BuildTree(leaves []*Node) error {
	leavesCount := len(leaves)
	if leavesCount < 1 {
		return fmt.Errorf("Leaves is empty")
	} else if leavesCount == 1 {
		t.Root = leaves[0]
	} else {
		parents := []*Node{}

		for i := 0; i < leavesCount; i += 2 {
			left := leaves[i]
			var right *Node
			if i+1 == leavesCount {
				right = nil
			} else {
				right = leaves[i+1]
			}

			parents = append(parents, NewParentNode(left, right))
		}

		t.BuildTree(parents)
	}

	t.Leaves = leaves
	return nil
}

// AuditProof used to find an trail of hashes to prove that the supplied hash is a part of the tree
func (t *Tree) AuditProof(hash []byte) ([]*ProofHash, error) {
	var auditTrail []*ProofHash
	var err error

	leafNode := t.FindLeaf(hash)

	if leafNode != nil {
		if leafNode.Parent == nil {
			return nil, fmt.Errorf("expected leaf hash to have a parent hash")
		}

		parent := leafNode.Parent
		auditTrail, err = t.BuildAuditTrail(auditTrail, parent, leafNode)
		if err != nil {
			return nil, err
		}
	}

	return auditTrail, nil
}

// BuildAuditTrail will build a trail composed of hash that are required to replicate the root hash
func (t *Tree) BuildAuditTrail(auditTrail []*ProofHash, parent *Node, child *Node) ([]*ProofHash, error) {
	if parent != nil {
		if child.Parent != parent {
			return nil, fmt.Errorf("parent of child is not expected parent")
		}

		sibling := parent.Left
		direction := LeftBranch
		if parent.Left == child {
			sibling = parent.Right
			direction = RightBranch
		}

		proof := ProofHash{
			Hash:      sibling.Hash,
			Direction: direction,
		}

		auditTrail = append(auditTrail, &proof)

		return t.BuildAuditTrail(auditTrail, parent.Parent, parent)
	}
	return auditTrail, nil
}

// AppendLeaf append a leaf node to Leaves,
// to be built into a tree with BuildTree
func (t *Tree) AppendLeaf(leaf *Node) {
	t.Leaves = append(t.Leaves, leaf)
}

// FindLeaf find a leaf node that match supplied hash
func (t *Tree) FindLeaf(hash []byte) *Node {
	for _, leaf := range t.Leaves {
		if reflect.DeepEqual(leaf.Hash, hash) {
			return leaf
		}
	}

	return nil
}
