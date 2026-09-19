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
 * YELLOW - didn't understood correctly the concept of a subtree. First solution found the first match of subRoot, but there could be duplicates. Then had to check the
 * solution.
 * Time: O(t * s), where t is one tree and s the other
 * Space O(h), where h is the height of root tree
 *
 */
func isSubtree(root *TreeNode, subRoot *TreeNode) bool {
	// is same tree? no, then is left same tree?, no, then is right same tree?
	if root == nil {
		return false
	}

	if dfsIsSame(root, subRoot) {
		return true
	} else {
		return isSubtree(root.Left, subRoot) || isSubtree(root.Right, subRoot)
	}
}

func dfsIsSame(n1 *TreeNode, n2 *TreeNode) bool {
	if n1 == nil && n2 == nil {
		return true
	}

	if n1 == nil || n2 == nil {
		return false
	}

	if n1.Val != n2.Val {
		return false
	}

	return dfsIsSame(n1.Left, n2.Left) && dfsIsSame(n1.Right, n2.Right)
}
