package tree

// Definition for a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func inorder(arr *[]int, root *TreeNode) {
	if root != nil {
		inorder(arr, root.Left)
		*arr = append(*arr, root.Val)
		inorder(arr, root.Right)
	}
}

func inorderTraversal(root *TreeNode) []int {
	arr := make([]int, 0)

	inorder(&arr, root)

	return arr
}
