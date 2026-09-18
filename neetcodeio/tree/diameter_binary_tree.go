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
 * YELLOW: tricky problem because your dfs calculates two different values. One is the height, the othe is the diameter. Then it does two
 * unintiuitive things: it keeps a global state variable for the max diameter while returning the height. Tricky problem.
 * Time: O(n)
 * Space: O(h)
 *
 */
func diameterOfBinaryTree(root *TreeNode) int {
	var result int

	dfs(root, &result)

	return result
}

func dfs(node *TreeNode, result *int) int {
	if node == nil {
		return 0
	}

	left := dfs(node.Left, result)
	right := dfs(node.Right, result)

	diameter := left + right

	// stores biggest diameter
	*result = max(diameter, *result)

	// returns the height
	return 1 + max(left, right)
}
