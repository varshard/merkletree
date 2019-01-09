package merkletree

// Branch enum indicate branch direction (left/right)
type Branch int

// Direction of branches
const (
	// OldRoot used for linear list of hashes to compute the old root in a consistency proof.
	OldRoot     Branch = 0
	LeftBranch  Branch = 1
	RightBranch Branch = 2
)

// ProofHash is a hash for auditing
type ProofHash struct {
	Hash      []byte
	Direction Branch
}
