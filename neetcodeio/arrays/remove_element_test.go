package arrays

import (
	"slices"
	"testing"
)

// The statement says "the order of the elements may be changed", so a test that
// compares nums[:k] element by element would reject a perfectly valid answer
// (e.g. the swap-with-last variant). Mirror the custom judge instead: assert
// k == len(expected), then sort the first k elements before comparing. Anything
// past index k is explicitly garbage and must not be inspected.
func checkRemoveElement(t *testing.T, nums []int, val int, want []int) {
	t.Helper()

	input := slices.Clone(nums)

	k := removeElement(nums, val)

	if k != len(want) {
		t.Errorf("removeElement(%v, %d) = %d, want %d", input, val, k, len(want))
		return
	}
	if k < 0 || k > len(nums) {
		t.Fatalf("removeElement(%v, %d) returned k = %d, out of range for len %d", input, val, k, len(nums))
	}

	prefix := slices.Clone(nums[:k])
	slices.Sort(prefix)

	sorted := slices.Clone(want)
	slices.Sort(sorted)

	if !slices.Equal(prefix, sorted) {
		t.Errorf("removeElement(%v, %d): nums[:%d] = %v (sorted %v), want %v", input, val, k, nums[:k], prefix, sorted)
	}
}

func TestRemoveElement(t *testing.T) {
	cases := []struct {
		nums []int
		val  int
		want []int
	}{
		// examples from the statement
		{[]int{3, 2, 2, 3}, 3, []int{2, 2}},
		{[]int{0, 1, 2, 2, 3, 0, 4, 2}, 2, []int{0, 1, 3, 0, 4}},

		// degenerate inputs -- constraints allow len == 0
		{[]int{}, 0, []int{}},
		{nil, 3, []int{}},
		{[]int{1}, 1, []int{}},
		{[]int{1}, 2, []int{1}},

		// val absent entirely -> nothing is removed, k == len(nums)
		{[]int{1, 2, 3}, 4, []int{1, 2, 3}},
		// val is outside nums' range: constraints allow val up to 100, nums[i] up to 50
		{[]int{1, 2, 3}, 100, []int{1, 2, 3}},

		// every element is val -> k == 0, and the write pointer must never advance
		{[]int{2, 2, 2}, 2, []int{}},
		{[]int{0, 0, 0, 0}, 0, []int{}},

		// val at the FRONT only: the first survivor is copied over a removed slot
		{[]int{5, 1, 2, 3}, 5, []int{1, 2, 3}},
		{[]int{5, 5, 1, 2}, 5, []int{1, 2}},

		// val at the BACK only: nothing moves, k just stops short
		{[]int{1, 2, 3, 5}, 5, []int{1, 2, 3}},
		{[]int{1, 2, 5, 5}, 5, []int{1, 2}},

		// val in the MIDDLE
		{[]int{1, 5, 2}, 5, []int{1, 2}},
		{[]int{1, 5, 5, 2}, 5, []int{1, 2}},

		// alternating, so read and write pointers drift apart steadily
		{[]int{1, 5, 1, 5, 1, 5}, 5, []int{1, 1, 1}},
		{[]int{5, 1, 5, 1, 5, 1}, 5, []int{1, 1, 1}},

		// duplicates among the SURVIVORS must all be kept -- this is not dedupe
		{[]int{4, 4, 5, 4, 4}, 5, []int{4, 4, 4, 4}},
		{[]int{7, 7, 7}, 3, []int{7, 7, 7}},

		// val == 0, the zero value: a solution that treats 0 as "empty" breaks here
		{[]int{0, 1, 0, 2, 0}, 0, []int{1, 2}},
		{[]int{1, 0, 2}, 0, []int{1, 2}},

		// survivors that happen to equal an index, to catch value/index mixups
		{[]int{0, 1, 2, 3, 4}, 2, []int{0, 1, 3, 4}},

		// longer mixed case
		{[]int{3, 1, 3, 2, 3, 3, 4, 3, 5}, 3, []int{1, 2, 4, 5}},
	}

	for _, c := range cases {
		checkRemoveElement(t, slices.Clone(c.nums), c.val, c.want)
	}
}

// "In-place" is the whole point: the caller keeps its own slice header and reads
// nums[:k] afterwards. A solution that builds a fresh slice and returns its length
// passes the table above (the table reads the slice it passed in... which is the
// same one) but would leave the caller's array untouched. Pin the backing array
// down explicitly.
func TestRemoveElementIsInPlace(t *testing.T) {
	nums := []int{0, 1, 2, 2, 3, 0, 4, 2}
	before := &nums[0]

	k := removeElement(nums, 2)

	if k != 5 {
		t.Fatalf("k = %d, want 5", k)
	}
	if &nums[0] != before {
		t.Error("removeElement reseated the caller's slice, expected a write into the original array")
	}
	if len(nums) != 8 {
		t.Errorf("len(nums) = %d, want 8: the array is not resized, only the first k entries matter", len(nums))
	}

	prefix := slices.Clone(nums[:k])
	slices.Sort(prefix)
	if !slices.Equal(prefix, []int{0, 0, 1, 3, 4}) {
		t.Errorf("nums[:5] sorted = %v, want [0 0 1 3 4]", prefix)
	}
}

// The forward two-pointer only ever writes to an index at or behind the reader,
// so no survivor can be destroyed before it is read. Verify the invariant holds
// across every possible val for a sizeable input rather than trusting the traces
// in the comments: for each val, the multiset of survivors must match a filter.
func TestRemoveElementAgainstFilter(t *testing.T) {
	base := []int{0, 3, 1, 3, 50, 0, 7, 3, 7, 1, 0, 22, 3, 50, 0}

	for val := 0; val <= 51; val++ {
		want := []int{}
		for _, n := range base {
			if n != val {
				want = append(want, n)
			}
		}
		checkRemoveElement(t, slices.Clone(base), val, want)
	}
}

// Constraints cap len at 100, but a long all-survivors run and a long all-val run
// are the two shapes where an off-by-one in the write pointer shows up as a wrong
// k rather than as wrong contents.
func TestRemoveElementLongRuns(t *testing.T) {
	n := 100

	allSurvive := make([]int, n)
	for i := range allSurvive {
		allSurvive[i] = 1
	}
	if k := removeElement(slices.Clone(allSurvive), 2); k != n {
		t.Errorf("no element equals val: k = %d, want %d", k, n)
	}

	allRemoved := make([]int, n)
	for i := range allRemoved {
		allRemoved[i] = 2
	}
	if k := removeElement(slices.Clone(allRemoved), 2); k != 0 {
		t.Errorf("every element equals val: k = %d, want 0", k)
	}

	// half and half: val in every other slot
	mixed := make([]int, n)
	for i := range mixed {
		if i%2 == 0 {
			mixed[i] = 2
		} else {
			mixed[i] = 9
		}
	}
	checkRemoveElement(t, mixed, 2, slices.Repeat([]int{9}, n/2))
}
