package linkedlist

type ListNode struct {
	Val  int
	Next *ListNode
}

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */

func reverseList(head *ListNode) *ListNode {
	// 1 -> 2 -> 3, 2 -> 1, but then where's the reference of 3?
	// need to keep track of them
	//
	var previous *ListNode
	current := head

	// nil -> head -> next -> next.next -> nil
	// prev   curr
	// current -> prev (can't lose reference to next - maybe store it)
	// head -> nil
	// nil -> head -> next -> next.next -> nil
	// 		  prev   curr
	// 1st: next = 2 (1 -> 2 -> 3) (2 is saved)
	// 2nd: current = 1, previous = nil, 1 -> nil, 2 -> 3 -> nil
	// 3rd: previous = current, current = next, prev = 1, current = 2

	// 1st: next = (current.next) 3
	// 2nd: current = 2, prev = 1, 2 -> 1 -> nil, 3 -> nil
	// 3rd: previous = current, current = next

	// 1st: next = (current.next) null
	// 2nd: current = 3, previous = 2, 3 -> 2 -> 1 -> nil
	// 3rd: previous = current, current = next  (3 -> 2 -> 1 -> nil)

	// 1 - nil, next = nil, current.next = nil, prev = 1, current = nil, 1 -> nil
	// 1 -> 2, next = 2, current.next = nil, prev = 1, current = 2, 1 -> nil, 2 -> nil
	// 1 -> nil, 2 -> nil, next = nil, current.next = 1, previous 2, current = nil
	// 2 -> 1 -> nil
	// 1 -> 2 -> 3
	// next = 2, current.next = nil, prev = 1, current = 2, 1 -> nil, 2 -> 3 -> nil
	// next = 3, current.next = 1, prev = 2, current = 3, 2 -> 1 -> nil, 3 -> nil
	// next = nil, current.next = 2, prev = 3, current = nil, 3 -> 2 -> 1 -> nil
	for current != nil {
		next := current.Next
		current.Next = previous
		previous = current
		current = next
	}

	return previous
}
