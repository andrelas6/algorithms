package addtwonumbers

import (
	"runtime/debug"
	"testing"
)

func buildList(vals []int) *ListNode {
	dummy := &ListNode{}
	tail := dummy
	for _, v := range vals {
		tail.Next = &ListNode{Val: v}
		tail = tail.Next
	}
	return dummy.Next
}

func toSlice(head *ListNode) []int {
	out := []int{}
	for n := head; n != nil; n = n.Next {
		out = append(out, n.Val)
	}
	return out
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestAddTwoNumbers(t *testing.T) {
	tests := []struct {
		name     string
		l1, l2   []int
		expected []int
	}{
		// samples
		{"sample 1: 342 + 465", []int{2, 4, 3}, []int{5, 6, 4}, []int{7, 0, 8}},
		{"sample 2: 0 + 0 stays [0]", []int{0}, []int{0}, []int{0}},
		{"sample 3: 9999999 + 9999", []int{9, 9, 9, 9, 9, 9, 9}, []int{9, 9, 9, 9}, []int{8, 9, 9, 9, 0, 0, 0, 1}},

		// single nodes
		{"single, no carry", []int{2}, []int{3}, []int{5}},
		{"single, final carry", []int{5}, []int{5}, []int{0, 1}},
		{"single, 9 + 9", []int{9}, []int{9}, []int{8, 1}},

		// no extra trailing node when there's no final carry
		{"no final carry, no trailing zero", []int{1, 2}, []int{3, 4}, []int{4, 6}},

		// different lengths
		{"l1 longer than l2", []int{1, 2, 3}, []int{4}, []int{5, 2, 3}},
		{"l2 longer than l1", []int{4}, []int{1, 2, 3}, []int{5, 2, 3}},
		{"carry chain through longer list: 99 + 1", []int{9, 9}, []int{1}, []int{0, 0, 1}},
		{"carry chain, l2 longer: 1 + 99", []int{1}, []int{9, 9}, []int{0, 0, 1}},

		// zeros inside the number must survive
		{"zero in the middle: 101 + 0", []int{1, 0, 1}, []int{0}, []int{1, 0, 1}},
		{"result digit is 0 mid-list: 15 + 5", []int{5, 1}, []int{5}, []int{0, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panicked: %v\n%s", r, debug.Stack())
				}
			}()

			l1 := buildList(tt.l1)
			l2 := buildList(tt.l2)

			// inputs may be modified, so only the result is checked
			got := toSlice(addTwoNumbers(l1, l2))
			if !equalInts(got, tt.expected) {
				t.Errorf("addTwoNumbers(%v, %v): expected %v, got %v", tt.l1, tt.l2, tt.expected, got)
			}
		})
	}
}
