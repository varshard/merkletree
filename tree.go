package merkletree

// TODO: Make a custom UnmashalJSON

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Tree a merkle tree structure
type Tree struct {
	Root   *Node   `json:"root"`
	Leaves []*Node `json:"leaves"`
}

// BuildTree build a tree out of Leaves
func (t *Tree) BuildTree() error {
	return t.buildTree(t.Leaves)
}

// buildTree build a tree out of a slice of leaves
func (t *Tree) buildTree(leaves []*Node) error {
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

		t.buildTree(parents)
	}

	t.Leaves = leaves
	return nil
}

// AuditProof returns the audit proof hashes to reconstruct the root hash.
func (t *Tree) AuditProof(hash []byte) ([]*ProofHash, error) {
	var auditTrail []*ProofHash
	var err error

	leafNode := t.FindLeaf(hash)

	if leafNode != nil {
		if leafNode.parent == nil {
			return nil, fmt.Errorf("expected leaf hash to have a parent hash")
		}

		parent := leafNode.parent
		auditTrail, err = t.BuildAuditTrail(auditTrail, parent, leafNode)
		if err != nil {
			return nil, err
		}
	}

	return auditTrail, nil
}

// VerifyAudit verify that if we walk up the tree from a particular leaf, we encounter the expected root hash.
func (t *Tree) VerifyAudit(rootHash []byte, targetHash []byte, auditTrail []*ProofHash) bool {
	testHash := targetHash

	for _, proof := range auditTrail {
		if proof.Direction == RightBranch {
			testHash = computeHash(append(testHash, proof.Hash...))
		} else {
			testHash = computeHash(append(proof.Hash, testHash...))
		}

	}

	return reflect.DeepEqual(rootHash, testHash)
}

// Verify will verify that the rootHash and the targetHash are valid.
func (t *Tree) Verify(rootHash []byte, targetHash []byte) bool {
	auditTrail, err := t.AuditProof(targetHash)
	if err != nil {
		return false
	}

	return t.VerifyAudit(rootHash, targetHash, auditTrail)
}

// BuildAuditTrail will build a trail composed of hash that are required to replicate the root hash
func (t *Tree) BuildAuditTrail(auditTrail []*ProofHash, parent *Node, child *Node) ([]*ProofHash, error) {
	if parent != nil {
		if child.parent != parent {
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

		return t.BuildAuditTrail(auditTrail, parent.parent, parent)
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

// FromJSON build a tree from JSON
func (t *Tree) FromJSON(raw []byte) error {
	if err := json.Unmarshal(raw, t); err != nil {
		return err
	}
	t.BuildTree()

	return nil
}
