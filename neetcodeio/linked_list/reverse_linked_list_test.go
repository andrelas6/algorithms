package linkedlist

import (
	"slices"
	"testing"
)

// build turns a slice into a list and hands back the head plus every node in
// order. The node slice is what lets the tests assert identity: reversal must
// re-point the existing nodes, not allocate new ones.
func build(vals []int) (*ListNode, []*ListNode) {
	if len(vals) == 0 {
		return nil, nil
	}

	nodes := make([]*ListNode, len(vals))
	for i, v := range vals {
		nodes[i] = &ListNode{Val: v}
	}
	for i := 0; i < len(nodes)-1; i++ {
		nodes[i].Next = nodes[i+1]
	}

	return nodes[0], nodes
}

// toSlice walks the list, but refuses to walk forever. The classic bug in this
// problem is dropping the `current.Next = previous` line for the original head,
// which leaves 1 -> 2 -> 1 and loops. A plain `for n != nil` walk would hang the
// test binary until the go test timeout instead of failing, so bound it and
// report the cycle as a normal failure.
func toSlice(t *testing.T, head *ListNode, limit int) []int {
	t.Helper()

	out := []int{}
	for n := head; n != nil; n = n.Next {
		if len(out) > limit {
			t.Fatalf("list did not terminate after %d nodes: the reversed list has a cycle (got %v...)", limit, out)
		}
		out = append(out, n.Val)
	}

	return out
}

func TestReverseList(t *testing.T) {
	cases := []struct {
		in   []int
		want []int
	}{
		// example from the statement
		{[]int{0, 1, 2, 3}, []int{3, 2, 1, 0}},

		// degenerate: constraints allow an empty list
		{nil, []int{}},
		{[]int{1}, []int{1}},

		// two nodes is the smallest case where anything actually moves
		{[]int{1, 2}, []int{2, 1}},

		// three nodes: the first case with a genuine middle node, so `next` has
		// to be saved before the link is rewritten
		{[]int{1, 2, 3}, []int{3, 2, 1}},

		{[]int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},

		// duplicate values: a solution that collects values and rebuilds passes
		// this, but TestReverseListReusesNodes below will not let it through
		{[]int{7, 7, 7}, []int{7, 7, 7}},
		{[]int{1, 2, 1, 2}, []int{2, 1, 2, 1}},

		// a palindrome reverses to itself -- catches a no-op implementation only
		// in combination with the other cases, kept for the traversal check
		{[]int{1, 2, 3, 2, 1}, []int{1, 2, 3, 2, 1}},

		// negatives and zero: constraints allow -5000 <= Val <= 5000
		{[]int{-1, 0, 1}, []int{1, 0, -1}},
		{[]int{0, 0}, []int{0, 0}},
		{[]int{-5000, 5000}, []int{5000, -5000}},

		// already descending, so a broken implementation cannot pass by accident
		{[]int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
	}

	for _, c := range cases {
		head, _ := build(c.in)

		got := toSlice(t, reverseList(head), len(c.in)+10)

		if !slices.Equal(got, c.want) {
			t.Errorf("reverseList(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

// The old head becomes the new tail, and its Next MUST be nil. Forgetting to
// rewrite that one link is the single most common bug here: the values still
// come out right for a couple of steps, then the walk loops forever. Assert it
// directly rather than relying on toSlice's bound to catch it.
func TestReverseListOldHeadBecomesTail(t *testing.T) {
	cases := [][]int{
		{1},
		{1, 2},
		{1, 2, 3},
		{1, 2, 3, 4, 5},
	}

	for _, vals := range cases {
		head, nodes := build(vals)
		oldHead := nodes[0]
		oldTail := nodes[len(nodes)-1]

		got := reverseList(head)

		if got != oldTail {
			t.Errorf("reverseList(%v) returned %v, want the old tail (Val %d)", vals, got, oldTail.Val)
			continue
		}
		if oldHead.Next != nil {
			t.Errorf("reverseList(%v): old head (Val %d) still points at Val %d, want nil", vals, oldHead.Val, oldHead.Next.Val)
		}
	}
}

// Reversal is a pointer rewrite, not a rebuild. The reversed list must contain
// exactly the original node pointers in reverse order -- same count, same
// identities, nothing allocated. This is what separates the intended O(1)-space
// solution from "read the values into a slice and build a new list", which is
// O(n) space and passes every value-based assertion above.
func TestReverseListReusesNodes(t *testing.T) {
	vals := []int{10, 20, 30, 40}
	head, nodes := build(vals)

	got := reverseList(head)

	seen := []*ListNode{}
	for n := got; n != nil; n = n.Next {
		if len(seen) > len(vals) {
			t.Fatal("reversed list has a cycle")
		}
		seen = append(seen, n)
	}

	if len(seen) != len(nodes) {
		t.Fatalf("reversed list has %d nodes, want %d", len(seen), len(nodes))
	}
	for i, n := range seen {
		want := nodes[len(nodes)-1-i]
		if n != want {
			t.Errorf("node %d of the reversed list is a different allocation (Val %d), want the original node holding Val %d", i, n.Val, want.Val)
		}
	}
}

// Reversing twice is the identity. Cheap property check that catches asymmetric
// off-by-one errors -- e.g. an implementation that drops the last node would
// shrink the list on every pass and fail here even if a single reversal happens
// to look plausible.
func TestReverseListTwiceIsIdentity(t *testing.T) {
	vals := []int{1, 2, 3, 4, 5, 6}
	head, nodes := build(vals)

	back := reverseList(reverseList(head))

	if back != nodes[0] {
		t.Fatalf("double reversal returned %v, want the original head (Val %d)", back, nodes[0].Val)
	}
	if got := toSlice(t, back, len(vals)+10); !slices.Equal(got, vals) {
		t.Errorf("double reversal = %v, want %v", got, vals)
	}
}

// Constraints cap the list at 5000 nodes. The iterative version handles that in
// O(1) space; a recursive one is fine at this depth too, but this pins down that
// nothing quadratic or stack-hungry sneaks in.
func TestReverseListLong(t *testing.T) {
	n := 5000

	vals := make([]int, n)
	for i := range vals {
		vals[i] = i
	}
	head, _ := build(vals)

	got := toSlice(t, reverseList(head), n+10)

	if len(got) != n {
		t.Fatalf("reversed list has %d nodes, want %d", len(got), n)
	}
	for i, v := range got {
		if want := n - 1 - i; v != want {
			t.Fatalf("node %d = %d, want %d", i, v, want)
		}
	}
}
