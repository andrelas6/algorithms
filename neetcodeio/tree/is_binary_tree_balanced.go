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
 * YELLOW: I understood the concept of balanced tree incorrectly so I got it wrong. Then as soon as I understood it, I realised I had to use DFS.
 * Then the other weird part is using DFS to return two variables. That's all right but I decided to use a pointer to store the variable.
 * Time: O(n)
 * Space: O(h) where h is the height (recursion)
 *
 */
func isBalanced(root *TreeNode) bool {
	// BFS
	// on each level, if there is a null, store that vlaue
	// then if there is another level, return false
	// otherwise, t rue (default)

	// bfs on the root node
	// node, add to queue
	// loop on queue
	// remove node from uqeue, add its children fi they are not null
	// if any children is null, if there is a next level, return false

	// I though tbalance was about nodes in a level, but it's actually about height
	// that's DFS
	if root == nil {
		return true
	}

	var unbalancedFound bool
	left := dfsHeight(root.Left, &unbalancedFound)
	right := dfsHeight(root.Right, &unbalancedFound)

	if unbalancedFound {
		return false
	}

	return calcHeightDiff(left, right) <= 1
}

func calcHeightDiff(a, b int) int {
	if a > b {
		return a - b
	} else {
		return b - a
	}
}

// return height of tree
func dfsHeight(node *TreeNode, unbalancedFound *bool) int {
	if node == nil {
		return 0
	}

	leftHeight := dfsHeight(node.Left, unbalancedFound)
	rightHeight := dfsHeight(node.Right, unbalancedFound)

	if calcHeightDiff(leftHeight, rightHeight) > 1 {
		*unbalancedFound = true
	}

	return max(leftHeight, rightHeight) + 1
}
