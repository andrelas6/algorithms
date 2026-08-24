package heightpointer

// Node is the ordinary pointer form of a binary tree node.
// A nil *Node means "there is no node here" — the equivalent of -1 in the
// packed-array form used by arrays/../trees/dfs_recursion.
type Node struct {
	Val         int32
	Left, Right *Node
}

// height returns the number of nodes on the longest path from n down to a leaf.
// A nil node has height 0; a leaf has height 1.
func height(n *Node) int32 {
	if n == nil {
		return 0
	}

	left := 1 + height(n.Left)
	right := 1 + height(n.Right)

	return max(left, right)
}
