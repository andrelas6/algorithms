package linkedlist

import (
	"slices"
	"testing"
)

// interleaveEnds is an independent oracle for the expected order: take the
// first value, then the last, then the second, then the second to last, and so
// on. It works on a plain slice, so it shares nothing with the pointer work
// under test.
func interleaveEnds(vals []int) []int {
	out := make([]int, 0, len(vals))
	for i, j := 0, len(vals)-1; i <= j; i, j = i+1, j-1 {
		out = append(out, vals[i])
		if i != j {
			out = append(out, vals[j])
		}
	}
	return out
}

func checkReorder(t *testing.T, name string, vals []int) {
	t.Helper()

	head, nodes := build(vals)
	reorderList(head)

	got := toSlice(t, head, len(vals)+5)
	want := interleaveEnds(vals)

	if !slices.Equal(got, want) {
		t.Errorf("%s: reorderList(%v) gave %v, want %v", name, vals, got, want)
		return
	}

	// the list must be rebuilt from the nodes it was given, not from copies
	given := map[*ListNode]bool{}
	for _, n := range nodes {
		given[n] = true
	}
	seen := 0
	for n := head; n != nil; n = n.Next {
		if !given[n] {
			t.Errorf("%s: a node in the result was not one of the original nodes", name)
			return
		}
		seen++
	}
	if seen != len(nodes) {
		t.Errorf("%s: result has %d nodes, want %d", name, seen, len(nodes))
	}
}

func TestReorderListExamples(t *testing.T) {
	// samples from the statement
	checkReorder(t, "sample 1", []int{1, 2, 3, 4})     // -> 1 4 2 3
	checkReorder(t, "sample 2", []int{1, 2, 3, 4, 5})  // -> 1 5 2 4 3
	checkReorder(t, "sample 3", []int{2, 4, 6, 8})     // -> 2 8 4 6
	checkReorder(t, "sample 4", []int{2, 4, 6, 8, 10}) // -> 2 10 4 8 6
}

func TestReorderListSmallInputs(t *testing.T) {
	checkReorder(t, "single node", []int{1})
	checkReorder(t, "two nodes", []int{1, 2})
	checkReorder(t, "three nodes", []int{1, 2, 3})
	checkReorder(t, "six nodes", []int{1, 2, 3, 4, 5, 6})
	checkReorder(t, "seven nodes", []int{1, 2, 3, 4, 5, 6, 7})
}

// Both parities, over a range of lengths: the midpoint lands differently for
// odd and even counts, and that is where this problem usually breaks.
func TestReorderListEveryLength(t *testing.T) {
	for n := 1; n <= 20; n++ {
		vals := make([]int, n)
		for i := range vals {
			vals[i] = i + 1
		}
		checkReorder(t, "length "+string(rune('a'+n-1)), vals)
	}
}

func TestReorderListValues(t *testing.T) {
	checkReorder(t, "duplicates", []int{1, 1, 1, 1, 1})
	checkReorder(t, "two values alternating", []int{1, 2, 1, 2, 1, 2})
	checkReorder(t, "negatives", []int{-1, -2, -3, -4, -5})
	checkReorder(t, "across zero", []int{-2, -1, 0, 1, 2})
	checkReorder(t, "range limits", []int{0, 1000, 0, 1000})
	checkReorder(t, "descending", []int{9, 7, 5, 3, 1})
}

// The reordered list must end. A missing `slow.Next = nil` leaves the two
// halves pointing at each other, which loops forever.
func TestReorderListTerminates(t *testing.T) {
	for _, n := range []int{2, 3, 4, 5, 8, 9} {
		vals := make([]int, n)
		for i := range vals {
			vals[i] = i
		}

		head, _ := build(vals)
		reorderList(head)

		count := 0
		for node := head; node != nil; node = node.Next {
			count++
			if count > n+5 {
				t.Fatalf("length %d: the list does not end after %d nodes", n, count)
			}
		}
		if count != n {
			t.Errorf("length %d: walked %d nodes, want %d", n, count, n)
		}
	}
}

// Values must not be rewritten: the problem says to relink nodes, not to copy
// values between them. Each node keeps the value it started with.
func TestReorderListMovesNodesNotValues(t *testing.T) {
	vals := []int{10, 20, 30, 40, 50, 60}

	head, nodes := build(vals)
	valueOf := map[*ListNode]int{}
	for i, n := range nodes {
		valueOf[n] = vals[i]
	}

	reorderList(head)

	for n := head; n != nil; n = n.Next {
		if n.Val != valueOf[n] {
			t.Errorf("a node's value changed from %d to %d", valueOf[n], n.Val)
		}
	}
}

// Constraints allow 5 * 10^4 nodes.
func TestReorderListLong(t *testing.T) {
	for _, n := range []int{10_000, 50_000, 49_999} {
		vals := make([]int, n)
		for i := range vals {
			vals[i] = i
		}

		head, _ := build(vals)
		reorderList(head)

		got := toSlice(t, head, n+5)
		want := interleaveEnds(vals)

		if !slices.Equal(got, want) {
			t.Errorf("length %d: reordering did not interleave the ends", n)
			break
		}
	}
}
