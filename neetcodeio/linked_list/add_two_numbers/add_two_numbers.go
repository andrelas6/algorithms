package addtwonumbers

type ListNode struct {
	Val  int
	Next *ListNode
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	// create linked list starting with 0
	// walk every element, sum l1, l2 and l0
	// if sum > 10, do modulo, set next = 1
	// continue until end of l1/l2 -> if sum > 10, add 1 to result
	// O(n)
	dummy := &ListNode{
		Val:  0,
		Next: nil,
	}
	c1 := l1
	c2 := l2

	current := dummy
	// build result ll with zeroes
	// add leading zeroes to l1 and l2
	// for loop continues except both next are nil
	for {
		current.Next = newNode()
		current = current.Next

		if c1.Next == nil && c2.Next == nil {
			break
		}

		if c1.Next == nil {
			c1.Next = newNode()
		}

		if c2.Next == nil {
			c2.Next = newNode()
		}

		c1 = c1.Next
		c2 = c2.Next
	}

	current.Next = newNode()

	// do the sum
	c1 = l1
	c2 = l2
	c3 := dummy.Next

	for c1 != nil {
		result := c1.Val + c2.Val + c3.Val

		if result >= 10 {
			c3.Val = result % 10
			c3.Next.Val = 1
		} else {
			c3.Val = result
		}

		// update pointers
		c1 = c1.Next
		c2 = c2.Next
		c3 = c3.Next
	}

	// remove leading zero
	c3 = dummy
	for c3 != nil && c3.Next != nil {
		if c3.Next.Val == 0 && c3.Next.Next == nil {
			c3.Next = nil
		}

		c3 = c3.Next
	}

	return dummy.Next
}

func newNode() *ListNode {
	node := &ListNode{
		Val:  0,
		Next: nil,
	}

	return node
}
