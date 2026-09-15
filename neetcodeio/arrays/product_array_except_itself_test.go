package arrays

import (
	"slices"
	"testing"
)

func TestProductExceptSelf(t *testing.T) {
	cases := []struct {
		name string
		in   []int
		want []int
	}{
		// samples from the statement
		{"neetcode sample", []int{1, 2, 4, 6}, []int{48, 24, 12, 8}},
		{"neetcode with zero", []int{-1, 0, 1, 2, 3}, []int{0, -6, 0, 0, 0}},
		{"leetcode sample", []int{1, 2, 3, 4}, []int{24, 12, 8, 6}},
		{"leetcode with zeroes", []int{-1, 1, 0, -3, 3}, []int{0, 0, 9, 0, 0}},

		// degenerate
		{"single element", []int{5}, []int{1}},
		{"two elements", []int{3, 7}, []int{7, 3}},
		{"two elements with zero", []int{0, 7}, []int{7, 0}},

		// zeroes: one zero leaves exactly one non-zero result, two or more
		// make everything zero
		{"one zero at the start", []int{0, 1, 2, 3}, []int{6, 0, 0, 0}},
		{"one zero in the middle", []int{1, 2, 0, 3}, []int{0, 0, 6, 0}},
		{"one zero at the end", []int{1, 2, 3, 0}, []int{0, 0, 0, 6}},
		{"two zeroes", []int{0, 0, 1, 2}, []int{0, 0, 0, 0}},
		{"two zeroes apart", []int{0, 1, 0, 2}, []int{0, 0, 0, 0}},
		{"three zeroes", []int{0, 0, 0}, []int{0, 0, 0}},
		{"all zeroes", []int{0, 0, 0, 0}, []int{0, 0, 0, 0}},
		{"single zero", []int{0}, []int{1}},

		// signs
		{"all negative", []int{-1, -2, -3}, []int{6, 3, 2}},
		{"all negative even count", []int{-1, -2, -3, -4}, []int{-24, -12, -8, -6}},
		{"mixed signs", []int{-2, 3, -4}, []int{-12, 8, -6}},
		{"one negative", []int{2, -3, 4}, []int{-12, 8, -6}},

		// ones and repeats
		{"all ones", []int{1, 1, 1, 1}, []int{1, 1, 1, 1}},
		{"repeated value", []int{2, 2, 2}, []int{4, 4, 4}},
		{"with a one", []int{1, 5}, []int{5, 1}},

		// bigger, still within the stated bounds
		{"eight values", []int{2, 3, 4, 5}, []int{60, 40, 30, 24}},
		{"negatives and zero", []int{-5, 0, -5}, []int{0, 25, 0}},
	}

	for _, c := range cases {
		in := slices.Clone(c.in)

		got := productExceptSelf(in)

		if !slices.Equal(got, c.want) {
			t.Errorf("%s: productExceptSelf(%v) = %v, want %v", c.name, c.in, got, c.want)
		}
	}
}

// obviously-correct oracle: for each index, multiply everything else
func productBrute(nums []int) []int {
	out := make([]int, len(nums))
	for i := range nums {
		p := 1
		for j := range nums {
			if j != i {
				p *= nums[j]
			}
		}
		out[i] = p
	}
	return out
}

func TestProductExceptSelfAgainstBruteForce(t *testing.T) {
	// small deterministic PRNG so failures reproduce
	seed := uint64(20260914)
	next := func(n int) int {
		seed = seed*6364136223846793005 + 1442695040888963407
		return int((seed >> 33) % uint64(n))
	}

	for trial := 0; trial < 5000; trial++ {
		n := 1 + next(9)
		in := make([]int, n)
		for i := range in {
			// -4..4, so zeroes and sign flips are common and products stay small
			in[i] = next(9) - 4
		}

		want := productBrute(slices.Clone(in))
		got := productExceptSelf(slices.Clone(in))

		if !slices.Equal(got, want) {
			t.Fatalf("productExceptSelf(%v) = %v, want %v", in, got, want)
		}
	}
}

// The invariant that defines the problem: for any index, multiplying the
// result by the element gives the product of everything. Checked only where
// no zero is present, since a zero makes the relation vacuous.
func TestProductExceptSelfInvariant(t *testing.T) {
	inputs := [][]int{
		{1, 2, 3, 4},
		{2, 3, 4, 5},
		{-1, 2, -3, 4},
		{5},
		{3, 7},
		{1, 1, 1},
	}

	for _, in := range inputs {
		total := 1
		for _, v := range in {
			total *= v
		}

		got := productExceptSelf(slices.Clone(in))

		for i, v := range in {
			if got[i]*v != total {
				t.Errorf("productExceptSelf(%v)[%d] = %d; %d * %d = %d, want the full product %d",
					in, i, got[i], got[i], v, got[i]*v, total)
			}
		}
	}
}

// The result is a fresh slice: same length as the input, and writing to it
// must not reach back into the caller's data.
func TestProductExceptSelfReturnsFreshSlice(t *testing.T) {
	in := []int{1, 2, 3, 4}
	before := slices.Clone(in)

	got := productExceptSelf(in)

	if len(got) != len(in) {
		t.Fatalf("result has length %d, want %d", len(got), len(in))
	}
	if !slices.Equal(in, before) {
		t.Errorf("productExceptSelf modified its input: %v became %v", before, in)
	}

	got[0] = 999
	if in[0] != before[0] {
		t.Errorf("result aliases the input: writing to the result changed nums[0] to %d", in[0])
	}
}

// Constraints allow 1000 elements. Keep the values at 1 and -1 so the products
// stay in range while every position still has to be computed.
func TestProductExceptSelfLong(t *testing.T) {
	n := 1000

	ones := make([]int, n)
	for i := range ones {
		ones[i] = 1
	}
	got := productExceptSelf(slices.Clone(ones))
	for i, v := range got {
		if v != 1 {
			t.Fatalf("all ones: result[%d] = %d, want 1", i, v)
		}
	}

	// alternating signs: an even count of -1 leaves the full product at 1, so
	// each result is 1 divided by its own element, i.e. the element itself
	alt := make([]int, n)
	for i := range alt {
		if i%2 == 0 {
			alt[i] = 1
		} else {
			alt[i] = -1
		}
	}
	got = productExceptSelf(slices.Clone(alt))
	for i, v := range got {
		want := alt[i]
		if v != want {
			t.Fatalf("alternating: result[%d] = %d, want %d", i, v, want)
		}
	}

	// a single zero somewhere in the middle of a long run
	withZero := make([]int, n)
	for i := range withZero {
		withZero[i] = 1
	}
	withZero[500] = 0
	got = productExceptSelf(slices.Clone(withZero))
	for i, v := range got {
		want := 0
		if i == 500 {
			want = 1
		}
		if v != want {
			t.Fatalf("single zero: result[%d] = %d, want %d", i, v, want)
		}
	}
}
