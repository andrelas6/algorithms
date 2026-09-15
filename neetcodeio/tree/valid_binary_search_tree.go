package tree

import "math"

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

//	  4
//	2   5
//
// 1 3    6
/*
 * RED - I only had the brute force solution in mind but didn't implemented. Needed to watch the solution video to solve. I knew DFS, but the trick to use a
 * left/right boundary to check the value I didn't know.
 * Time O(n)
 * Space O(h) because of recursion and height of tree
 */
func isValidBST(root *TreeNode) bool {
	leftBoundary := math.MinInt
	rightBoundary := math.MaxInt

	isTreeValid := validWithDFS(root, leftBoundary, rightBoundary)

	return isTreeValid
}

func validWithDFS(node *TreeNode, leftBoundary, rightBoundary int) bool {
	if node == nil {
		return true
	}

	isNodeValid := compare(node.Val, leftBoundary, rightBoundary)

	if !isNodeValid {
		return false
	}

	leftNodeValid := validWithDFS(node.Left, leftBoundary, node.Val)
	rightNodeValid := validWithDFS(node.Right, node.Val, rightBoundary)

	return leftNodeValid && rightNodeValid
}

func compare(value, leftBoundary, rightBoundary int) bool {
	return rightBoundary > value && leftBoundary < value
}
