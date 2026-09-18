package linkedlist

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */

/*
 * RED: Conceptually it was not hard to come up with a solution in my head, but implmenting it I found confusing. Two pointers, moving
 * on adjacent lists, was a bit complicated to keep track of. Then, I also confused again the correct boolean comparison. I put || instead of &&
 * on the for loop condition for no reason.
 * Time: O(n + m) -> n is the nodes of list 1 and m is nodes of list 2
 * Space: O(1) -> just pointers
 */
func mergeTwoLists(list1 *ListNode, list2 *ListNode) (result *ListNode) {
	// visit every node of every list, alternating
	// a.1, then b.1, then compare -> higher go to the tail
	// then proceed to a.2, compare with tail -> higher go to the tail
	// then proceed to b.2

	dummy := &ListNode{
		Val:  0,
		Next: nil,
	}

	var p1, p2 *ListNode
	p1 = list1
	p2 = list2

	result = dummy
	tail := dummy

	// edge case -> one list reaches the end, the other does not, no need to compare anymore
	// 1
	// nil
	// 1
	// 1 3 5 6
	// mistake: again on confusing || with &&
	for p1 != nil && p2 != nil {
		if p1.Val < p2.Val {
			tail.Next = p1
			tail = p1
			p1 = p1.Next
		} else {
			tail.Next = p2
			tail = p2
			p2 = p2.Next
		}
	}

	if p1 == nil {
		tail.Next = p2
	} else if p2 == nil {
		tail.Next = p1
	}

	return dummy.Next
}
