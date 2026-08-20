package greedybyend

import (
	"math/bits"
	"math/rand"
	"testing"
)

func TestMaximizeNonOverlappingMeetings(t *testing.T) {
	cases := []struct {
		name     string
		meetings [][]int32
		want     int32
	}{
		// examples from the statement
		{"example 1", [][]int32{{1, 2}, {2, 3}, {3, 4}, {1, 3}}, 3},
		{"example 2", [][]int32{{0, 5}, {0, 1}, {1, 2}, {2, 3}, {3, 5}, {4, 6}}, 4},
		{"input format sample", [][]int32{{1, 2}, {3, 4}, {0, 6}, {5, 7}, {8, 9}, {5, 9}}, 4},

		// samples
		{"single meeting", [][]int32{{5, 10}}, 1},
		{"three touching", [][]int32{{1, 2}, {2, 3}, {3, 4}}, 3},

		// degenerate
		{"empty", [][]int32{}, 0},

		// touching is NOT overlapping
		{"two touching", [][]int32{{1, 2}, {2, 3}}, 2},
		{"two overlapping", [][]int32{{1, 3}, {2, 4}}, 1},

		// everything overlaps everything
		{"all overlap", [][]int32{{0, 10}, {1, 9}, {2, 8}}, 1},
		{"identical meetings", [][]int32{{1, 5}, {1, 5}, {1, 5}}, 1},

		// nothing overlaps
		{"all disjoint", [][]int32{{0, 1}, {2, 3}, {4, 5}, {6, 7}}, 4},

		// a long meeting that must be skipped in favour of two short ones
		{"skip the long one", [][]int32{{0, 5}, {1, 2}, {3, 4}}, 2},
		{"shortest-first would fail", [][]int32{{0, 10}, {9, 11}, {10, 20}}, 2},

		// already sorted by end vs reverse order
		{"reverse order input", [][]int32{{6, 7}, {4, 5}, {2, 3}, {0, 1}}, 4},

		// large values, per the constraints
		{"large values", [][]int32{{0, 1000000000}, {999999999, 1000000000}}, 1},
	}

	for _, c := range cases {
		if got := maximizeNonOverlappingMeetings(c.meetings); got != c.want {
			t.Errorf("%s: maximizeNonOverlappingMeetings(%v) = %d, want %d", c.name, c.meetings, got, c.want)
		}
	}
}

// Exhaustive oracle: try every subset, keep the largest pairwise-compatible one.
// Only usable for tiny n, which is the point.
func bruteForce(meetings [][]int32) int32 {
	n := len(meetings)
	best := 0
	for mask := 0; mask < 1<<n; mask++ {
		ok := true
		for i := 0; i < n && ok; i++ {
			if mask&(1<<i) == 0 {
				continue
			}
			for j := i + 1; j < n && ok; j++ {
				if mask&(1<<j) == 0 {
					continue
				}
				a, b := meetings[i], meetings[j]
				// disjoint iff one ends at or before the other starts
				if !(a[1] <= b[0] || b[1] <= a[0]) {
					ok = false
				}
			}
		}
		if ok && bits.OnesCount(uint(mask)) > best {
			best = bits.OnesCount(uint(mask))
		}
	}
	return int32(best)
}

func TestAgainstBruteForce(t *testing.T) {
	r := rand.New(rand.NewSource(7))

	for trial := 0; trial < 3000; trial++ {
		n := r.Intn(9) // 0..8 meetings
		meetings := make([][]int32, n)
		for i := range meetings {
			start := int32(r.Intn(10))
			end := start + 1 + int32(r.Intn(5))
			meetings[i] = []int32{start, end}
		}

		want := bruteForce(meetings)
		if got := maximizeNonOverlappingMeetings(meetings); got != want {
			t.Fatalf("maximizeNonOverlappingMeetings(%v) = %d, want %d", meetings, got, want)
		}
	}
}
