package linkedlist

import (
	"slices"
	"testing"
)

func TestMergeTwoListsExamples(t *testing.T) {
	cases := []struct {
		name         string
		list1, list2 []int
		want         []int
	}{
		// samples from the statement
		{"neetcode sample", []int{1, 2, 4}, []int{1, 3, 5}, []int{1, 1, 2, 3, 4, 5}},
		{"leetcode sample", []int{1, 2, 4}, []int{1, 3, 4}, []int{1, 1, 2, 3, 4, 4}},
		{"both empty", []int{}, []int{}, []int{}},
		{"first empty", []int{}, []int{0}, []int{0}},

		// degenerate
		{"second empty", []int{0}, []int{}, []int{0}},
		{"one node each, first smaller", []int{1}, []int{2}, []int{1, 2}},
		{"one node each, second smaller", []int{2}, []int{1}, []int{1, 2}},
		{"one node each, equal", []int{1}, []int{1}, []int{1, 1}},
	}

	for _, c := range cases {
		head1, _ := build(c.list1)
		head2, _ := build(c.list2)

		got := toSlice(t, mergeTwoLists(head1, head2), len(c.want)+5)
		if !slices.Equal(got, c.want) {
			t.Errorf("%s: mergeTwoLists(%v, %v) = %v, want %v", c.name, c.list1, c.list2, got, c.want)
		}
	}
}

// One list running out before the other is the case the loop condition decides.
// With && where || belongs, or a missing tail splice, the leftover nodes are
// dropped.
func TestMergeTwoListsLeftovers(t *testing.T) {
	cases := []struct {
		name         string
		list1, list2 []int
		want         []int
	}{
		{"first list runs out early", []int{1, 2}, []int{3, 4, 5, 6}, []int{1, 2, 3, 4, 5, 6}},
		{"second list runs out early", []int{3, 4, 5, 6}, []int{1, 2}, []int{1, 2, 3, 4, 5, 6}},
		{"all of the first is smaller", []int{1, 2, 3}, []int{4, 5, 6}, []int{1, 2, 3, 4, 5, 6}},
		{"all of the second is smaller", []int{4, 5, 6}, []int{1, 2, 3}, []int{1, 2, 3, 4, 5, 6}},
		{"one node against many", []int{5}, []int{1, 2, 3, 4}, []int{1, 2, 3, 4, 5}},
		{"many against one node", []int{1, 2, 3, 4}, []int{5}, []int{1, 2, 3, 4, 5}},
		{"leftover is a single node", []int{1, 2, 3}, []int{2, 9}, []int{1, 2, 2, 3, 9}},
		{"last value decides the tail", []int{1, 9}, []int{2, 3}, []int{1, 2, 3, 9}},
	}

	for _, c := range cases {
		head1, _ := build(c.list1)
		head2, _ := build(c.list2)

		got := toSlice(t, mergeTwoLists(head1, head2), len(c.want)+5)
		if !slices.Equal(got, c.want) {
			t.Errorf("%s: mergeTwoLists(%v, %v) = %v, want %v", c.name, c.list1, c.list2, got, c.want)
		}
	}
}

func TestMergeTwoListsValues(t *testing.T) {
	cases := []struct {
		name         string
		list1, list2 []int
		want         []int
	}{
		{"interleaved", []int{1, 3, 5, 7}, []int{2, 4, 6, 8}, []int{1, 2, 3, 4, 5, 6, 7, 8}},
		{"equal values throughout", []int{2, 2, 2}, []int{2, 2}, []int{2, 2, 2, 2, 2}},
		{"duplicates within one list", []int{1, 1, 5}, []int{1, 2}, []int{1, 1, 1, 2, 5}},
		{"negatives", []int{-100, -5, 0}, []int{-50, -1, 100}, []int{-100, -50, -5, -1, 0, 100}},
		{"range limits", []int{-100, 100}, []int{-100, 100}, []int{-100, -100, 100, 100}},
		{"zeroes", []int{0, 0}, []int{0}, []int{0, 0, 0}},
	}

	for _, c := range cases {
		head1, _ := build(c.list1)
		head2, _ := build(c.list2)

		got := toSlice(t, mergeTwoLists(head1, head2), len(c.want)+5)
		if !slices.Equal(got, c.want) {
			t.Errorf("%s: mergeTwoLists(%v, %v) = %v, want %v", c.name, c.list1, c.list2, got, c.want)
		}
	}
}

// The merged list must be built from the nodes it was given: no copies, and the
// placeholder node used to start the list must not end up in the result.
func TestMergeTwoListsReusesNodes(t *testing.T) {
	vals1 := []int{1, 4, 6}
	vals2 := []int{2, 3, 5}

	head1, nodes1 := build(vals1)
	head2, nodes2 := build(vals2)

	original := map[*ListNode]bool{}
	for _, n := range append(slices.Clone(nodes1), nodes2...) {
		original[n] = true
	}

	count := 0
	for node := mergeTwoLists(head1, head2); node != nil; node = node.Next {
		if !original[node] {
			t.Fatalf("node %d at position %d is not one of the input nodes", node.Val, count)
		}
		count++
	}

	if want := len(nodes1) + len(nodes2); count != want {
		t.Errorf("merged list has %d nodes, want %d", count, want)
	}
}

// The result is sorted and holds exactly the values of both inputs together,
// checked without rebuilding the merge.
func TestMergeTwoListsSortedAndComplete(t *testing.T) {
	inputs := []struct {
		list1, list2 []int
	}{
		{[]int{1, 2, 4}, []int{1, 3, 4}},
		{[]int{-3, -1, 7}, []int{-2, 0, 0, 9}},
		{[]int{5}, []int{-5, -4, -3, 100}},
		{[]int{1, 1, 1, 1}, []int{1, 1}},
		{[]int{}, []int{-7, 3}},
		{[]int{-7, 3}, []int{}},
	}

	for _, in := range inputs {
		head1, _ := build(in.list1)
		head2, _ := build(in.list2)

		got := toSlice(t, mergeTwoLists(head1, head2), len(in.list1)+len(in.list2)+5)

		if !slices.IsSorted(got) {
			t.Errorf("merging %v and %v gave %v, which is not sorted", in.list1, in.list2, got)
		}

		want := append(slices.Clone(in.list1), in.list2...)
		slices.Sort(want)
		gotSorted := slices.Clone(got)
		slices.Sort(gotSorted)

		if !slices.Equal(gotSorted, want) {
			t.Errorf("merging %v and %v gave %v, which is not both lists together", in.list1, in.list2, got)
		}
	}
}

// Constraints allow up to 50 nodes per list, values from -100 to 100.
func TestMergeTwoListsLong(t *testing.T) {
	const n = 50

	evens := make([]int, n)
	odds := make([]int, n)
	for i := range n {
		evens[i] = 2*i - 100
		odds[i] = 2*i - 99
	}

	head1, _ := build(evens)
	head2, _ := build(odds)

	got := toSlice(t, mergeTwoLists(head1, head2), 2*n+5)

	want := append(slices.Clone(evens), odds...)
	slices.Sort(want)

	if !slices.Equal(got, want) {
		t.Errorf("merging two lists of %d: got %v, want %v", n, got, want)
	}

	// one list entirely below the other, so the whole second list is spliced on
	low := make([]int, n)
	high := make([]int, n)
	for i := range n {
		low[i] = i - 100
		high[i] = i + 50
	}

	head1, _ = build(low)
	head2, _ = build(high)

	got = toSlice(t, mergeTwoLists(head1, head2), 2*n+5)
	want = append(slices.Clone(low), high...)

	if !slices.Equal(got, want) {
		t.Errorf("merging %d low values with %d high values: got %v, want %v", n, n, got, want)
	}
}
