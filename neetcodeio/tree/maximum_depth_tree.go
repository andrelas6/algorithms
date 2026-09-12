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
 * Space: O(h) which would be O(n) for a very unbalanced tree
 */
func maxDepth(root *TreeNode) int {
	// level_right = traverse right
	// level_left = traverse left
	// return max(left, right)
	if root == nil {
		return 0
	}

	tree_depth_left := depthFrom(root.Left, 1)
	tree_depth_right := depthFrom(root.Right, 1)

	return max(tree_depth_left, tree_depth_right)
}

func depthFrom(root *TreeNode, level int) int {
	if root == nil {
		return level
	}

	level_left := depthFrom(root.Left, level+1)
	level_right := depthFrom(root.Right, level+1)

	return max(level_left, level_right)
}
