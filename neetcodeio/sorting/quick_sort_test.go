package sorting

import (
	"math/rand/v2"
	"slices"
	"testing"
)

// byKey, stableSorted, itoa and formatPairs live in the other test files in this
// package and are reused here.
//
// NOTE ON THE ORACLE. Quick Sort is explicitly NOT stable -- the statement says
// so and Example 1 demonstrates it. So stableSorted() is the wrong oracle for
// this problem: the order of equal keys is decided by the partition scheme, not
// by input order, and asserting stability would fail a correct solution.
//
// Instead the tests below assert what the statement actually requires:
//   - the two worked examples, matched EXACTLY (they pin the Lomuto scheme with
//     a right-most pivot and a strict `<` comparison)
//   - the output is non-decreasing by key
//   - the output is a permutation of the input (nothing lost or duplicated)
//   - the sort happens in place
//
// Together those pin down everything except the arrangement of tied values,
// which the two examples cover.

func isSortedByKey(pairs []Pair) bool {
	for i := 1; i < len(pairs); i++ {
		if pairs[i-1].Key > pairs[i].Key {
			return false
		}
	}
	return true
}

// checkPermutation reports every element that was lost or duplicated. Values are
// unique in these tests, so Pair is usable as a map key.
func checkPermutation(t *testing.T, in, got []Pair) {
	t.Helper()

	if len(got) != len(in) {
		t.Errorf("got %d elements, want %d\nin  %s\ngot %s", len(got), len(in), formatPairs(in), formatPairs(got))
		return
	}

	counts := make(map[Pair]int, len(in))
	for _, p := range in {
		counts[p]++
	}
	for _, p := range got {
		counts[p]--
	}
	for p, n := range counts {
		if n > 0 {
			t.Errorf("QuickSort(%s) dropped %v (%d times)", formatPairs(in), p, n)
		}
		if n < 0 {
			t.Errorf("QuickSort(%s) duplicated %v (%d times)", formatPairs(in), p, -n)
		}
	}
}

// The two worked examples, matched exactly. These are the only cases where the
// arrangement of EQUAL keys is asserted, and they are what distinguishes the
// required scheme (right-most pivot, strict `<`, Lomuto partition) from any
// other correct quicksort. A solution using `<=`, or a different pivot choice,
// still sorts -- but produces a different tie order and fails here.
func TestQuickSortExamples(t *testing.T) {
	cases := []struct {
		name string
		in   []Pair
		want []Pair
	}{
		{
			// the statement calls this out as proof the sort is NOT stable:
			// "bird" was after "cat" in the input and ends up before it
			name: "example 1",
			in:   []Pair{{3, "cat"}, {2, "dog"}, {3, "bird"}},
			want: []Pair{{2, "dog"}, {3, "bird"}, {3, "cat"}},
		},
		{
			// three pairs share key 9 and come out fully reversed
			name: "example 2",
			in: []Pair{
				{5, "apple"}, {9, "banana"}, {9, "cherry"}, {1, "date"}, {9, "elderberry"},
			},
			want: []Pair{
				{1, "date"}, {5, "apple"}, {9, "elderberry"}, {9, "cherry"}, {9, "banana"},
			},
		},
	}

	for _, c := range cases {
		got := QuickSort(slices.Clone(c.in))

		if !slices.Equal(got, c.want) {
			t.Errorf("%s: QuickSort(%s) = %s, want %s",
				c.name, formatPairs(c.in), formatPairs(got), formatPairs(c.want))
		}
	}
}

// Ordering by key across the shapes that stress the partition differently. Tie
// order is deliberately not asserted here -- only that keys are non-decreasing
// and no element went missing.
func TestQuickSortOrdersByKey(t *testing.T) {
	cases := []struct {
		name string
		in   []Pair
	}{
		// degenerate: constraints allow an empty list
		{"empty", []Pair{}},
		{"nil", nil},
		{"single", []Pair{{1, "a"}}},

		// two elements: the smallest input that partitions
		{"two sorted", []Pair{{1, "a"}, {2, "b"}}},
		{"two reversed", []Pair{{2, "a"}, {1, "b"}}},
		{"two equal", []Pair{{1, "a"}, {1, "b"}}},

		{"three sorted", []Pair{{1, "a"}, {2, "b"}, {3, "c"}}},
		{"three reversed", []Pair{{3, "a"}, {2, "b"}, {1, "c"}}},
		{"three equal", []Pair{{1, "a"}, {1, "b"}, {1, "c"}}},

		// ALREADY SORTED is quicksort's worst case with a right-most pivot: the
		// pivot is always the largest element, so every partition puts n-1
		// elements on the left and recursion depth is O(n)
		{"sorted eight", []Pair{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}, {5, "e"}, {6, "f"}, {7, "g"}, {8, "h"}}},

		// REVERSE SORTED is the mirror worst case: the pivot is always the
		// smallest, so left stays at start and everything lands on the right
		{"reversed eight", []Pair{{8, "a"}, {7, "b"}, {6, "c"}, {5, "d"}, {4, "e"}, {3, "f"}, {2, "g"}, {1, "h"}}},

		// ALL KEYS EQUAL: with a strict `<` the comparison never fires, `left`
		// never advances, and the pivot swaps to position `start` every time --
		// the most degenerate partition possible. An implementation that recurses
		// on the wrong sub-range hangs or overflows here rather than failing.
		{"all equal", []Pair{{5, "a"}, {5, "b"}, {5, "c"}, {5, "d"}, {5, "e"}, {5, "f"}}},

		// pivot is the smallest element: left sub-array is empty, so the first
		// recursive call gets an inverted range (start > end)
		{"pivot smallest", []Pair{{5, "a"}, {4, "b"}, {3, "c"}, {2, "d"}, {1, "e"}}},
		// pivot is the largest element: right sub-array is empty
		{"pivot largest", []Pair{{1, "a"}, {4, "b"}, {3, "c"}, {2, "d"}, {5, "e"}}},

		// negatives and zero
		{"negatives", []Pair{{0, "a"}, {-3, "b"}, {5, "c"}, {-3, "d"}, {0, "e"}}},
		{"all negative", []Pair{{-1, "a"}, {-5, "b"}, {-3, "c"}}},

		// duplicates mixed with distinct keys, so partitions split unevenly at
		// several depths
		{"mixed duplicates", []Pair{{2, "a"}, {1, "b"}, {2, "c"}, {3, "d"}, {1, "e"}, {3, "f"}, {2, "g"}}},

		// larger shuffled input
		{"shuffled twelve", []Pair{
			{7, "a"}, {2, "b"}, {11, "c"}, {4, "d"}, {9, "e"}, {1, "f"},
			{12, "g"}, {5, "h"}, {3, "i"}, {10, "j"}, {6, "k"}, {8, "l"},
		}},
	}

	for _, c := range cases {
		in := slices.Clone(c.in)

		got := QuickSort(slices.Clone(c.in))

		if !isSortedByKey(got) {
			t.Errorf("%s: QuickSort(%s) = %s, keys are not non-decreasing",
				c.name, formatPairs(in), formatPairs(got))
		}

		// the key sequence itself is fully determined even though tie order is not
		wantKeys := make([]int, len(in))
		for i, p := range stableSorted(in) {
			wantKeys[i] = p.Key
		}
		gotKeys := make([]int, len(got))
		for i, p := range got {
			gotKeys[i] = p.Key
		}
		if !slices.Equal(gotKeys, wantKeys) {
			t.Errorf("%s: key sequence = %v, want %v", c.name, gotKeys, wantKeys)
		}

		checkPermutation(t, in, got)
	}
}

// QuickSort partitions by swapping inside the caller's slice and returns that
// same slice. A solution that builds a new array and returns it would satisfy
// every ordering assertion above while quietly changing the contract -- and
// QuickSortHelper, which takes explicit start/end indices, only makes sense if
// the work happens in place.
func TestQuickSortIsInPlace(t *testing.T) {
	in := []Pair{{3, "a"}, {1, "b"}, {4, "c"}, {1, "d"}, {5, "e"}}
	original := slices.Clone(in)
	before := &in[0]

	got := QuickSort(in)

	if len(got) == 0 {
		t.Fatal("QuickSort returned an empty slice")
	}
	if &got[0] != before {
		t.Error("QuickSort returned a different backing array, expected an in-place partition")
	}
	if !isSortedByKey(in) {
		t.Errorf("the caller's slice was not sorted: %s", formatPairs(in))
	}
	checkPermutation(t, original, in)
}

// QuickSortHelper's contract is a HALF-OPEN-free, fully inclusive range
// [start, end]. Sorting a window must leave everything outside it untouched --
// which is what makes the recursive calls on [start, left-1] and [left+1, end]
// safe. Getting the boundaries wrong (off by one on either side) shows up here
// as a disturbed neighbour rather than as a wrong sort.
func TestQuickSortHelperSortsOnlyItsRange(t *testing.T) {
	cases := []struct {
		name       string
		start, end int
	}{
		{"middle window", 2, 5},
		{"prefix", 0, 3},
		{"suffix", 4, 7},
		{"whole range", 0, 7},
		{"single element", 3, 3},
		// end < start: an empty range, must be a no-op rather than a panic
		{"inverted range", 4, 3},
	}

	base := []Pair{
		{8, "a"}, {3, "b"}, {7, "c"}, {1, "d"}, {6, "e"}, {2, "f"}, {5, "g"}, {4, "h"},
	}

	for _, c := range cases {
		pairs := slices.Clone(base)

		QuickSortHelper(pairs, c.start, c.end)

		// everything before the window is byte-for-byte unchanged
		if c.start > 0 && !slices.Equal(pairs[:c.start], base[:c.start]) {
			t.Errorf("%s: QuickSortHelper(_, %d, %d) disturbed the prefix [:%d]: got %s, want %s",
				c.name, c.start, c.end, c.start, formatPairs(pairs[:c.start]), formatPairs(base[:c.start]))
		}
		// and everything after it
		if c.end+1 < len(base) && c.end >= c.start && !slices.Equal(pairs[c.end+1:], base[c.end+1:]) {
			t.Errorf("%s: QuickSortHelper(_, %d, %d) disturbed the suffix [%d:]: got %s, want %s",
				c.name, c.start, c.end, c.end+1, formatPairs(pairs[c.end+1:]), formatPairs(base[c.end+1:]))
		}

		if c.end <= c.start {
			// nothing to sort; the whole slice must be untouched
			if !slices.Equal(pairs, base) {
				t.Errorf("%s: QuickSortHelper(_, %d, %d) modified the slice on an empty/single range: got %s, want %s",
					c.name, c.start, c.end, formatPairs(pairs), formatPairs(base))
			}
			continue
		}

		window := pairs[c.start : c.end+1]
		if !isSortedByKey(window) {
			t.Errorf("%s: QuickSortHelper(_, %d, %d) left the window unsorted: %s",
				c.name, c.start, c.end, formatPairs(window))
		}
		checkPermutation(t, base[c.start:c.end+1], window)
	}
}

// Randomised sweep. Fixed seed, so any failure reproduces. The key range is kept
// narrow so ties are frequent: with a strict `<` comparison, runs of equal keys
// are where a partition most easily loses or duplicates an element.
func TestQuickSortRandomised(t *testing.T) {
	r := rand.New(rand.NewPCG(7, 99))

	for n := 0; n <= 40; n++ {
		for trial := range 20 {
			in := make([]Pair, n)
			for i := range in {
				in[i] = Pair{Key: r.IntN(6), Value: itoa(i)}
			}
			original := slices.Clone(in)

			got := QuickSort(slices.Clone(in))

			if !isSortedByKey(got) {
				t.Fatalf("n=%d trial=%d: not sorted\nin  %s\ngot %s",
					n, trial, formatPairs(original), formatPairs(got))
			}
			checkPermutation(t, original, got)
			if t.Failed() {
				t.Fatalf("n=%d trial=%d: input was %s", n, trial, formatPairs(original))
			}
		}
	}
}

// Constraints cap the input at 100 pairs. Sorted, reversed and all-equal inputs
// are the three shapes that drive a right-most-pivot quicksort to O(n) recursion
// depth -- the cases where a bad base case blows the stack instead of returning
// a wrong answer.
func TestQuickSortLong(t *testing.T) {
	n := 100

	shapes := []struct {
		name string
		key  func(i int) int
	}{
		{"already sorted", func(i int) int { return i }},
		{"reverse sorted", func(i int) int { return n - i }},
		{"all keys equal", func(int) int { return 7 }},
		{"two distinct keys", func(i int) int { return i % 2 }},
	}

	for _, s := range shapes {
		in := make([]Pair, n)
		for i := range in {
			in[i] = Pair{Key: s.key(i), Value: itoa(i)}
		}
		original := slices.Clone(in)

		got := QuickSort(in)

		if len(got) != n {
			t.Errorf("%s: got %d elements, want %d", s.name, len(got), n)
			continue
		}
		if !isSortedByKey(got) {
			t.Errorf("%s: result is not sorted by key", s.name)
		}
		checkPermutation(t, original, got)
	}
}
