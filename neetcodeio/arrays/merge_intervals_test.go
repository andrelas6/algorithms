package arrays

import (
	"math/rand"
	"slices"
	"testing"
)

func intervalsEqual(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !slices.Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

func cloneIntervals(in [][]int) [][]int {
	out := make([][]int, len(in))
	for i, iv := range in {
		out[i] = slices.Clone(iv)
	}
	return out
}

func TestMergeExamples(t *testing.T) {
	cases := []struct {
		name string
		in   [][]int
		want [][]int
	}{
		// samples from the statement
		{"sample 1", [][]int{{1, 3}, {1, 5}, {6, 7}}, [][]int{{1, 5}, {6, 7}}},
		{"sample 2", [][]int{{1, 2}, {2, 3}}, [][]int{{1, 3}}},
		{"leetcode sample", [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}, [][]int{{1, 6}, {8, 10}, {15, 18}}},

		// degenerate
		{"single interval", [][]int{{5, 10}}, [][]int{{5, 10}}},
		{"single point", [][]int{{4, 4}}, [][]int{{4, 4}}},

		// touching endpoints merge; adjacent-but-separate do not
		{"touching merges", [][]int{{1, 2}, {2, 3}, {3, 4}}, [][]int{{1, 4}}},
		{"adjacent stays apart", [][]int{{1, 2}, {3, 4}}, [][]int{{1, 2}, {3, 4}}},

		// nothing overlaps
		{"all disjoint", [][]int{{1, 2}, {4, 5}, {7, 8}}, [][]int{{1, 2}, {4, 5}, {7, 8}}},
		// everything overlaps
		{"all one blob", [][]int{{1, 10}, {2, 3}, {4, 9}}, [][]int{{1, 10}}},

		// input arrives out of order
		{"unsorted", [][]int{{5, 10}, {1, 3}, {2, 6}}, [][]int{{1, 10}}},
		{"reverse sorted", [][]int{{8, 9}, {5, 6}, {1, 2}}, [][]int{{1, 2}, {5, 6}, {8, 9}}},

		// one interval swallows a later one entirely - end must not shrink
		{"contained", [][]int{{1, 10}, {2, 3}}, [][]int{{1, 10}}},
		{"contained unsorted", [][]int{{2, 3}, {1, 10}}, [][]int{{1, 10}}},

		// equal starts
		{"same start", [][]int{{1, 4}, {1, 2}, {1, 7}}, [][]int{{1, 7}}},
		// zero-width intervals
		{"points on a line", [][]int{{1, 1}, {2, 2}, {3, 3}}, [][]int{{1, 1}, {2, 2}, {3, 3}}},
		{"point inside an interval", [][]int{{1, 5}, {3, 3}}, [][]int{{1, 5}}},

		// duplicates
		{"exact duplicates", [][]int{{1, 3}, {1, 3}, {1, 3}}, [][]int{{1, 3}}},

		// constraint bounds
		{"full range", [][]int{{0, 1000}, {500, 600}}, [][]int{{0, 1000}}},
		{"at the edges", [][]int{{0, 0}, {1000, 1000}}, [][]int{{0, 0}, {1000, 1000}}},
	}

	for _, c := range cases {
		in := cloneIntervals(c.in)

		got := merge(in)

		if !intervalsEqual(got, c.want) {
			t.Errorf("%s: merge(%v) = %v, want %v", c.name, c.in, got, c.want)
		}
	}
}

// Independent oracle. Coordinates are doubled so that "touching" (shared
// endpoint) and "adjacent" (2 then 3) stop looking alike on an integer grid:
// [1,2] and [2,3] become [2,4] and [4,6] and share point 4, while [1,2] and
// [3,4] become [2,4] and [6,8] with a gap at 5. Mark every covered point, then
// read off the maximal runs.
func mergeByCoverage(intervals [][]int) [][]int {
	const limit = 2*1000 + 2

	covered := make([]bool, limit)
	for _, iv := range intervals {
		for p := 2 * iv[0]; p <= 2*iv[1]; p++ {
			covered[p] = true
		}
	}

	var out [][]int
	p := 0
	for p < limit {
		if !covered[p] {
			p++
			continue
		}
		start := p
		for p < limit && covered[p] {
			p++
		}
		out = append(out, []int{start / 2, (p - 1) / 2})
	}

	return out
}

func TestMergeAgainstCoverage(t *testing.T) {
	rng := rand.New(rand.NewSource(7))

	for trial := 0; trial < 2000; trial++ {
		n := 1 + rng.Intn(12)
		in := make([][]int, n)
		for i := range in {
			// small coordinate space so overlaps actually happen
			a := rng.Intn(20)
			b := a + rng.Intn(6)
			in[i] = []int{a, b}
		}

		want := mergeByCoverage(cloneIntervals(in))
		got := merge(cloneIntervals(in))

		if !intervalsEqual(got, want) {
			t.Fatalf("merge(%v) = %v, want %v", in, got, want)
		}
	}
}

// Structural properties of any correct answer.
func TestMergeOutputIsSortedAndDisjoint(t *testing.T) {
	rng := rand.New(rand.NewSource(8))

	for trial := 0; trial < 1000; trial++ {
		n := 1 + rng.Intn(15)
		in := make([][]int, n)
		for i := range in {
			a := rng.Intn(50)
			b := a + rng.Intn(10)
			in[i] = []int{a, b}
		}

		got := merge(cloneIntervals(in))

		for i, iv := range got {
			if len(iv) != 2 {
				t.Fatalf("interval %d is %v, want two values", i, iv)
			}
			if iv[0] > iv[1] {
				t.Fatalf("interval %d is %v: start after end", i, iv)
			}
			if i > 0 && got[i-1][1] >= iv[0] {
				t.Fatalf("merge(%v) = %v: %v and %v still overlap or touch", in, got, got[i-1], iv)
			}
		}
	}
}

// Every input interval must be fully inside exactly one output interval.
func TestMergeCoversEveryInput(t *testing.T) {
	rng := rand.New(rand.NewSource(9))

	for trial := 0; trial < 1000; trial++ {
		n := 1 + rng.Intn(10)
		in := make([][]int, n)
		for i := range in {
			a := rng.Intn(30)
			b := a + rng.Intn(8)
			in[i] = []int{a, b}
		}

		reference := cloneIntervals(in)
		got := merge(cloneIntervals(in))

		for _, iv := range reference {
			found := false
			for _, m := range got {
				if m[0] <= iv[0] && iv[1] <= m[1] {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("merge(%v) = %v: input %v is not covered by any output interval", reference, got, iv)
			}
		}
	}
}

// merge sorts in place, so the caller's ORDER changes - that is expected and
// worth being explicit about. The VALUES inside each interval should survive,
// though: nothing about merging says the caller's numbers may be rewritten.
func TestMergeDoesNotRewriteCallerValues(t *testing.T) {
	in := [][]int{{1, 3}, {2, 6}, {8, 10}}

	// what the caller holds onto
	first := in[0]
	before := slices.Clone(first)

	merge(in)

	if !slices.Equal(first, before) {
		t.Errorf("merge rewrote the caller's interval in place: %v became %v", before, first)
	}
}
