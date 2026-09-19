package tree

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

// DFS on both nodes at the same time
// each recursion, compare both values
// if they are the smae, proceed, otherwise, return false
// base case -> the nodes are not both null or nodes have different values
// dfs returns boolean

/*
 * GREEN
 * Time: O(n)
 * Space: O(h)
 */
func isSameTree(p *TreeNode, q *TreeNode) bool {
	if p == nil && q == nil {
		return true
	}
	if p == nil || q == nil {
		return false
	}
	if p.Val != q.Val {
		return false
	}

	return isSameTree(p.Left, q.Left) && isSameTree(p.Right, q.Right)
}
