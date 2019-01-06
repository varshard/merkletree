# MerkleTree implemented in Go

This is a port of https://github.com/cliftonm/MerkleTree from c# to Go.

This project doesn't duplicate an odd node like some other implementatoin

## Example
```go
tree := Tree{}
	  leaves := []*Node{
		NewNode([]byte("1")),
	}

tree.BuildTree(leaves)

// Get root's hash
fmt.Println(tree.Root.Hash)
```