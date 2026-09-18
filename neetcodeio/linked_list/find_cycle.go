package linkedlist

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */

/*
 * GREEN
 * Time: O(n)
 * Space: O(1) due to recursion + map. Could be O(1) with a for loop though and pointers.
 * Challenge: https://neetcode.io/problems/linked-list-cycle-detection/question
 */
func hasCycle(head *ListNode) bool {
	// visit every node
	// record the node is visited
	// go ahead, check if the node has already been visited
	// then that's a cycle

	// should this be the pointer instead of val?
	//
	var slow, fast *ListNode

	slow = head
	fast = head

	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next

		if slow == fast {
			return true
		}
	}

	return false
}
