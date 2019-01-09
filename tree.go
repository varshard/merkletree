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

	return nil
}

// BuildAuditTrail will build a trail composed of hash that are required to replicate the root hash
func (t *Tree) BuildAuditTrail(auditTrail []*ProofHash, parent *Node, child *Node) ([]*ProofHash, error) {
	if parent != nil {
		if child.Parent != parent {
			return nil, fmt.Errorf("parent of child is not expected parent")
		}

		sibling := parent.Left
		direction := RightBranch
		if parent.Left == child {
			sibling = parent.Right
			direction = LeftBranch
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
