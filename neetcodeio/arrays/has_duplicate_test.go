package arrays

import (
	"math/rand"
	"testing"
)

func TestHasDuplicate(t *testing.T) {
	cases := []struct {
		name string
		in   []int
		want bool
	}{
		// samples
		{"has a repeat", []int{1, 2, 3, 3}, true},
		{"all distinct", []int{1, 2, 3, 4}, false},

		// degenerate
		{"empty", []int{}, false},
		{"nil", nil, false},
		{"single", []int{7}, false},

		// duplicate at the very start / very end
		{"adjacent at start", []int{5, 5, 1, 2}, true},
		{"far apart", []int{5, 1, 2, 3, 5}, true},

		// negatives and zero
		{"negatives distinct", []int{-1, -2, -3}, false},
		{"negatives repeat", []int{-1, 2, -1}, true},
		{"zeroes", []int{0, 0}, true},
		{"zero and negative zero are the same int", []int{0, 1, 0}, true},

		// everything the same
		{"all identical", []int{4, 4, 4, 4}, true},
	}

	for _, c := range cases {
		in := append([]int(nil), c.in...)
		if got := hasDuplicate(in); got != c.want {
			t.Errorf("%s: hasDuplicate(%v) = %v, want %v", c.name, c.in, got, c.want)
		}
	}
}

func TestHasDuplicateDoesNotMutate(t *testing.T) {
	in := []int{3, 1, 2, 1}
	before := append([]int(nil), in...)
	hasDuplicate(in)
	for i := range in {
		if in[i] != before[i] {
			t.Fatalf("input mutated: got %v, want %v", in, before)
		}
	}
}

func TestHasDuplicateAgainstBruteForce(t *testing.T) {
	rng := rand.New(rand.NewSource(11))
	for trial := 0; trial < 2000; trial++ {
		n := rng.Intn(15)
		in := make([]int, n)
		for i := range in {
			in[i] = rng.Intn(8) - 4
		}

		want := false
		for i := 0; i < len(in) && !want; i++ {
			for j := i + 1; j < len(in); j++ {
				if in[i] == in[j] {
					want = true
					break
				}
			}
		}

		if got := hasDuplicate(append([]int(nil), in...)); got != want {
			t.Fatalf("hasDuplicate(%v) = %v, want %v", in, got, want)
		}
	}
}
