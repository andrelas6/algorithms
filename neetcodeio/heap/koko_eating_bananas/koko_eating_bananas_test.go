package kokoeatingbananas

import (
	"math/rand/v2"
	"testing"
)

// hoursAt counts hours at speed k with integer ceiling, so no floats are involved.
func hoursAt(piles []int, k int) int {
	total := 0
	for _, p := range piles {
		total += (p + k - 1) / k
	}
	return total
}

// slowestSpeedByScan is an independent oracle: try every speed from 1 upward.
// O(max * n), so only for small piles.
func slowestSpeedByScan(piles []int, h int) int {
	for k := 1; ; k++ {
		if hoursAt(piles, k) <= h {
			return k
		}
	}
}

func TestMinEatingSpeed(t *testing.T) {
	tests := []struct {
		name  string
		piles []int
		h     int
		want  int
	}{
		{"sample 1", []int{1, 4, 3, 2}, 9, 2},
		{"sample 2", []int{25, 10, 23, 4}, 4, 25},
		{"leetcode sample 1", []int{3, 6, 7, 11}, 8, 4},
		{"leetcode sample 2", []int{30, 11, 23, 4, 20}, 5, 30},
		{"leetcode sample 3", []int{30, 11, 23, 4, 20}, 6, 23},
		{"single pile of 1", []int{1}, 1, 1},
		{"single pile, one hour: must eat it all", []int{1_000_000_000}, 1, 1_000_000_000},
		{"single pile, plenty of time: speed 1", []int{7}, 100, 1},
		{"h equals number of piles: answer is max", []int{5, 9, 2}, 3, 9},
		{"all piles equal", []int{4, 4, 4, 4}, 8, 2},
		{"near miss: one hour short of speed 1", []int{3, 3}, 5, 2},
		{"exact fit at speed 1", []int{3, 3}, 6, 1},
		{"answer sits exactly on a divisor", []int{6, 6}, 4, 3},
		{"answer is one above a divisor", []int{7, 6}, 4, 4},
		{
			name:  "huge piles, huge h: speed 1 sums past 32-bit",
			piles: []int{1_000_000_000, 1_000_000_000, 1_000_000_000},
			h:     3_000_000_000,
			want:  1,
		},
		{
			name:  "huge piles where float ceiling must be exact",
			piles: []int{999_999_999, 1_000_000_000},
			h:     3,
			want:  999_999_999, // one pile in 1 hour, the other in 2
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minEatingSpeed(tt.piles, tt.h); got != tt.want {
				t.Errorf("minEatingSpeed(%v, %d) = %d, want %d", tt.piles, tt.h, got, tt.want)
			}
		})
	}
}

func TestMinEatingSpeedRandomAgainstOracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(7, 11))
	for i := range 2000 {
		n := 1 + rng.IntN(8)
		piles := make([]int, n)
		total := 0
		for j := range piles {
			piles[j] = 1 + rng.IntN(50)
			total += piles[j]
		}
		h := n + rng.IntN(total+5) // h >= len(piles), sometimes more than enough
		got := minEatingSpeed(piles, h)
		want := slowestSpeedByScan(piles, h)
		if got != want {
			t.Fatalf("case %d: minEatingSpeed(%v, %d) = %d, want %d", i, piles, h, got, want)
		}
	}
}

// Max size: 10^4 piles of 10^9. Checks the answer is feasible and minimal.
func TestMinEatingSpeedMaxSize(t *testing.T) {
	rng := rand.New(rand.NewPCG(3, 5))
	piles := make([]int, 10_000)
	for i := range piles {
		piles[i] = 1 + rng.IntN(1_000_000_000)
	}
	for _, h := range []int{10_000, 50_000, 1_000_000, 1_000_000_000} {
		k := minEatingSpeed(piles, h)
		if hoursAt(piles, k) > h {
			t.Errorf("h=%d: speed %d is too slow (%d hours)", h, k, hoursAt(piles, k))
		}
		if k > 1 && hoursAt(piles, k-1) <= h {
			t.Errorf("h=%d: speed %d is not minimal, %d also works", h, k, k-1)
		}
	}
}
