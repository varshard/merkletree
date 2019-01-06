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
