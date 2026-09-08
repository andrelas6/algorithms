package slidingwindow

import (
	"math/rand"
	"testing"
)

func TestMaxProfit(t *testing.T) {
	cases := []struct {
		name string
		in   []int
		want int
	}{
		// samples
		{"buy low sell high", []int{10, 1, 5, 6, 7, 1}, 6},
		{"only falls", []int{10, 8, 7, 5, 2}, 0},
		{"leetcode sample", []int{7, 1, 5, 3, 6, 4}, 5},
		{"strictly decreasing", []int{7, 6, 4, 3, 1}, 0},

		// degenerate
		{"empty", []int{}, 0},
		{"nil", nil, 0},
		{"single day", []int{5}, 0},
		{"two days up", []int{1, 5}, 4},
		{"two days down", []int{5, 1}, 0},

		// flat
		{"all equal", []int{3, 3, 3, 3}, 0},

		// best buy is first, best sell is last
		{"monotonic up", []int{1, 2, 3, 4, 5}, 4},

		// the min comes after the max - must not sell before buying
		{"max then min", []int{9, 1}, 0},
		{"tempting early peak", []int{5, 9, 1, 3}, 4},

		// min and max adjacent at the end
		{"late spike", []int{5, 4, 3, 2, 1, 100}, 99},

		// zeros and equal extremes
		{"zeros", []int{0, 0, 0}, 0},
		{"dip then recover", []int{3, 1, 3}, 2},
	}

	for _, c := range cases {
		in := append([]int(nil), c.in...)
		if got := maxProfit(in); got != c.want {
			t.Errorf("%s: maxProfit(%v) = %d, want %d", c.name, c.in, got, c.want)
		}
	}
}

// every buy/sell pair, obviously correct
func maxProfitPairs(prices []int) int {
	best := 0
	for buy := 0; buy < len(prices); buy++ {
		for sell := buy + 1; sell < len(prices); sell++ {
			if p := prices[sell] - prices[buy]; p > best {
				best = p
			}
		}
	}
	return best
}

func TestMaxProfitAgainstAllPairs(t *testing.T) {
	rng := rand.New(rand.NewSource(21))
	for trial := 0; trial < 3000; trial++ {
		n := rng.Intn(25)
		in := make([]int, n)
		for i := range in {
			in[i] = rng.Intn(20)
		}
		want := maxProfitPairs(in)
		if got := maxProfit(append([]int(nil), in...)); got != want {
			t.Fatalf("maxProfit(%v) = %d, want %d", in, got, want)
		}
	}
}

// LeetCode constraint: up to 10^5 prices.
func maxProfitBenchInput(n int) []int {
	rng := rand.New(rand.NewSource(22))
	out := make([]int, n)
	for i := range out {
		out[i] = rng.Intn(10000)
	}
	return out
}

func BenchmarkMaxProfit(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		in := maxProfitBenchInput(n)
		b.Run(map[int]string{1000: "n=1000", 10000: "n=10000", 100000: "n=100000"}[n], func(b *testing.B) {
			for b.Loop() {
				maxProfit(in)
			}
		})
	}
}
