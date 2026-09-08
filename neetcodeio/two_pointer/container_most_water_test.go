package twopointer

import (
	"math/rand"
	"testing"
)

func TestMaxArea(t *testing.T) {
	cases := []struct {
		name string
		in   []int
		want int
	}{
		// samples from the statement
		{"neetcode sample 1", []int{1, 7, 2, 5, 4, 7, 3, 6}, 36},
		{"neetcode sample 2", []int{2, 2, 2}, 4},
		{"leetcode sample", []int{1, 8, 6, 2, 5, 4, 8, 3, 7}, 49},

		// degenerate: smallest legal input
		{"two bars", []int{1, 1}, 1},
		{"two bars unequal", []int{1, 1000}, 1},

		// zeros are legal heights
		{"all zeros", []int{0, 0, 0, 0}, 0},
		{"zeros at the ends", []int{0, 5, 5, 0}, 5},
		{"single tall bar is useless alone", []int{0, 0, 10, 0, 0}, 0},

		// widest pair is NOT the answer
		{"narrow beats wide", []int{1, 9, 9, 1}, 9},

		// widest pair IS the answer
		{"wide wins", []int{9, 1, 1, 9}, 27},

		// monotonic: pointer must walk all the way in
		{"increasing", []int{1, 2, 3, 4, 5}, 6},
		{"decreasing", []int{5, 4, 3, 2, 1}, 6},

		// ties everywhere - the case André asked about
		{"all equal", []int{3, 3, 3, 3, 3, 3}, 15},
		{"tie at the ends, taller inside", []int{4, 9, 9, 9, 4}, 18},
		{"tie at the ends, nothing inside", []int{4, 1, 1, 1, 4}, 16},
		{"tie ends beat everything", []int{6, 1, 2, 1, 6}, 24},

		// plateau of equal maxima
		{"equal maxima far apart", []int{7, 1, 1, 1, 1, 1, 7}, 42},

		// answer is an adjacent pair
		{"adjacent pair wins", []int{1, 100, 100, 1}, 100},
	}

	for _, c := range cases {
		in := append([]int(nil), c.in...)
		if got := maxArea(in); got != c.want {
			t.Errorf("%s: maxArea(%v) = %d, want %d", c.name, c.in, got, c.want)
		}
	}
}

// brute force: every legal pair. O(n^2), obviously correct, used as the oracle.
func maxAreaBrute(heights []int) int {
	best := 0
	for i := 0; i < len(heights); i++ {
		for j := i + 1; j < len(heights); j++ {
			if a := (j - i) * min(heights[i], heights[j]); a > best {
				best = a
			}
		}
	}
	return best
}

// This is the one that actually answers the "what about ties?" question:
// heights capped low so equal values collide constantly.
func TestMaxAreaAgainstBruteForce(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	for _, maxH := range []int{1, 2, 3, 10, 10000} {
		for trial := 0; trial < 500; trial++ {
			n := 2 + rng.Intn(30)
			in := make([]int, n)
			for i := range in {
				in[i] = rng.Intn(maxH + 1)
			}

			want := maxAreaBrute(in)
			if got := maxArea(append([]int(nil), in...)); got != want {
				t.Fatalf("maxArea(%v) = %d, want %d (maxH=%d)", in, got, want, maxH)
			}
		}
	}
}
