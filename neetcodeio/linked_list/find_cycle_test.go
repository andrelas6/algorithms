package linkedlist

import "testing"

// buildWithCycle builds the list and, when pos >= 0, points the tail's Next at
// the node at that index. pos == -1 leaves the list ending in nil. This is the
// `pos` from the statement, which is not passed to the function.
func buildWithCycle(vals []int, pos int) *ListNode {
	head, nodes := build(vals)
	if pos >= 0 && len(nodes) > 0 {
		nodes[len(nodes)-1].Next = nodes[pos]
	}
	return head
}

func TestHasCycleExamples(t *testing.T) {
	cases := []struct {
		name string
		vals []int
		pos  int
		want bool
	}{
		// samples from the statement
		{"tail points back to index 1", []int{3, 2, 0, -4}, 1, true},
		{"two nodes pointing at each other", []int{1, 2}, 0, true},
		{"single node, no cycle", []int{1}, -1, false},

		// degenerate
		{"empty list", []int{}, -1, false},
		{"single node pointing at itself", []int{1}, 0, true},
		{"two nodes, no cycle", []int{1, 2}, -1, false},
	}

	for _, c := range cases {
		if got := hasCycle(buildWithCycle(c.vals, c.pos)); got != c.want {
			t.Errorf("%s: hasCycle(%v, pos=%d) = %v, want %v", c.name, c.vals, c.pos, got, c.want)
		}
	}
}

// Where the cycle starts must not matter: the whole list, the last node only,
// or anything between.
func TestHasCyclePosition(t *testing.T) {
	vals := []int{10, 20, 30, 40, 50}

	for pos := range vals {
		if !hasCycle(buildWithCycle(vals, pos)) {
			t.Errorf("cycle back to index %d: got false, want true", pos)
		}
	}

	if hasCycle(buildWithCycle(vals, -1)) {
		t.Error("no cycle: got true, want false")
	}
}

// The trap: a cycle is two nodes being the SAME node, not two nodes holding the
// same value. Every list here repeats values and none of them has a cycle.
func TestHasCycleRepeatedValues(t *testing.T) {
	cases := []struct {
		name string
		vals []int
	}{
		{"all the same value", []int{1, 1, 1, 1}},
		{"two of the same value", []int{1, 2, 1}},
		{"value repeated at the ends", []int{7, 3, 9, 7}},
		{"all zeroes", []int{0, 0, 0}},
		{"head value repeated next door", []int{5, 5}},
		{"long run of one value", []int{2, 2, 2, 2, 2, 2, 2, 2}},
	}

	for _, c := range cases {
		if hasCycle(buildWithCycle(c.vals, -1)) {
			t.Errorf("%s: hasCycle(%v) = true, want false (same values, different nodes)", c.name, c.vals)
		}
	}
}

// Repeated values AND a real cycle: the answer is still true, but not because
// of the repeats.
func TestHasCycleRepeatedValuesWithCycle(t *testing.T) {
	cases := []struct {
		name string
		vals []int
		pos  int
	}{
		{"all the same value, cycle at the head", []int{1, 1, 1}, 0},
		{"all the same value, cycle at the tail", []int{1, 1, 1}, 2},
		{"repeated values, cycle in the middle", []int{4, 9, 4, 9, 4}, 2},
	}

	for _, c := range cases {
		if !hasCycle(buildWithCycle(c.vals, c.pos)) {
			t.Errorf("%s: hasCycle(%v, pos=%d) = false, want true", c.name, c.vals, c.pos)
		}
	}
}

// Detecting a cycle is read-only: the values and the links must come back
// unchanged, so a later caller sees the same list.
func TestHasCycleDoesNotModifyList(t *testing.T) {
	vals := []int{1, 2, 3, 4, 5}

	for _, pos := range []int{-1, 0, 2, 4} {
		head, nodes := build(vals)
		if pos >= 0 {
			nodes[len(nodes)-1].Next = nodes[pos]
		}

		nextBefore := make([]*ListNode, len(nodes))
		valsBefore := make([]int, len(nodes))
		for i, n := range nodes {
			nextBefore[i], valsBefore[i] = n.Next, n.Val
		}

		hasCycle(head)

		for i, n := range nodes {
			if n.Next != nextBefore[i] {
				t.Errorf("pos=%d: node %d had its Next rewritten", pos, i)
			}
			if n.Val != valsBefore[i] {
				t.Errorf("pos=%d: node %d had its Val changed from %d to %d", pos, i, valsBefore[i], n.Val)
			}
		}
	}
}

// Constraints allow 10^4 nodes and values from -10^5 to 10^5.
func TestHasCycleLong(t *testing.T) {
	const n = 10_000

	vals := make([]int, n)
	for i := range vals {
		vals[i] = i - n/2
	}

	if hasCycle(buildWithCycle(vals, -1)) {
		t.Error("10^4 nodes, no cycle: got true, want false")
	}
	for _, pos := range []int{0, 1, n / 2, n - 2, n - 1} {
		if !hasCycle(buildWithCycle(vals, pos)) {
			t.Errorf("10^4 nodes, cycle back to index %d: got false, want true", pos)
		}
	}

	// every value the same, so only node identity can tell them apart
	same := make([]int, n)
	if hasCycle(buildWithCycle(same, -1)) {
		t.Error("10^4 identical values, no cycle: got true, want false")
	}
	if !hasCycle(buildWithCycle(same, n/2)) {
		t.Error("10^4 identical values with a cycle: got false, want true")
	}
}
