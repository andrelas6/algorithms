package search

import (
	"math/rand/v2"
	"slices"
	"testing"
)

// The statement guarantees nums is sorted ascending with DISTINCT values, so
// every target has exactly one correct answer: its index, or -1. That makes the
// oracle trivial (slices.Index) and lets the tests assert an exact int rather
// than a property.

func TestSearchExamples(t *testing.T) {
	nums := []int{-1, 0, 2, 4, 6, 8}

	if got := search(nums, 4); got != 3 {
		t.Errorf("search(%v, 4) = %d, want 3", nums, got)
	}
	if got := search(nums, 3); got != -1 {
		t.Errorf("search(%v, 3) = %d, want -1", nums, got)
	}
}

// Every element of a given array must be findable at its own index. Missing the
// FIRST or LAST element is the classic symptom of a bad loop bound (`left < right`
// instead of `<=`) or a bad recursion bound (`middle` instead of `middle±1`), and
// those two positions are only reachable when the search window has narrowed to
// one element.
func TestSearchFindsEveryElement(t *testing.T) {
	cases := [][]int{
		// constraints say length >= 1, but empty costs nothing to cover
		{},
		{5},
		{1, 2},
		{1, 2, 3},
		{1, 2, 3, 4},
		// even length: middle never lands exactly in the centre
		{10, 20, 30, 40, 50, 60},
		// odd length
		{10, 20, 30, 40, 50, 60, 70},
		// negatives and zero, spanning the sign boundary
		{-9999, -500, -1, 0, 1, 500, 9999},
		// large gaps between neighbours
		{-1000, 0, 1000, 2000, 9999},
		// tightly packed consecutive values
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		// all negative
		{-50, -40, -30, -20, -10},
	}

	for _, nums := range cases {
		for i, v := range nums {
			if got := search(nums, v); got != i {
				t.Errorf("search(%v, %d) = %d, want %d", nums, v, got, i)
			}
		}
	}
}

// Absent targets must return -1, not the index of a near neighbour. The three
// interesting kinds of absence are: below everything, above everything, and in a
// gap between two present values -- the last one being where an implementation
// that returns `left` or `middle` on failure gives a plausible-looking wrong
// answer instead of -1.
func TestSearchMissingTargets(t *testing.T) {
	cases := []struct {
		nums   []int
		absent []int
	}{
		{[]int{}, []int{0, 1, -1}},
		{[]int{5}, []int{4, 6, 0, -5}},
		{[]int{1, 3}, []int{0, 2, 4}},
		{[]int{1, 3, 5}, []int{0, 2, 4, 6}},
		{[]int{2, 4, 6, 8}, []int{1, 3, 5, 7, 9}},
		// gaps at both ends and in the middle
		{[]int{-1, 0, 2, 4, 6, 8}, []int{-2, 1, 3, 5, 7, 9, -9999, 9999}},
		// just outside the bounds by one
		{[]int{10, 20, 30}, []int{9, 11, 19, 21, 29, 31}},
	}

	for _, c := range cases {
		for _, target := range c.absent {
			if got := search(c.nums, target); got != -1 {
				t.Errorf("search(%v, %d) = %d, want -1 (target is not present)", c.nums, target, got)
			}
		}
	}
}

// The first and last positions get their own test because they are the two the
// window-narrowing logic reaches last, and a bug there is easy to miss when the
// array is small enough that the midpoint happens to land on them early.
func TestSearchBoundaries(t *testing.T) {
	for n := 1; n <= 64; n++ {
		nums := make([]int, n)
		for i := range nums {
			nums[i] = i * 3 // distinct, ascending, with gaps
		}

		if got := search(nums, nums[0]); got != 0 {
			t.Errorf("n=%d: search for first element %d = %d, want 0", n, nums[0], got)
		}
		if got := search(nums, nums[n-1]); got != n-1 {
			t.Errorf("n=%d: search for last element %d = %d, want %d", n, nums[n-1], got, n-1)
		}

		// one below the first and one above the last must both miss
		if got := search(nums, nums[0]-1); got != -1 {
			t.Errorf("n=%d: search below the range = %d, want -1", n, got)
		}
		if got := search(nums, nums[n-1]+1); got != -1 {
			t.Errorf("n=%d: search above the range = %d, want -1", n, got)
		}
		// a value in a gap must miss
		if n > 1 {
			if got := search(nums, nums[0]+1); got != -1 {
				t.Errorf("n=%d: search for a gap value = %d, want -1", n, got)
			}
		}
	}
}

// BinarySearch takes explicit inclusive bounds, so it can be pointed at a window
// of the array. Searching a sub-range must find only what is inside it and report
// -1 for values that exist in the array but lie outside the given bounds. This is
// what makes the recursive calls on [middle+1, right] and [left, middle-1] sound.
func TestBinarySearchRespectsBounds(t *testing.T) {
	nums := []int{0, 10, 20, 30, 40, 50, 60, 70}

	cases := []struct {
		name        string
		left, right int
	}{
		{"whole array", 0, 7},
		{"first half", 0, 3},
		{"second half", 4, 7},
		{"middle window", 2, 5},
		{"single element", 3, 3},
		// an empty range: the loop guard must reject it rather than loop or panic
		{"inverted range", 4, 3},
		{"inverted range at zero", 0, -1},
	}

	for _, c := range cases {
		for i, v := range nums {
			got := BinarySearch(nums, c.left, c.right, v)

			inWindow := i >= c.left && i <= c.right
			want := -1
			if inWindow {
				want = i
			}

			if got != want {
				t.Errorf("%s: BinarySearch(nums, %d, %d, %d) = %d, want %d (index %d is %sin [%d,%d])",
					c.name, c.left, c.right, v, got, want, i,
					map[bool]string{true: "", false: "not "}[inWindow], c.left, c.right)
			}
		}
	}
}

// Randomised sweep against slices.Index as the oracle. Every array is built from
// a sorted set of distinct values, matching the stated constraints; targets
// include both present and absent values so hits and misses are both exercised.
// Fixed seed, so a failure reproduces.
func TestSearchRandomised(t *testing.T) {
	r := rand.New(rand.NewPCG(11, 2024))

	for n := 0; n <= 60; n++ {
		for trial := range 20 {
			// distinct ascending values via a random strictly-increasing walk
			nums := make([]int, n)
			v := -9999 + r.IntN(100)
			for i := range nums {
				nums[i] = v
				v += 1 + r.IntN(5) // gaps of 1..5 leave room for absent targets
			}

			// probe every value in a range that straddles the array
			lo, hi := -10000, 9999
			if n > 0 {
				lo, hi = nums[0]-3, nums[n-1]+3
			}
			for target := lo; target <= hi && target-lo < 200; target++ {
				want := slices.Index(nums, target)

				if got := search(nums, target); got != want {
					t.Fatalf("n=%d trial=%d: search(%v, %d) = %d, want %d",
						n, trial, nums, target, got, want)
				}
			}
		}
	}
}

// Constraints allow 10000 elements. This is where O(log n) matters: a linear or
// accidentally-quadratic implementation still returns the right answer, so the
// point here is to confirm correctness holds at the full input size and that the
// recursion depth (~log2(10000) = 14) stays sane.
func TestSearchLong(t *testing.T) {
	n := 10000

	// values are distinct and ascending, spread across the allowed range
	nums := make([]int, n)
	for i := range nums {
		nums[i] = -9999 + i*2
	}

	for i := range nums {
		if got := search(nums, nums[i]); got != i {
			t.Fatalf("n=%d: search for %d = %d, want %d", n, nums[i], got, i)
		}
	}

	// every odd value between the entries is absent
	for i := 0; i < n-1; i++ {
		if got := search(nums, nums[i]+1); got != -1 {
			t.Fatalf("n=%d: search for absent %d = %d, want -1", n, nums[i]+1, got)
		}
	}

	if got := search(nums, nums[0]-1); got != -1 {
		t.Errorf("below range = %d, want -1", got)
	}
	if got := search(nums, nums[n-1]+1); got != -1 {
		t.Errorf("above range = %d, want -1", got)
	}
}

// The array must not be modified: search is a read-only operation, and a caller
// reusing the slice afterwards depends on that.
func TestSearchDoesNotModifyInput(t *testing.T) {
	nums := []int{-1, 0, 2, 4, 6, 8}
	before := slices.Clone(nums)

	for _, target := range []int{-1, 4, 8, 3, -100, 100} {
		search(nums, target)
	}

	if !slices.Equal(nums, before) {
		t.Errorf("search modified the input: got %v, want %v", nums, before)
	}
}

func BenchmarkSearch(b *testing.B) {
	n := 10000
	nums := make([]int, n)
	for i := range nums {
		nums[i] = i * 2
	}

	b.ResetTimer()
	for i := range b.N {
		search(nums, (i%n)*2)
	}
}
