package twopointer

import "testing"

func TestTwoSum(t *testing.T) {
	cases := []struct {
		name    string
		numbers []int
		target  int
		want    []int
	}{
		// samples - answers are 1-indexed
		{"sample", []int{1, 2, 3, 4}, 3, []int{1, 2}},
		{"leetcode sample", []int{2, 7, 11, 15}, 9, []int{1, 2}},
		{"middle pair", []int{2, 3, 4}, 6, []int{1, 3}},
		{"negatives", []int{-1, 0}, -1, []int{1, 2}},

		// smallest legal input
		{"two elements", []int{1, 2}, 3, []int{1, 2}},

		// duplicates
		{"identical pair", []int{3, 3}, 6, []int{1, 2}},
		{"duplicates around the answer", []int{1, 1, 2, 3}, 5, []int{3, 4}},

		// answer at the ends
		{"first and last", []int{1, 2, 3, 9}, 10, []int{1, 4}},
		{"last two", []int{1, 2, 8, 9}, 17, []int{3, 4}},

		// zeros and negatives
		{"zero sum", []int{-3, 0, 3}, 0, []int{1, 3}},
		{"both negative", []int{-5, -3, -1}, -8, []int{1, 2}},
	}

	for _, c := range cases {
		numbers := append([]int(nil), c.numbers...)
		got := twoSum(numbers, c.target)

		if len(got) != 2 {
			t.Errorf("%s: twoSum(%v, %d) returned %v, want two indices", c.name, c.numbers, c.target, got)
			continue
		}
		if got[0] != c.want[0] || got[1] != c.want[1] {
			t.Errorf("%s: twoSum(%v, %d) = %v, want %v", c.name, c.numbers, c.target, got, c.want)
		}
	}
}

// the indices must be 1-based, ascending, and actually sum to the target
func TestTwoSumIndicesAreConsistent(t *testing.T) {
	cases := []struct {
		numbers []int
		target  int
	}{
		{[]int{1, 2, 3, 4}, 3},
		{[]int{2, 7, 11, 15}, 9},
		{[]int{-3, 0, 3}, 0},
		{[]int{1, 1, 2, 3}, 5},
		{[]int{3, 3}, 6},
	}

	for _, c := range cases {
		got := twoSum(append([]int(nil), c.numbers...), c.target)
		if got[0] < 1 || got[1] > len(c.numbers) {
			t.Errorf("twoSum(%v, %d) = %v: indices out of 1..n range", c.numbers, c.target, got)
			continue
		}
		if got[0] >= got[1] {
			t.Errorf("twoSum(%v, %d) = %v: indices must be ascending", c.numbers, c.target, got)
			continue
		}
		if sum := c.numbers[got[0]-1] + c.numbers[got[1]-1]; sum != c.target {
			t.Errorf("twoSum(%v, %d) = %v: values sum to %d", c.numbers, c.target, got, sum)
		}
	}
}

func TestTwoSumDoesNotMutate(t *testing.T) {
	in := []int{1, 2, 3, 4}
	before := append([]int(nil), in...)
	twoSum(in, 5)
	for i := range in {
		if in[i] != before[i] {
			t.Fatalf("input mutated: got %v, want %v", in, before)
		}
	}
}
