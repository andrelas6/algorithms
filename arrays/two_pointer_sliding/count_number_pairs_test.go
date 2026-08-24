package twopointersliding

import "testing"

// Spec under test: count index pairs (i, j) with i < j and prices[i]+prices[j] <= budget.
// prices is sorted non-decreasing, n <= 1000, values and budget up to 1e9.
// n < 2 -> 0.
func TestCountAffordablePairs(t *testing.T) {
	cases := []struct {
		prices []int32
		budget int32
		want   int32
	}{
		// the worked example from the statement
		{[]int32{1, 2, 3, 4, 5}, 7, 8},

		// samples 0 and 1 from the statement
		{[]int32{}, 100, 0},
		{[]int32{5}, 5, 0},

		// nil slice is the same as empty
		{nil, 100, 0},

		// smallest real cases
		{[]int32{1, 2}, 3, 1},
		{[]int32{1, 2}, 2, 0},
		{[]int32{1, 1}, 2, 1},
		{[]int32{1, 1}, 1, 0},

		// every pair qualifies -> n*(n-1)/2
		{[]int32{1, 2, 3, 4, 5}, 100, 10},
		{[]int32{1, 1, 1}, 2, 3},
		{[]int32{1, 1, 1, 1}, 2, 6},
		{[]int32{2, 2, 2, 2}, 4, 6},

		// no pair qualifies
		{[]int32{1, 2, 3, 4, 5}, 2, 0},
		{[]int32{10, 20, 30}, 15, 0},

		// partial: the valid pairs are a staircase, not a prefix of the array
		{[]int32{1, 2, 3, 4}, 7, 6},
		{[]int32{1, 2, 3, 4}, 5, 4},
		{[]int32{1, 3, 5, 7}, 8, 4},
		{[]int32{1, 2, 3, 4, 5, 6}, 7, 9},

		// a small element pairs with far elements while the middle ones do not
		{[]int32{1, 5, 6}, 7, 2},

		// duplicates around the boundary
		{[]int32{1, 2, 2, 3}, 5, 6},

		// upper end of the value range: 1e9 + 1e9 must not be mishandled
		{[]int32{1000000000, 1000000000}, 1000000000, 0},
		{[]int32{1, 1000000000}, 1000000000, 0},
		{[]int32{1, 999999999}, 1000000000, 1},
	}

	for _, c := range cases {
		if got := countAffordablePairs(c.prices, c.budget); got != c.want {
			t.Errorf("countAffordablePairs(%v, %d) = %d, want %d", c.prices, c.budget, got, c.want)
		}
	}
}
