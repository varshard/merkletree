# MerkleTree

Merkle Tree implemented in Go.
This is a port of https://github.com/cliftonm/MerkleTree from C# to Go.

**Note:** This project doesn't duplicate an odd node like some other implementation

## Usages

### Import
```go
import "github.com/varshard/merkletree"
```

### Building a tree

```go
func main() {
	tree := Tree{}

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

### Verify that a hash and root are of the same tree
```go
target := sha256.Sum256([]byte("2")])

// return boolean
tree.verify(target)
```
