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
 * GREEN: Solved with BFS (neeed a bit of a refresher but not much - just checked my notes and voila).
 * Time: O(n)
 * Space: O(n) -> queue keeps adding each node but it's a slice, so backing array will contain the values. The returned result does not count as space afaik.
 */
func rightSideView(root *TreeNode) (result []int) {
	// loop to visit a level (loop leaves when queue is empty)
	// in loop, add nodes of level to the queue
	// remove from queue, process
	if root == nil {
		return result
	}

	queue := make([]*TreeNode, 0)

	queue = append(queue, root)

	for len(queue) != 0 {
		// visit node
		// add left if not null
		// add right if not null
		// how do I know that it's visible? get last element of queue after the level finishes in BFS
		result = append(result, queue[len(queue)-1].Val)
		for _ = range len(queue) {
			node := queue[0]
			queue = queue[1:]

			if node.Left != nil {
				queue = append(queue, node.Left)
			}

			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
	}

	return result
}
