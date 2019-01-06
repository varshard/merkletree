package merkletree

import (
	"crypto/sha256"
)

// Node is a node in a merkle tree
type Node struct {
	Hash   []byte
	Left   *Node
	Right  *Node
	Parent *Node
}

// NewNode create a leave node and hash its data
func NewNode(data []byte) *Node {
	n := Node{
		Hash: computeHash(data),
	}

	return &n
}

// NewParentNode construct a new parent node out of leave nodes
func NewParentNode(left *Node, right *Node) *Node {
	n := Node{}
	n.Left = left
	n.Right = right
	left.Parent = &n

	if right != nil {
		right.Parent = &n
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
