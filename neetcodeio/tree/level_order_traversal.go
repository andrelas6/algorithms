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
 * Space: O(n + w) two times (queue w + result 2d array n)
 */
func levelOrder(root *TreeNode) [][]int {
	var result [][]int
	// null check heer
	if root == nil {
		return result
	}
	return levelOrderTraversal(root, result)
}

func levelOrderTraversal(root *TreeNode, result [][]int) [][]int {
	// queue
	// add root to queue and start
	// loop per level - add all to the queue
	// in the loop, get each left/right of each node inthe queue, add to queue
	// next loop processes these
	// BFS
	queue := make([]*TreeNode, 0)
	queue = append(queue, root)

	// while queue has elements, continue
	for len(queue) > 0 {
		// level nodes
		nodeValuesAtLevel := make([]int, 0)

		for _ = range len(queue) {
			// dequeue
			node := queue[0] // get fifo
			nodeValuesAtLevel = append(nodeValuesAtLevel, node.Val)
			queue = queue[1:] // remove from the queue

			// adding to queue
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		result = append(result, nodeValuesAtLevel)
	}

	return result
}
