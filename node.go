package merkletree

// TODO: Make a custom UnmashalJSON

import (
	"crypto/sha256"
	"encoding/json"
)

// Node is a node in a merkle tree
type Node struct {
	Hash   []byte `json:"hash"`
	Left   *Node  `json:"left"`
	Right  *Node  `json:"right"`
	parent *Node
}

// NewNode create a leave node and hash its data
func NewNode(data []byte) *Node {
	n := Node{
		Hash: computeHash(data),
	}

	return &n
}

// FromJSON will construct the Node from parsed raw json
func (n *Node) FromJSON(raw []byte) error {
	if err := json.Unmarshal(raw, n); err != nil {
		return err
	}

	populateParent(n, n.Left)
	populateParent(n, n.Right)

	return nil
}

func populateParent(parent *Node, node *Node) {
	if node != nil {
		node.parent = parent
		populateParent(node, node.Left)
		populateParent(node, node.Right)
	}
}

// NewParentNode construct a new parent node out of leave nodes
func NewParentNode(left *Node, right *Node) *Node {
	n := Node{}
	n.Left = left
	n.Right = right
	left.parent = &n

	if right != nil {
		right.parent = &n
	}

	n.Hash = computeConcatinatedHash(left, right)
	return &n
}

func computeHash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func computeConcatinatedHash(left *Node, right *Node) []byte {
	// Carry up left node if the tree has ood node
	if right == nil {
		return left.Hash
	}

	rightHash := right.Hash
	hash := append(left.Hash, rightHash...)
	return computeHash(hash)
}
