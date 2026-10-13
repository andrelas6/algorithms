package copyrandomlist

type Node struct {
	Val    int
	Next   *Node
	Random *Node
}

func copyRandomList(head *Node) *Node {
	// do one pass, copy node, add to hashtable oldnode -> newnode
	// do a second pass, get oldnode.Random (which points to the new node in the hash table), assign to newnode.Random

	originalToCopy := make(map[*Node]*Node)
	// create copy linked list without random
	curr := head
	dummy := Node{
		Val:    -1,
		Next:   nil,
		Random: nil,
	}
	var currCopy *Node
	currCopy = &dummy

	for curr != nil {
		cp := Node{
			Val:    curr.Val,
			Next:   nil,
			Random: nil,
		}
		currCopy.Next = &cp
		currCopy = currCopy.Next

		originalToCopy[curr] = &cp

		curr = curr.Next
	}

	// assign random
	curr = head
	currCopy = dummy.Next
	for curr != nil {
		cp := originalToCopy[curr.Random]

		if cp != nil {
			currCopy.Random = cp
		}

		curr = curr.Next
		currCopy = currCopy.Next
	}

	return dummy.Next
}
