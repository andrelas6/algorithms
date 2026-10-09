package carfleet

import "testing"

func TestCarFleet(t *testing.T) {
	cases := []struct {
		name     string
		target   int
		position []int
		speed    []int
		want     int
	}{
		// samples from the statement
		{"sample 1", 12, []int{10, 8, 0, 5, 3}, []int{2, 4, 1, 1, 3}, 3},
		{"sample 2", 10, []int{3}, []int{3}, 1},
		{"sample 3", 100, []int{0, 2, 4}, []int{4, 2, 1}, 1},

		// degenerate
		{"single car at start", 10, []int{0}, []int{1}, 1},
		{"single car one step away", 1, []int{0}, []int{1}, 1},

		// nobody catches anybody
		{"same speed, all separate", 10, []int{0, 1, 2}, []int{1, 1, 1}, 3},
		{"front car is faster", 10, []int{6, 8}, []int{3, 2}, 2},
		{"back car just too slow (catches after target)", 20, []int{0, 10}, []int{19, 10}, 2},

		// ties: arriving at the target together counts as one fleet
		{"tie at the target, whole times", 10, []int{0, 5}, []int{2, 1}, 1},
		{"tie at the target, fractional times", 10, []int{0, 5}, []int{6, 3}, 1},

		// one slow car blocks everyone behind it
		{"slow leader absorbs a chain", 10, []int{9, 2, 1, 0}, []int{1, 10, 10, 10}, 1},
		// a new fleet becomes the thing later cars merge into
		{"merge into the second fleet, not the first", 10, []int{7, 4, 1}, []int{1, 1, 10}, 2},
		// back car would catch the front car but the middle one is slower still
		{"middle car is the bottleneck", 10, []int{8, 5, 0}, []int{1, 1, 100}, 2},

		// input not sorted by position
		{"unsorted input", 12, []int{3, 5, 0, 8, 10}, []int{3, 1, 1, 4, 2}, 3},
		{"ascending positions", 10, []int{0, 1, 2, 3}, []int{1, 2, 3, 4}, 4},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := carFleet(c.target, c.position, c.speed); got != c.want {
				t.Errorf("carFleet(%d, %v, %v) = %d, want %d", c.target, c.position, c.speed, got, c.want)
			}
		})
	}
}
