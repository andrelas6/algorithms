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
 * Yellow - I had a solution that was not optimal. Needed the hint to find the easy solution.
 * Time: O(h) where h is the height of the tree
 * Space: O(h) because each recursion goes a level
 */
func lowestCommonAncestor(root *TreeNode, p *TreeNode, q *TreeNode) *TreeNode {
	// create two slices, one for p and one for q
	// DFS to find both p & q (traversal type does not matter much)
	// slices contain order of visited root nodes
	// solution
	// DFS
	// go left, go right
	// when p is found, return slices of visited values
	// when q is found, return slices of visited
	// on main, iterate on the biggest slice, until sliceQ[i] != sliceQ[i]
	// then, return value

	// optimize
	// LCA is a lowest common ancestor, meaning that p and q have the same common parent
	// until they diverge (p < node and q > node), meaning that
	// 1. p is left, q is right. When that happens, node is a common ancestor
	// 2. but the other case: if p == node, then the LCA of q is node too
	// 1 or 2
	// traverse only if p & q > node or p & q < node
	// otherwise, value is either the parent or p/q, if p == node or q == node

	// p 3
	// q 4
	// root = 5
	if p.Val > root.Val && q.Val > root.Val {
		return lowestCommonAncestor(root.Right, p, q)
	} else if p.Val < root.Val && q.Val < root.Val {
		return lowestCommonAncestor(root.Left, p, q)
	} else {
		return root
	}
}

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

func lowestCommonAncestorMoreImproved(root *TreeNode, p *TreeNode, q *TreeNode) *TreeNode {
	// create two slices, one for p and one for q
	// DFS to find both p & q (traversal type does not matter much)
	// slices contain order of visited root nodes
	// solution
	// DFS
	// go left, go right
	// when p is found, return slices of visited values
	// when q is found, return slices of visited
	// on main, iterate on the biggest slice, until sliceQ[i] != sliceQ[i]
	// then, return value

	// optimize
	// LCA is a lowest common ancestor, meaning that p and q have the same common parent
	// until they diverge (p < node and q > node), meaning that
	// 1. p is left, q is right. When that happens, node is a common ancestor
	// 2. but the other case: if p == node, then the LCA of q is node too
	// 1 or 2
	// traverse only if p & q > node or p & q < node
	// otherwise, value is either the parent or p/q, if p == node or q == node

	// p 3
	// q 4
	// root = 5

	// optimize 2
	// recursion keeps frames in the stack, which for the function, is memory
	// to make this O(1) space, use a for loop instead.
	// for loop works fine as well
	// for needs no condition because p & q are guaranteed to be there
	// so they will be found before the traversal reaches a null value
	node := root
	for {
		if p.Val > node.Val && q.Val > node.Val {
			node = node.Right
		} else if p.Val < node.Val && q.Val < node.Val {
			node = node.Left
		} else {
			return node
		}
	}

}
