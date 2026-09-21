package linkedlist

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */

/*
 * RED: pointer assignments are definitely my weakness. I couldn't keep track of them during the challenge. My logic was sound, but the implementation with the pointers
 * was definitely not optimal. I needed to check the solution.
 * Time: O(n)
 * Space: O(1)
 */
func reorderList(head *ListNode) {
	// find mid point to "break up" list in two
	if head.Next == nil {
		return
	}
	var slow, fast *ListNode
	slow = head
	fast = head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}

	// from midpoint, revert linked list
	current := slow.Next // 8
	slow.Next = nil
	var previous *ListNode
	for current != nil {
		tmp := current.Next // nil
		current.Next = previous
		previous = current
		current = tmp
	}

	// iterate on head, put head of reverted between 1ist and second, then continue
	second := previous // head of the reverted list

	first := head
	for second != nil {
		tmp1, tmp2 := first.Next, second.Next
		first.Next = second
		second.Next = tmp1
		first, second = tmp1, tmp2
	}
}
