package heap

import (
	"math/rand/v2"
	"slices"
	"testing"
	"time"
)

// Both implementations answer the same question, so every table runs against
// both. Each call gets its own copy: both reorder nums in place.
var kthLargestImpls = []struct {
	name string
	fn   func(nums []int, k int) int
}{
	{"findKthLargest", findKthLargest},
	{"findKthLargestImproved", findKthLargestImproved},
}

// isKthLargest is an independent oracle straight from the definition, with no
// sorting and no partitioning: v is the kth largest when fewer than k values are
// strictly greater than v, and at least k values are greater than or equal to it.
// Duplicates count separately, which is what the problem asks for.
func isKthLargest(nums []int, k, v int) bool {
	greater, greaterOrEqual := 0, 0
	for _, x := range nums {
		if x > v {
			greater++
		}
		if x >= v {
			greaterOrEqual++
		}
	}
	return greater < k && k <= greaterOrEqual
}

func TestFindKthLargestExamples(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		// samples from the statement
		{"neetcode sample 1", []int{2, 3, 1, 5, 4}, 2, 4},
		{"neetcode sample 2", []int{2, 3, 1, 1, 5, 5, 4}, 3, 4},
		{"leetcode sample 1", []int{3, 2, 1, 5, 6, 4}, 2, 5},
		{"leetcode sample 2", []int{3, 2, 3, 1, 2, 4, 5, 5, 6}, 4, 4},

		// degenerate
		{"single element", []int{7}, 1, 7},
		{"two elements, largest", []int{1, 2}, 1, 2},
		{"two elements, smallest", []int{1, 2}, 2, 1},

		// the ends of the range of k
		{"k = 1 is the maximum", []int{4, 9, 1, 7, 3}, 1, 9},
		{"k = n is the minimum", []int{4, 9, 1, 7, 3}, 5, 1},
		{"k in the middle", []int{4, 9, 1, 7, 3}, 3, 4},

		// input order must not matter
		{"already ascending", []int{1, 2, 3, 4, 5}, 2, 4},
		{"already descending", []int{5, 4, 3, 2, 1}, 2, 4},

		// negatives and zero: constraints allow -10^4 <= nums[i] <= 10^4
		{"all negative", []int{-5, -1, -3, -2, -4}, 2, -2},
		{"across zero", []int{-10000, 0, 10000, -1, 1}, 3, 0},
		{"range limits", []int{10000, -10000}, 2, -10000},
	}

	for _, impl := range kthLargestImpls {
		for _, c := range cases {
			if got := impl.fn(slices.Clone(c.nums), c.k); got != c.want {
				t.Errorf("%s: %s(%v, k=%d) = %d, want %d", c.name, impl.name, c.nums, c.k, got, c.want)
			}
		}
	}
}

// The trap in this problem: it is the kth largest in sorted order, NOT the kth
// largest distinct value. A solution that dedupes first passes on unique inputs
// and fails every one of these.
func TestFindKthLargestDuplicatesCount(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"all equal", []int{3, 3, 3, 3}, 3, 3},
		{"duplicate maximum", []int{5, 5, 1}, 2, 5},
		{"duplicate minimum", []int{1, 1, 5}, 3, 1},
		// distinct-values answer would be 1
		{"duplicate hides the next value", []int{2, 2, 2, 1}, 3, 2},
		// sorted: 1 2 2 3 3 3 -> 3 3 3 are the top three
		{"run of the largest", []int{3, 1, 3, 2, 3, 2}, 3, 3},
		{"just past the run", []int{3, 1, 3, 2, 3, 2}, 4, 2},
		{"many duplicates, pivot-heavy", []int{1, 2, 1, 2, 1, 2, 1, 2}, 5, 1},
	}

	for _, impl := range kthLargestImpls {
		for _, c := range cases {
			if got := impl.fn(slices.Clone(c.nums), c.k); got != c.want {
				t.Errorf("%s: %s(%v, k=%d) = %d, want %d", c.name, impl.name, c.nums, c.k, got, c.want)
			}
		}
	}
}

// Every k from 1 to n on a spread of inputs, checked against the definition.
// Small value ranges force lots of duplicates; a fixed seed keeps it repeatable.
func TestFindKthLargestEveryKAgainstDefinition(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	inputs := [][]int{
		{1},
		{2, 1},
		{5, 5, 5},
		{9, 8, 7, 6, 5, 4, 3, 2, 1},
		{1, 2, 3, 4, 5, 6, 7, 8, 9},
		{0, -1, 1, -2, 2, -3, 3},
	}
	for _, n := range []int{10, 31, 100} {
		for _, spread := range []int{3, 1000} {
			nums := make([]int, n)
			for i := range nums {
				nums[i] = rng.IntN(spread) - spread/2
			}
			inputs = append(inputs, nums)
		}
	}

	for _, impl := range kthLargestImpls {
		for _, nums := range inputs {
			for k := 1; k <= len(nums); k++ {
				got := impl.fn(slices.Clone(nums), k)
				if !isKthLargest(nums, k, got) {
					t.Errorf("%s(%v, k=%d) = %d, which is not the kth largest", impl.name, nums, k, got)
				}
			}
		}
	}
}

// Both functions rearrange nums in place. That is allowed, but they must only
// move values around: nothing lost, nothing duplicated.
func TestFindKthLargestKeepsTheSameValues(t *testing.T) {
	nums := []int{3, 2, 3, 1, 2, 4, 5, 5, 6, -7, 0}

	for _, impl := range kthLargestImpls {
		for k := 1; k <= len(nums); k++ {
			work := slices.Clone(nums)
			impl.fn(work, k)

			gotSorted, wantSorted := slices.Clone(work), slices.Clone(nums)
			slices.Sort(gotSorted)
			slices.Sort(wantSorted)

			if !slices.Equal(gotSorted, wantSorted) {
				t.Errorf("%s(k=%d): values after the call %v are not a rearrangement of %v", impl.name, k, work, nums)
			}
		}
	}
}

// Quickselect leaves the answer at its sorted position, index n-k, with nothing
// larger to its left and nothing smaller to its right. That is the invariant
// that lets it stop without sorting the rest.
func TestFindKthLargestImprovedPartitionsAroundAnswer(t *testing.T) {
	inputs := [][]int{
		{3, 2, 1, 5, 6, 4},
		{3, 2, 3, 1, 2, 4, 5, 5, 6},
		{7, 7, 1, 7, 2, 7},
		{10, -3, 8, 0, 8, -3, 5, 1},
	}

	for _, nums := range inputs {
		for k := 1; k <= len(nums); k++ {
			work := slices.Clone(nums)
			got := findKthLargestImproved(work, k)
			at := len(work) - k

			if work[at] != got {
				t.Errorf("nums=%v k=%d: returned %d but nums[%d] = %d after the call (%v)", nums, k, got, at, work[at], work)
				continue
			}
			for i := range at {
				if work[i] > got {
					t.Errorf("nums=%v k=%d: nums[%d] = %d is left of the answer %d but larger (%v)", nums, k, i, work[i], got, work)
				}
			}
			for i := at + 1; i < len(work); i++ {
				if work[i] < got {
					t.Errorf("nums=%v k=%d: nums[%d] = %d is right of the answer %d but smaller (%v)", nums, k, i, work[i], got, work)
				}
			}
		}
	}
}

// quickSelect is the recursive worker: it finds the value that belongs at index
// k of nums[left..right] once sorted, and must leave everything outside that
// window alone.
func TestQuickSelectStaysInsideWindow(t *testing.T) {
	nums := []int{100, 5, 1, 4, 2, 3, -100}

	for k := 1; k <= 5; k++ {
		work := slices.Clone(nums)

		if got := quickSelect(work, 1, 5, k); got != k {
			t.Errorf("quickSelect(window 1..5, k=%d) = %d, want %d", k, got, k)
		}
		if work[0] != 100 || work[6] != -100 {
			t.Errorf("quickSelect(window 1..5, k=%d) touched values outside the window: %v", k, work)
		}
	}
}

// Constraints allow n up to 10^5. The random input is the comfortable case. The
// ascending and all-equal inputs are the ones that decide whether a quickselect
// that always takes the last element as pivot runs in linear or quadratic time,
// so they get a time budget. Skipped with -short.
func TestFindKthLargestAtMaxSize(t *testing.T) {
	if testing.Short() {
		t.Skip("large inputs skipped in -short mode")
	}

	const n = 100_000
	const budget = time.Second

	rng := rand.New(rand.NewPCG(3, 4))

	random := make([]int, n)
	ascending := make([]int, n)
	allEqual := make([]int, n)
	for i := range n {
		random[i] = rng.IntN(20_001) - 10_000
		ascending[i] = i - n/2
		allEqual[i] = 42
	}

	inputs := []struct {
		name string
		nums []int
	}{
		{"random", random},
		{"ascending", ascending},
		{"all equal", allEqual},
	}

	for _, impl := range kthLargestImpls {
		for _, in := range inputs {
			for _, k := range []int{1, n / 2, n} {
				start := time.Now()
				got := impl.fn(slices.Clone(in.nums), k)
				elapsed := time.Since(start)

				if !isKthLargest(in.nums, k, got) {
					t.Errorf("%s on %s n=%d k=%d: got %d, which is not the kth largest", impl.name, in.name, n, k, got)
				}
				if elapsed > budget {
					t.Errorf("%s on %s n=%d k=%d: took %v, over the %v budget", impl.name, in.name, n, k, elapsed.Round(time.Millisecond), budget)
				}
			}
		}
	}
}
