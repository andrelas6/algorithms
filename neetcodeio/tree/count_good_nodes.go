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
GREEN
Time: O(n)
Space: O(h)
Link: https://neetcode.io/problems/count-good-nodes-in-binary-tree/question
*/
func goodNodes(root *TreeNode) int {
	// dfs, passing the highest node down
	// compare curernt element to highes tnode
	// if node > highest: valid
	// if node < highest: invalid
	// goodNodesCount := 0
	return dfsHighest(root, -101)

}

func dfsHighest(node *TreeNode, highestVal int) int {
	// base cases
	if node == nil {
		return 0
	}

	result := 0
	if isValid(node.Val, highestVal) {
		result = 1
	} else {
		result = 0
	}

	result += dfsHighest(node.Left, max(node.Val, highestVal))
	result += dfsHighest(node.Right, max(node.Val, highestVal))

	return result
}

func isValid(nodeVal, highestVal int) bool {
	return nodeVal >= highestVal
}
