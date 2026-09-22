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
Yellow: I knew I had to do in-order (left, root, right) traversal, but the trick of using a slice to order it I didn't know.
Time: O(n)
Space: O(h)
Link: https://neetcode.io/problems/kth-smallest-integer-in-bst/question
*/
func kthSmallest(root *TreeNode, k int) int {
	// dfs in order
	// left, val, right
	// add those elements to a slice
	// check k
	result := make([]int, 0)
	sorted := dfsInOrder(root, result)

	return sorted[k-1]
}

func dfsInOrder(node *TreeNode, result []int) []int {
	if node == nil {
		return result
	}

	result = dfsInOrder(node.Left, result)
	result = append(result, node.Val)
	result = dfsInOrder(node.Right, result)

	return result
}
