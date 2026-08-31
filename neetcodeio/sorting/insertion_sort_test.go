package sorting

import (
	"cmp"
	"slices"
	"testing"
)

// byKey is the ordering the problem specifies: compare on Key only, and use a
// STABLE sort so equal keys keep their input order. Used as an oracle below --
// it is the stdlib doing the sorting, not a second copy of insertion sort, so
// the tests stay independent of the implementation under test.
func byKey(a, b Pair) int { return cmp.Compare(a.Key, b.Key) }

func stableSorted(pairs []Pair) []Pair {
	out := slices.Clone(pairs)
	slices.SortStableFunc(out, byKey)
	return out
}

func formatStates(states [][]Pair) string {
	if states == nil {
		return "nil"
	}

	out := "["
	for i, s := range states {
		if i > 0 {
			out += "\n "
		}
		out += "["
		for j, p := range s {
			if j > 0 {
				out += " "
			}
			out += "(" + itoa(p.Key) + ", " + p.Value + ")"
		}
		out += "]"
	}
	return out + "]"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	digits := ""
	for n > 0 {
		digits = string(rune('0'+n%10)) + digits
		n /= 10
	}
	if neg {
		return "-" + digits
	}
	return digits
}

// Both worked examples from the statement, compared exactly. These pin down the
// two things prose leaves ambiguous: the FIRST state is the array after
// "inserting" element 0 (i.e. the untouched input), and there are exactly
// len(pairs) states -- one per insertion, not one per swap.
func TestInsertionSortExamples(t *testing.T) {
	cases := []struct {
		name string
		in   []Pair
		want [][]Pair
	}{
		{
			name: "example 1",
			in: []Pair{
				{5, "apple"}, {2, "banana"}, {9, "cherry"},
			},
			want: [][]Pair{
				{{5, "apple"}, {2, "banana"}, {9, "cherry"}},
				{{2, "banana"}, {5, "apple"}, {9, "cherry"}},
				{{2, "banana"}, {5, "apple"}, {9, "cherry"}},
			},
		},
		{
			// the stability example: "cat" was before "bird" in the input and
			// must still be before it at the end, even though both have key 3
			name: "example 2",
			in: []Pair{
				{3, "cat"}, {3, "bird"}, {2, "dog"},
			},
			want: [][]Pair{
				{{3, "cat"}, {3, "bird"}, {2, "dog"}},
				{{3, "cat"}, {3, "bird"}, {2, "dog"}},
				{{2, "dog"}, {3, "cat"}, {3, "bird"}},
			},
		},
	}

	for _, c := range cases {
		got := insertionSort(slices.Clone(c.in))

		if len(got) != len(c.want) {
			t.Errorf("%s: got %d states, want %d\ngot:  %s\nwant: %s",
				c.name, len(got), len(c.want), formatStates(got), formatStates(c.want))
			continue
		}
		for i := range c.want {
			if !slices.Equal(got[i], c.want[i]) {
				t.Errorf("%s: state %d = %v, want %v\ngot:  %s\nwant: %s",
					c.name, i, got[i], c.want[i], formatStates(got), formatStates(c.want))
			}
		}
	}
}

// "There should be pairs.length states in total." One state per element, no
// matter how many swaps each insertion needed. An implementation that snapshots
// per comparison or per swap produces a variable number of states and fails
// here on the first input that needs more than one shift.
func TestInsertionSortStateCount(t *testing.T) {
	cases := [][]Pair{
		{},
		{{1, "a"}},
		{{2, "a"}, {1, "b"}},
		{{3, "a"}, {2, "b"}, {1, "c"}},
		// reverse sorted: the worst case, every insertion shifts all the way down
		{{5, "a"}, {4, "b"}, {3, "c"}, {2, "d"}, {1, "e"}},
		// already sorted: the best case, no insertion shifts at all
		{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}, {5, "e"}},
		// all equal: stability means no shifts either
		{{1, "a"}, {1, "b"}, {1, "c"}, {1, "d"}},
	}

	for _, in := range cases {
		got := insertionSort(slices.Clone(in))

		if len(got) != len(in) {
			t.Errorf("insertionSort(%v): got %d states, want %d (one per element)\ngot: %s",
				in, len(got), len(in), formatStates(got))
		}
	}
}

// The empty list has zero insertions, so it has zero states. Called out on its
// own because the natural way to write this ("seed the result with the input,
// then loop") emits one state for an empty input and is off by one everywhere.
func TestInsertionSortEmpty(t *testing.T) {
	got := insertionSort([]Pair{})

	if len(got) != 0 {
		t.Errorf("insertionSort([]) = %s, want 0 states", formatStates(got))
	}

	if got := insertionSort(nil); len(got) != 0 {
		t.Errorf("insertionSort(nil) = %s, want 0 states", formatStates(got))
	}
}

// The defining invariant of insertion sort, stated without reimplementing it:
// after insertion i, the first i+1 elements are the first i+1 INPUT elements in
// stable-sorted order, and everything from i+1 onwards is still untouched input.
// Checking every state this way catches an implementation that reaches the right
// final answer through the wrong intermediate states.
func TestInsertionSortStateInvariant(t *testing.T) {
	cases := [][]Pair{
		{{5, "apple"}, {2, "banana"}, {9, "cherry"}},
		{{3, "cat"}, {3, "bird"}, {2, "dog"}},
		{{4, "a"}, {3, "b"}, {2, "c"}, {1, "d"}},
		{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}},
		{{2, "a"}, {1, "b"}, {4, "c"}, {3, "d"}, {6, "e"}, {5, "f"}},
		// duplicates scattered through, so stability is exercised at every step
		{{2, "a"}, {1, "b"}, {2, "c"}, {1, "d"}, {2, "e"}},
		// negatives and zero
		{{0, "a"}, {-3, "b"}, {5, "c"}, {-3, "d"}, {0, "e"}},
		// one element out of place at the very end: the last insertion does all
		// the work, every earlier state is the untouched input
		{{1, "a"}, {2, "b"}, {3, "c"}, {0, "d"}},
		// one element out of place at the front
		{{9, "a"}, {1, "b"}, {2, "c"}, {3, "d"}},
	}

	for _, in := range cases {
		original := slices.Clone(in)

		states := insertionSort(slices.Clone(in))

		if len(states) != len(original) {
			t.Errorf("insertionSort(%v): got %d states, want %d", original, len(states), len(original))
			continue
		}

		for i, state := range states {
			if len(state) != len(original) {
				t.Errorf("insertionSort(%v): state %d has length %d, want %d (states are snapshots of the whole array)",
					original, i, len(state), len(original))
				continue
			}

			wantPrefix := stableSorted(original[:i+1])
			if !slices.Equal(state[:i+1], wantPrefix) {
				t.Errorf("insertionSort(%v): state %d prefix [:%d] = %v, want %v (the first %d elements must be sorted)",
					original, i, i+1, state[:i+1], wantPrefix, i+1)
			}

			wantSuffix := original[i+1:]
			if !slices.Equal(state[i+1:], wantSuffix) {
				t.Errorf("insertionSort(%v): state %d suffix [%d:] = %v, want %v (the unsorted tail must be untouched)",
					original, i, i+1, state[i+1:], wantSuffix)
			}
		}
	}
}

// Stability is stated explicitly in the problem, and it is the one property a
// naive swap-based implementation quietly breaks: using `>` instead of `>=` when
// deciding whether to keep shifting swaps equal keys past each other. Values are
// unique per key here so the original order is recoverable from the output.
func TestInsertionSortStability(t *testing.T) {
	in := []Pair{
		{2, "a"}, {1, "b"}, {2, "c"}, {1, "d"}, {2, "e"}, {3, "f"}, {1, "g"},
	}

	states := insertionSort(slices.Clone(in))

	if len(states) == 0 {
		t.Fatal("insertionSort returned no states")
	}
	final := states[len(states)-1]

	want := stableSorted(in)
	if !slices.Equal(final, want) {
		t.Errorf("final state = %v, want %v (equal keys must keep their input order)", final, want)
	}
}

// Each state is a separate snapshot in time. If the implementation appends the
// live slice instead of a copy, every entry in the result aliases the same
// backing array, so all of them show the FINAL array and the history is lost.
// The example tables would catch that too, but only as a confusing value
// mismatch -- this names the actual cause.
func TestInsertionSortStatesAreIndependentSnapshots(t *testing.T) {
	in := []Pair{{3, "a"}, {1, "b"}, {2, "c"}}

	states := insertionSort(slices.Clone(in))

	if len(states) < 2 {
		t.Fatalf("got %d states, want %d", len(states), len(in))
	}

	for i := range states {
		for j := i + 1; j < len(states); j++ {
			if len(states[i]) == 0 || len(states[j]) == 0 {
				continue
			}
			if &states[i][0] == &states[j][0] {
				t.Fatalf("states %d and %d share a backing array: the result holds one slice aliased %d times, not %d independent snapshots",
					i, j, len(states), len(states))
			}
		}
	}

	// mutating one state must not disturb any other
	before := slices.Clone(states[len(states)-1])
	states[0][0] = Pair{999, "mutated"}
	if !slices.Equal(states[len(states)-1], before) {
		t.Error("mutating state 0 changed the final state: the snapshots are not independent")
	}
}

// Constraints cap the input at 100 pairs. Reverse-sorted is the worst case for
// insertion sort (every element shifts the full distance), which is where a
// snapshot-per-swap implementation blows up the state count most visibly.
func TestInsertionSortLong(t *testing.T) {
	n := 100

	in := make([]Pair, n)
	for i := range in {
		in[i] = Pair{Key: n - i, Value: itoa(n - i)}
	}
	original := slices.Clone(in)

	states := insertionSort(in)

	if len(states) != n {
		t.Fatalf("got %d states, want %d", len(states), n)
	}

	final := states[n-1]
	want := stableSorted(original)
	if !slices.Equal(final, want) {
		t.Errorf("final state is not sorted:\ngot  %v\nwant %v", final, want)
	}
}
