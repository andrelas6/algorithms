package tree

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

/*
 * GREEN
 * Time: O(n)
 * Space: O(h) height
 */
func invertTree(root *TreeNode) *TreeNode {
	// invert children of root and proceed
	// this a level traversal
	// look at root, get its children, invert pointers, go to them
	// BFS seems more applicable
	// recursion
	// root.Left, root.Right = root.Right, root.Left
	// dfs(root.Left) dfs(root.Right)

	invertTreeDFS(root)

	return root
}

func invertTreeDFS(node *TreeNode) {
	if node == nil {
		return
	}

	node.Left, node.Right = node.Right, node.Left

	invertTreeDFS(node.Left)
	invertTreeDFS(node.Right)
}
