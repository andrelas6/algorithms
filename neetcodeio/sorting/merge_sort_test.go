package sorting

import (
	"math/rand/v2"
	"slices"
	"testing"
)

// byKey, stableSorted and itoa live in insertion_sort_test.go -- same package,
// so they are reused here rather than duplicated. stableSorted is the oracle:
// stdlib SortStableFunc on Key, which is exactly the spec (sort by key, equal
// keys keep input order).

func formatPairs(pairs []Pair) string {
	if pairs == nil {
		return "nil"
	}

	out := "["
	for i, p := range pairs {
		if i > 0 {
			out += " "
		}
		out += "(" + itoa(p.Key) + ", " + p.Value + ")"
	}
	return out + "]"
}

func TestMergeSort(t *testing.T) {
	cases := []struct {
		name string
		in   []Pair
		want []Pair
	}{
		{
			name: "example 1",
			in: []Pair{
				{5, "apple"}, {2, "banana"}, {9, "cherry"}, {1, "date"}, {9, "elderberry"},
			},
			want: []Pair{
				{1, "date"}, {2, "banana"}, {5, "apple"}, {9, "cherry"}, {9, "elderberry"},
			},
		},
		{
			name: "example 2",
			in:   []Pair{{3, "cat"}, {2, "dog"}, {3, "bird"}},
			want: []Pair{{2, "dog"}, {3, "cat"}, {3, "bird"}},
		},

		// degenerate: constraints allow an empty list
		{name: "empty", in: []Pair{}, want: []Pair{}},
		{name: "nil", in: nil, want: nil},

		// the recursion base case
		{name: "single", in: []Pair{{1, "a"}}, want: []Pair{{1, "a"}}},

		// two elements: the smallest input that actually splits and merges
		{name: "two sorted", in: []Pair{{1, "a"}, {2, "b"}}, want: []Pair{{1, "a"}, {2, "b"}}},
		{name: "two reversed", in: []Pair{{2, "a"}, {1, "b"}}, want: []Pair{{1, "b"}, {2, "a"}}},
		// equal keys at length 2: the tightest possible stability check
		{name: "two equal", in: []Pair{{1, "a"}, {1, "b"}}, want: []Pair{{1, "a"}, {1, "b"}}},

		// three elements: odd length, so the split is uneven (1 left, 2 right)
		{name: "three sorted", in: []Pair{{1, "a"}, {2, "b"}, {3, "c"}}, want: []Pair{{1, "a"}, {2, "b"}, {3, "c"}}},
		{name: "three reversed", in: []Pair{{3, "a"}, {2, "b"}, {1, "c"}}, want: []Pair{{1, "c"}, {2, "b"}, {3, "a"}}},
		{name: "three all equal", in: []Pair{{1, "a"}, {1, "b"}, {1, "c"}}, want: []Pair{{1, "a"}, {1, "b"}, {1, "c"}}},

		// already sorted and fully reversed at a larger size
		{
			name: "sorted six",
			in:   []Pair{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}, {5, "e"}, {6, "f"}},
			want: []Pair{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}, {5, "e"}, {6, "f"}},
		},
		{
			name: "reversed six",
			in:   []Pair{{6, "a"}, {5, "b"}, {4, "c"}, {3, "d"}, {2, "e"}, {1, "f"}},
			want: []Pair{{1, "f"}, {2, "e"}, {3, "d"}, {4, "c"}, {5, "b"}, {6, "a"}},
		},

		// one half entirely below the other: the merge drains `left` first and
		// then falls into the leftover-`right` loop without ever interleaving
		{
			name: "left half all smaller",
			in:   []Pair{{1, "a"}, {2, "b"}, {3, "c"}, {7, "d"}, {8, "e"}, {9, "f"}},
			want: []Pair{{1, "a"}, {2, "b"}, {3, "c"}, {7, "d"}, {8, "e"}, {9, "f"}},
		},
		// and the mirror image: `right` drains first, exercising the other
		// leftover loop
		{
			name: "right half all smaller",
			in:   []Pair{{7, "a"}, {8, "b"}, {9, "c"}, {1, "d"}, {2, "e"}, {3, "f"}},
			want: []Pair{{1, "d"}, {2, "e"}, {3, "f"}, {7, "a"}, {8, "b"}, {9, "c"}},
		},

		// perfectly interleaved halves, so the merge alternates every step
		{
			name: "interleaved",
			in:   []Pair{{1, "a"}, {3, "b"}, {5, "c"}, {2, "d"}, {4, "e"}, {6, "f"}},
			want: []Pair{{1, "a"}, {2, "d"}, {3, "b"}, {4, "e"}, {5, "c"}, {6, "f"}},
		},

		// negatives and zero
		{
			name: "negatives",
			in:   []Pair{{0, "a"}, {-3, "b"}, {5, "c"}, {-3, "d"}, {0, "e"}},
			want: []Pair{{-3, "b"}, {-3, "d"}, {0, "a"}, {0, "e"}, {5, "c"}},
		},

		// power-of-two length: every split is even all the way down
		{
			name: "eight elements",
			in:   []Pair{{8, "a"}, {3, "b"}, {5, "c"}, {1, "d"}, {7, "e"}, {2, "f"}, {6, "g"}, {4, "h"}},
			want: []Pair{{1, "d"}, {2, "f"}, {3, "b"}, {4, "h"}, {5, "c"}, {6, "g"}, {7, "e"}, {8, "a"}},
		},

		// odd length that splits unevenly at several levels (7 -> 3 + 4 -> ...)
		{
			name: "seven elements",
			in:   []Pair{{7, "a"}, {1, "b"}, {6, "c"}, {2, "d"}, {5, "e"}, {3, "f"}, {4, "g"}},
			want: []Pair{{1, "b"}, {2, "d"}, {3, "f"}, {4, "g"}, {5, "e"}, {6, "c"}, {7, "a"}},
		},
	}

	for _, c := range cases {
		got := mergeSort(slices.Clone(c.in))

		if !slices.Equal(got, c.want) {
			t.Errorf("%s: mergeSort(%s) = %s, want %s",
				c.name, formatPairs(c.in), formatPairs(got), formatPairs(c.want))
		}
	}
}

// Stability is stated explicitly in the problem and is the property most easily
// lost: it hinges entirely on `<=` rather than `<` in the merge comparison. With
// `<`, a tie hands the slot to `right`, which jumps a later element ahead of an
// earlier one with the same key.
//
// Every value here is unique and records its input position, so the expected
// output is fully determined and any reordering of equal keys is visible.
func TestMergeSortStability(t *testing.T) {
	// keys chosen so that ties land in different halves at several recursion
	// depths, not just at the final merge
	keys := []int{2, 1, 2, 1, 3, 2, 1, 3, 2, 1, 3, 3}

	in := make([]Pair, len(keys))
	for i, k := range keys {
		in[i] = Pair{Key: k, Value: itoa(i)}
	}

	got := mergeSort(slices.Clone(in))
	want := stableSorted(in)

	if !slices.Equal(got, want) {
		t.Fatalf("mergeSort(%s)\n = %s\nwant %s (equal keys must keep their input order)",
			formatPairs(in), formatPairs(got), formatPairs(want))
	}

	// Say it a second way, so a failure names the actual violation rather than
	// just showing two slices: within each run of equal keys, the original
	// indices must be strictly increasing.
	index := make(map[Pair]int, len(in))
	for i, p := range in {
		index[p] = i
	}
	for i := 1; i < len(got); i++ {
		if got[i-1].Key != got[i].Key {
			continue
		}
		if index[got[i-1]] > index[got[i]] {
			t.Errorf("key %d: output position %d holds input element %d, but position %d holds input element %d -- equal keys were reordered",
				got[i].Key, i-1, index[got[i-1]], i, index[got[i]])
		}
	}
}

// The merge writes into a fresh slice using three indices (i, j, k). An
// off-by-one in k, or dropping one of the two leftover-drain loops, shows up as
// a missing or duplicated element. Comparing the output as a multiset catches
// that independently of ordering.
func TestMergeSortIsAPermutation(t *testing.T) {
	cases := [][]Pair{
		{{3, "a"}, {1, "b"}, {2, "c"}},
		{{5, "a"}, {5, "b"}, {5, "c"}, {5, "d"}, {5, "e"}},
		{{9, "a"}, {8, "b"}, {7, "c"}, {6, "d"}, {5, "e"}, {4, "f"}, {3, "g"}},
		{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}, {5, "e"}, {6, "f"}, {7, "g"}, {8, "h"}, {9, "i"}},
	}

	for _, in := range cases {
		got := mergeSort(slices.Clone(in))

		if len(got) != len(in) {
			t.Errorf("mergeSort(%s) returned %d elements, want %d",
				formatPairs(in), len(got), len(in))
			continue
		}

		gotSorted := slices.Clone(got)
		inSorted := slices.Clone(in)
		slices.SortStableFunc(gotSorted, byKey)
		slices.SortStableFunc(inSorted, byKey)

		counts := make(map[Pair]int, len(in))
		for _, p := range in {
			counts[p]++
		}
		for _, p := range got {
			counts[p]--
		}
		for p, n := range counts {
			if n > 0 {
				t.Errorf("mergeSort(%s) dropped %v (%d times)", formatPairs(in), p, n)
			}
			if n < 0 {
				t.Errorf("mergeSort(%s) duplicated %v (%d times)", formatPairs(in), p, -n)
			}
		}
	}
}

// merge() only produces a sorted result if BOTH inputs are already sorted --
// that is the contract the recursion is responsible for upholding. Exercise it
// directly on hand-built halves, including the lopsided and empty ones the
// recursion generates at odd lengths.
func TestMerge(t *testing.T) {
	cases := []struct {
		name        string
		left, right []Pair
		want        []Pair
	}{
		{"both empty", []Pair{}, []Pair{}, []Pair{}},
		{"left empty", []Pair{}, []Pair{{1, "a"}, {2, "b"}}, []Pair{{1, "a"}, {2, "b"}}},
		{"right empty", []Pair{{1, "a"}, {2, "b"}}, []Pair{}, []Pair{{1, "a"}, {2, "b"}}},

		{"singletons in order", []Pair{{1, "a"}}, []Pair{{2, "b"}}, []Pair{{1, "a"}, {2, "b"}}},
		{"singletons out of order", []Pair{{2, "a"}}, []Pair{{1, "b"}}, []Pair{{1, "b"}, {2, "a"}}},
		// the tie: left must win, or stability is gone
		{"singleton tie", []Pair{{1, "a"}}, []Pair{{1, "b"}}, []Pair{{1, "a"}, {1, "b"}}},

		{
			name:  "interleaving",
			left:  []Pair{{1, "a"}, {3, "b"}, {5, "c"}},
			right: []Pair{{2, "d"}, {4, "e"}, {6, "f"}},
			want:  []Pair{{1, "a"}, {2, "d"}, {3, "b"}, {4, "e"}, {5, "c"}, {6, "f"}},
		},
		{
			// left drains completely first, then the leftover-right loop runs
			name:  "left exhausted first",
			left:  []Pair{{1, "a"}, {2, "b"}},
			right: []Pair{{3, "c"}, {4, "d"}, {5, "e"}},
			want:  []Pair{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}, {5, "e"}},
		},
		{
			// mirror image: the leftover-left loop runs
			name:  "right exhausted first",
			left:  []Pair{{3, "a"}, {4, "b"}, {5, "c"}},
			right: []Pair{{1, "d"}, {2, "e"}},
			want:  []Pair{{1, "d"}, {2, "e"}, {3, "a"}, {4, "b"}, {5, "c"}},
		},
		{
			// every element ties: the whole of left must precede the whole of right
			name:  "all keys equal",
			left:  []Pair{{1, "a"}, {1, "b"}},
			right: []Pair{{1, "c"}, {1, "d"}},
			want:  []Pair{{1, "a"}, {1, "b"}, {1, "c"}, {1, "d"}},
		},
		{
			// ties spread through both halves
			name:  "repeated ties",
			left:  []Pair{{1, "a"}, {2, "b"}, {3, "c"}},
			right: []Pair{{1, "d"}, {2, "e"}, {3, "f"}},
			want:  []Pair{{1, "a"}, {1, "d"}, {2, "b"}, {2, "e"}, {3, "c"}, {3, "f"}},
		},
		{
			name:  "lopsided",
			left:  []Pair{{5, "a"}},
			right: []Pair{{1, "b"}, {2, "c"}, {3, "d"}, {4, "e"}, {6, "f"}},
			want:  []Pair{{1, "b"}, {2, "c"}, {3, "d"}, {4, "e"}, {5, "a"}, {6, "f"}},
		},
	}

	for _, c := range cases {
		got := merge(slices.Clone(c.left), slices.Clone(c.right))

		if !slices.Equal(got, c.want) {
			t.Errorf("%s: merge(%s, %s) = %s, want %s",
				c.name, formatPairs(c.left), formatPairs(c.right), formatPairs(got), formatPairs(c.want))
		}
	}
}

// Randomised sweep against the stdlib oracle. Fixed seed, so a failure is
// reproducible. Covers every length from 0 to 40 with a key range narrow enough
// that ties are frequent -- which is where stability bugs actually live.
func TestMergeSortAgainstStdlib(t *testing.T) {
	r := rand.New(rand.NewPCG(42, 1024))

	for n := 0; n <= 40; n++ {
		for trial := 0; trial < 20; trial++ {
			in := make([]Pair, n)
			for i := range in {
				// keys collide often on purpose
				in[i] = Pair{Key: r.IntN(6), Value: itoa(i)}
			}

			got := mergeSort(slices.Clone(in))
			want := stableSorted(in)

			if !slices.Equal(got, want) {
				t.Fatalf("n=%d trial=%d\nin   %s\ngot  %s\nwant %s",
					n, trial, formatPairs(in), formatPairs(got), formatPairs(want))
			}
		}
	}
}

// Constraints cap the input at 100 pairs. Two shapes that stress the recursion
// differently: fully reversed (every merge interleaves nothing, one half always
// drains first) and all-equal keys (every comparison is a tie, so stability is
// under maximum pressure).
func TestMergeSortLong(t *testing.T) {
	n := 100

	reversed := make([]Pair, n)
	for i := range reversed {
		reversed[i] = Pair{Key: n - i, Value: itoa(i)}
	}
	if got, want := mergeSort(slices.Clone(reversed)), stableSorted(reversed); !slices.Equal(got, want) {
		t.Errorf("reversed input of %d: got %s, want %s", n, formatPairs(got), formatPairs(want))
	}

	allEqual := make([]Pair, n)
	for i := range allEqual {
		allEqual[i] = Pair{Key: 7, Value: itoa(i)}
	}
	got := mergeSort(slices.Clone(allEqual))
	if !slices.Equal(got, allEqual) {
		t.Errorf("all keys equal: the input order must be preserved exactly\ngot  %s\nwant %s",
			formatPairs(got), formatPairs(allEqual))
	}
}
