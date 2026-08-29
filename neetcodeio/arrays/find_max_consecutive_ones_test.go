package arrays

import "testing"

// Spec: longest run of consecutive 1s. Constraints say len >= 1 and values are
// only 0 or 1, but empty is covered anyway since it costs nothing.
func TestFindMaxConsecutiveOnes(t *testing.T) {
	cases := []struct {
		nums []int
		want int
	}{
		// examples from the statement
		{[]int{1, 1, 0, 1, 1, 1}, 3},
		{[]int{1, 0, 1, 1, 0, 1}, 2},

		// degenerate inputs
		{[]int{}, 0},
		{nil, 0},
		{[]int{1}, 1},
		{[]int{0}, 0},

		// no 1s at all
		{[]int{0, 0, 0}, 0},

		// all 1s
		{[]int{1, 1, 1}, 3},

		// the longest run is at the START -> must survive later resets
		{[]int{1, 1, 1, 0, 1}, 3},

		// the longest run is at the END -> must be counted without a trailing 0
		{[]int{1, 0, 1, 1}, 2},
		{[]int{0, 1, 1, 1}, 3},

		// the longest run is in the MIDDLE
		{[]int{0, 0, 1, 1, 1, 0, 0}, 3},
		{[]int{1, 0, 1, 1, 1, 0, 1}, 3},

		// runs of length 1 separated by 0s
		{[]int{1, 0, 1}, 1},
		{[]int{0, 1, 0}, 1},
		{[]int{1, 0, 1, 0, 1}, 1},

		// leading and trailing zeros
		{[]int{0, 1, 1, 0}, 2},
		{[]int{1, 1, 0, 1}, 2},

		// two runs of equal length -> the tie must not double-count
		{[]int{1, 1, 0, 1, 1}, 2},

		// several runs, longest is neither first nor last
		{[]int{1, 0, 0, 1, 1, 0, 1, 1, 1, 0}, 3},
		{[]int{0, 1, 1, 1, 1, 1, 0, 1}, 5},

		// consecutive zeros must not accumulate anything
		{[]int{0, 0, 1, 0, 0}, 1},
	}

	for _, c := range cases {
		if got := findMaxConsecutiveOnes(c.nums); got != c.want {
			t.Errorf("findMaxConsecutiveOnes(%v) = %d, want %d", c.nums, got, c.want)
		}
	}
}

// The constraint allows 100,000 elements, so make sure a long run is counted
// correctly and nothing overflows or resets partway.
func TestFindMaxConsecutiveOnesLong(t *testing.T) {
	n := 100000
	nums := make([]int, n)
	for i := range nums {
		nums[i] = 1
	}
	if got := findMaxConsecutiveOnes(nums); got != n {
		t.Errorf("all ones, len %d = %d, want %d", n, got, n)
	}

	// one 0 in the middle splits it into two halves
	nums[n/2] = 0
	if got := findMaxConsecutiveOnes(nums); got != n/2 {
		t.Errorf("split at %d = %d, want %d", n/2, got, n/2)
	}
}
