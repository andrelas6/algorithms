package twopointer

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
	"time"
)

// normalize sorts each triplet and then the list of triplets, so answers can be
// compared whatever order they came out in.
func normalize(triplets [][]int) [][]int {
	out := make([][]int, len(triplets))
	for i, t := range triplets {
		c := slices.Clone(t)
		slices.Sort(c)
		out[i] = c
	}
	slices.SortFunc(out, func(a, b []int) int { return slices.Compare(a, b) })
	return out
}

// threeSumBrute is an independent oracle: try every combination of three
// indices and keep the sorted triplets that sum to zero, deduplicated with a
// map. No sorting of the input, no pointers. O(n^3), fine for small inputs.
func threeSumBrute(nums []int) [][]int {
	seen := map[string]bool{}
	var out [][]int

	for i := range nums {
		for j := i + 1; j < len(nums); j++ {
			for k := j + 1; k < len(nums); k++ {
				if nums[i]+nums[j]+nums[k] != 0 {
					continue
				}
				triplet := []int{nums[i], nums[j], nums[k]}
				slices.Sort(triplet)

				key := fmt.Sprint(triplet)
				if !seen[key] {
					seen[key] = true
					out = append(out, triplet)
				}
			}
		}
	}

	return normalize(out)
}

func TestThreeSumExamples(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want [][]int
	}{
		// samples from the statement
		{"sample 1", []int{-1, 0, 1, 2, -1, -4}, [][]int{{-1, -1, 2}, {-1, 0, 1}}},
		{"sample 2", []int{0, 1, 1}, nil},
		{"sample 3", []int{0, 0, 0}, [][]int{{0, 0, 0}}},

		// degenerate
		{"empty", []int{}, nil},
		{"single value", []int{0}, nil},
		{"two values", []int{0, 0}, nil},
		{"three values that work", []int{-1, 0, 1}, [][]int{{-1, 0, 1}}},
		{"three values that do not", []int{1, 2, 3}, nil},

		// nothing can sum to zero
		{"all positive", []int{1, 2, 3, 4, 5}, nil},
		{"all negative", []int{-1, -2, -3, -4}, nil},
	}

	for _, c := range cases {
		got := normalize(threeSum(slices.Clone(c.nums)))
		want := normalize(c.want)

		if !slices.EqualFunc(got, want, func(a, b []int) bool { return slices.Equal(a, b) }) {
			t.Errorf("%s: threeSum(%v) = %v, want %v", c.name, c.nums, got, want)
		}
	}
}

// Duplicates are where this problem bites: the same triplet must not appear
// twice, even when the input holds several copies of its values.
func TestThreeSumDuplicates(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want [][]int
	}{
		{"four zeroes", []int{0, 0, 0, 0}, [][]int{{0, 0, 0}}},
		{"repeated pair around zero", []int{-2, 0, 0, 2, 2}, [][]int{{-2, 0, 2}}},
		{"same triplet reachable many ways", []int{-1, -1, -1, 2, 2, 2}, [][]int{{-1, -1, 2}}},
		{"two distinct triplets with repeats", []int{-2, -2, 0, 2, 2, 4}, [][]int{{-2, -2, 4}, {-2, 0, 2}}},
		{"duplicates that do not help", []int{1, 1, 1, 1}, nil},
		{"zero needs three copies", []int{0, 0, 1, -1}, [][]int{{-1, 0, 1}}},
		{"triple zero alongside a pair", []int{-3, 0, 0, 0, 3}, [][]int{{-3, 0, 3}, {0, 0, 0}}},
	}

	for _, c := range cases {
		got := normalize(threeSum(slices.Clone(c.nums)))
		want := normalize(c.want)

		if !slices.EqualFunc(got, want, func(a, b []int) bool { return slices.Equal(a, b) }) {
			t.Errorf("%s: threeSum(%v) = %v, want %v", c.name, c.nums, got, want)
		}
	}
}

// Random inputs against the brute force. Small value ranges force plenty of
// duplicates and plenty of triplets.
func TestThreeSumAgainstBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(31, 32))

	for range 2000 {
		n := rng.IntN(12)
		spread := []int{3, 5, 9}[rng.IntN(3)]

		nums := make([]int, n)
		for i := range nums {
			nums[i] = rng.IntN(spread) - spread/2
		}

		got := normalize(threeSum(slices.Clone(nums)))
		want := threeSumBrute(nums)

		if !slices.EqualFunc(got, want, func(a, b []int) bool { return slices.Equal(a, b) }) {
			t.Fatalf("threeSum(%v) = %v, brute force = %v", nums, got, want)
		}
	}
}

// Every triplet must sum to zero, must be buildable from the input values, and
// must appear only once in the answer.
func TestThreeSumOutputIsWellFormed(t *testing.T) {
	inputs := [][]int{
		{-1, 0, 1, 2, -1, -4},
		{-4, -2, -2, -2, 0, 1, 2, 2, 2, 3, 3, 4, 4, 6, 6},
		{0, 0, 0, 0, 0},
		{-5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5},
	}

	for _, nums := range inputs {
		got := threeSum(slices.Clone(nums))

		available := map[int]int{}
		for _, v := range nums {
			available[v]++
		}

		seen := map[string]bool{}
		for _, triplet := range got {
			if len(triplet) != 3 {
				t.Errorf("%v: got a group of %d values, want 3", nums, len(triplet))
				continue
			}
			if sum := triplet[0] + triplet[1] + triplet[2]; sum != 0 {
				t.Errorf("%v: triplet %v sums to %d, want 0", nums, triplet, sum)
			}

			need := map[int]int{}
			for _, v := range triplet {
				need[v]++
			}
			for v, count := range need {
				if available[v] < count {
					t.Errorf("%v: triplet %v needs %d copies of %d, input has %d", nums, triplet, count, v, available[v])
				}
			}

			sorted := slices.Clone(triplet)
			slices.Sort(sorted)
			key := fmt.Sprint(sorted)
			if seen[key] {
				t.Errorf("%v: triplet %v appears more than once", nums, sorted)
			}
			seen[key] = true
		}
	}
}

// Sorting the caller's slice is allowed, but no value may be lost or invented.
func TestThreeSumKeepsTheSameValues(t *testing.T) {
	nums := []int{-1, 0, 1, 2, -1, -4, 3, -3, 0}
	before := slices.Clone(nums)
	slices.Sort(before)

	threeSum(nums)

	after := slices.Clone(nums)
	slices.Sort(after)

	if !slices.Equal(after, before) {
		t.Errorf("values after the call %v are not the values before it %v", after, before)
	}
}

// Constraints allow 3000 values. The answer can hold a lot of triplets, so this
// checks the shape of the output rather than listing it.
func TestThreeSumLarge(t *testing.T) {
	if testing.Short() {
		t.Skip("large inputs skipped in -short mode")
	}

	const n = 3000
	const budget = 2 * time.Second

	rng := rand.New(rand.NewPCG(33, 34))
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.IntN(2001) - 1000
	}

	start := time.Now()
	got := threeSum(slices.Clone(nums))
	elapsed := time.Since(start)

	if elapsed > budget {
		t.Errorf("n=%d: took %v, over the %v budget", n, elapsed.Round(time.Millisecond), budget)
	}

	seen := map[string]bool{}
	for _, triplet := range got {
		if sum := triplet[0] + triplet[1] + triplet[2]; sum != 0 {
			t.Fatalf("triplet %v sums to %d, want 0", triplet, sum)
		}
		sorted := slices.Clone(triplet)
		slices.Sort(sorted)
		key := fmt.Sprint(sorted)
		if seen[key] {
			t.Fatalf("triplet %v appears more than once", sorted)
		}
		seen[key] = true
	}
	t.Logf("n=%d produced %d triplets in %v", n, len(got), elapsed.Round(time.Millisecond))
}
