# MerkleTree (WIP)

Merkle Tree implemented in Go.
This is a port of https://github.com/cliftonm/MerkleTree from C# to Go.

This project doesn't duplicate an odd node like some other implementation

## Example

```go
import "github.com/varshard/merkletree"

tree := Tree{}

func main() {
	leaves := []*Node{
		NewNode([]byte("1")),
		NewNode([]byte("2")),
		NewNode([]byte("3")),
	}

	tree.BuildTree(leaves)

	// Get root's hash
	fmt.Println(tree.Root.Hash)
}
```