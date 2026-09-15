package hashmaps

import (
	"math/rand/v2"
	"slices"
	"testing"
	"time"
)

// longestConsecutiveBySorting is an independent oracle: sort, drop duplicates,
// then count the longest run where each value is one more than the last. No map
// and no sequence starters, so a bug in that logic cannot hide behind the same
// mistake here. O(n log n), which is fine for a test.
func longestConsecutiveBySorting(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	sorted := slices.Clone(nums)
	slices.Sort(sorted)
	sorted = slices.Compact(sorted)

	best, run := 1, 1
	for i := 1; i < len(sorted); i++ {
		if sorted[i] == sorted[i-1]+1 {
			run++
		} else {
			run = 1
		}
		best = max(best, run)
	}
	return best
}

func TestLongestConsecutiveExamples(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want int
	}{
		// samples from the statement
		{"neetcode sample 1", []int{2, 20, 4, 10, 3, 4, 5}, 4},
		{"neetcode sample 2", []int{0, 3, 2, 5, 4, 6, 1, 1}, 7},
		{"leetcode sample 1", []int{100, 4, 200, 1, 3, 2}, 4},
		{"leetcode sample 2", []int{0, 3, 7, 2, 5, 8, 4, 6, 0, 1}, 9},
		{"leetcode sample 3", []int{1, 0, 1, 2}, 3},

		// degenerate
		{"empty", []int{}, 0},
		{"nil", nil, 0},
		{"single element", []int{7}, 1},
		{"two consecutive", []int{8, 7}, 2},
		{"two not consecutive", []int{7, 9}, 1},
		{"nothing consecutive", []int{10, 30, 20, 50, 40}, 1},

		// input order must not matter
		{"ascending", []int{1, 2, 3, 4, 5}, 5},
		{"descending", []int{5, 4, 3, 2, 1}, 5},
		{"interleaved runs", []int{1, 11, 2, 12, 3, 13, 14}, 4},
	}

	for _, c := range cases {
		if got := longestConsecutive(c.nums); got != c.want {
			t.Errorf("%s: longestConsecutive(%v) = %d, want %d", c.name, c.nums, got, c.want)
		}
	}
}

// Duplicates must not stretch a sequence: [1, 2, 2, 3] is length 3, not 4. A
// solution that sorts and counts adjacent pairs without skipping repeats resets
// or overcounts on every one of these.
func TestLongestConsecutiveDuplicates(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want int
	}{
		{"all equal", []int{5, 5, 5, 5}, 1},
		{"duplicate in the middle", []int{1, 2, 2, 3}, 3},
		{"duplicate starter", []int{1, 1, 2, 3}, 3},
		{"duplicate end", []int{1, 2, 3, 3}, 3},
		{"every value twice", []int{3, 1, 2, 3, 2, 1}, 3},
		{"duplicates across two runs", []int{1, 1, 2, 10, 10, 11, 12, 12}, 3},
	}

	for _, c := range cases {
		if got := longestConsecutive(c.nums); got != c.want {
			t.Errorf("%s: longestConsecutive(%v) = %d, want %d", c.name, c.nums, got, c.want)
		}
	}
}

// Several runs: the longest can be anywhere, and a run can be broken by a single
// missing value.
func TestLongestConsecutiveSeveralRuns(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want int
	}{
		{"longest first", []int{1, 2, 3, 4, 10, 11, 20}, 4},
		{"longest last", []int{20, 10, 11, 1, 2, 3, 4}, 4},
		{"longest in the middle", []int{50, 1, 2, 3, 4, 5, 30, 31}, 5},
		{"ties", []int{1, 2, 3, 10, 11, 12}, 3},
		{"one gap splits a run", []int{1, 2, 3, 5, 6, 7, 8}, 4},
		{"gap filled later in the input", []int{1, 2, 3, 5, 6, 7, 4}, 7},
	}

	for _, c := range cases {
		if got := longestConsecutive(c.nums); got != c.want {
			t.Errorf("%s: longestConsecutive(%v) = %d, want %d", c.name, c.nums, got, c.want)
		}
	}
}

// Negatives and zero: constraints allow -10^9 <= nums[i] <= 10^9. Runs that
// cross zero, and values at the very limits, must behave like any other.
func TestLongestConsecutiveNegativesAndLimits(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want int
	}{
		{"crosses zero", []int{1, -1, 0}, 3},
		{"all negative", []int{-3, -1, -2, -10}, 3},
		{"run ends at -1", []int{-5, -4, -3, -2, -1, 1}, 5},
		{"min and max only", []int{-1_000_000_000, 1_000_000_000}, 1},
		{"run at the max", []int{999_999_999, 1_000_000_000, 999_999_998}, 3},
		{"run at the min", []int{-1_000_000_000, -999_999_999, -999_999_998}, 3},
	}

	for _, c := range cases {
		if got := longestConsecutive(c.nums); got != c.want {
			t.Errorf("%s: longestConsecutive(%v) = %d, want %d", c.name, c.nums, got, c.want)
		}
	}
}

// Random inputs against the sort-based oracle. Small value ranges force lots of
// duplicates and long runs; a wide range makes runs rare. Fixed seed keeps it
// repeatable.
func TestLongestConsecutiveAgainstSorting(t *testing.T) {
	rng := rand.New(rand.NewPCG(7, 8))

	for range 2000 {
		n := rng.IntN(40)
		spread := []int{5, 50, 1_000_000}[rng.IntN(3)]

		nums := make([]int, n)
		for i := range nums {
			nums[i] = rng.IntN(spread) - spread/2
		}

		got, want := longestConsecutive(nums), longestConsecutiveBySorting(nums)
		if got != want {
			t.Fatalf("longestConsecutive(%v) = %d, sorting oracle = %d", nums, got, want)
		}
	}
}

// The answer depends only on which values are present, so any shuffle of the
// same input gives the same answer.
func TestLongestConsecutiveIgnoresOrder(t *testing.T) {
	rng := rand.New(rand.NewPCG(9, 10))
	nums := []int{9, 1, 4, 7, 3, -1, 0, 5, 8, 2, 2, 6, 100}
	want := longestConsecutive(slices.Clone(nums))

	for range 50 {
		shuffled := slices.Clone(nums)
		rng.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

		if got := longestConsecutive(shuffled); got != want {
			t.Errorf("longestConsecutive(%v) = %d, but %d for the same values in another order", shuffled, got, want)
		}
	}
}

// Finding the length is read-only: the caller's slice must come back untouched.
func TestLongestConsecutiveDoesNotModifyInput(t *testing.T) {
	nums := []int{100, 4, 200, 1, 3, 2, 2}
	original := slices.Clone(nums)

	longestConsecutive(nums)

	if !slices.Equal(nums, original) {
		t.Errorf("longestConsecutive modified its input: got %v, want %v", nums, original)
	}
}

// Constraints allow n up to 10^5, and the whole point of the problem is O(n). Each
// input gets a time budget. The call runs in a goroutine so a quadratic run fails
// at the budget instead of holding up the suite for a minute. Skipped with -short.
func TestLongestConsecutiveAtMaxSize(t *testing.T) {
	if testing.Short() {
		t.Skip("large inputs skipped in -short mode")
	}

	const n = 100_000
	const budget = time.Second

	rng := rand.New(rand.NewPCG(11, 12))

	// one run of every value, shuffled
	oneRun := make([]int, n)
	for i := range oneRun {
		oneRun[i] = i - n/2
	}
	rng.Shuffle(n, func(i, j int) { oneRun[i], oneRun[j] = oneRun[j], oneRun[i] })

	// nothing consecutive: even numbers only
	noRuns := make([]int, n)
	for i := range noRuns {
		noRuns[i] = 2 * i
	}

	// every value the same
	allEqual := make([]int, n)
	for i := range allEqual {
		allEqual[i] = 42
	}

	// half the input is copies of the start of one long run: 0 repeated n/2
	// times, then 1..n/2
	repeatedStart := make([]int, 0, n)
	for range n / 2 {
		repeatedStart = append(repeatedStart, 0)
	}
	for v := 1; v <= n/2; v++ {
		repeatedStart = append(repeatedStart, v)
	}

	cases := []struct {
		name string
		nums []int
		want int
	}{
		{"one shuffled run", oneRun, n},
		{"no runs", noRuns, 1},
		{"all equal", allEqual, 1},
		{"start of a long run repeated", repeatedStart, n/2 + 1},
	}

	for _, c := range cases {
		done := make(chan int, 1)
		start := time.Now()
		go func() { done <- longestConsecutive(c.nums) }()

		select {
		case got := <-done:
			if got != c.want {
				t.Errorf("%s (n=%d): got %d, want %d", c.name, len(c.nums), got, c.want)
			}
			t.Logf("%s (n=%d): %v", c.name, len(c.nums), time.Since(start).Round(time.Millisecond))
		case <-time.After(budget):
			t.Errorf("%s (n=%d): still running after %v, over budget", c.name, len(c.nums), budget)
		}
	}
}
