package findduplicate

import (
	"math/rand"
	"slices"
	"testing"
)

func TestFindDuplicate(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		// samples
		{"sample 1", []int{1, 3, 4, 2, 2}, 2},
		{"sample 2", []int{3, 1, 3, 4, 2}, 3},
		{"sample 3: all the same", []int{3, 3, 3, 3, 3}, 3},

		// smallest input: n = 1
		{"n = 1", []int{1, 1}, 1},

		// duplicate repeated many times, some values missing
		{"repeated 3 times", []int{2, 2, 2, 1}, 2},
		{"repeated, others missing", []int{1, 4, 4, 2, 4}, 4},

		// position of the duplicate
		{"duplicate at the ends", []int{2, 1, 3, 4, 2}, 2},
		{"duplicate next to each other", []int{1, 2, 3, 4, 4}, 4},
		{"duplicate is 1", []int{1, 1, 2, 3, 4}, 1},
		{"duplicate is n", []int{4, 1, 2, 3, 4}, 4},

		// long tail before the cycle
		{"long tail", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 5}, 5},
		// cycle closes right back on the entry
		{"self loop at the entry", []int{2, 2, 1}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nums := slices.Clone(tt.nums)
			got := findDuplicate(nums)
			if got != tt.expected {
				t.Errorf("findDuplicate(%v): expected %d, got %d", tt.nums, tt.expected, got)
			}
			if !slices.Equal(nums, tt.nums) {
				t.Errorf("input was modified: before %v, after %v", tt.nums, nums)
			}
		})
	}
}

// Random inputs where we plant the duplicate ourselves, so we know the answer.
func TestFindDuplicateRandom(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	for range 500 {
		n := 1 + r.Intn(50)
		dup := 1 + r.Intn(n)

		// a permutation of 1..n, plus extra copies of dup overwriting some values
		nums := make([]int, 0, n+1)
		for v := 1; v <= n; v++ {
			nums = append(nums, v)
		}
		nums = append(nums, dup)
		extra := r.Intn(n)
		for range extra {
			j := r.Intn(len(nums))
			nums[j] = dup
		}
		r.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })

		input := slices.Clone(nums)
		got := findDuplicate(nums)
		if got != dup {
			t.Fatalf("findDuplicate(%v): expected %d, got %d", input, dup, got)
		}
	}
}
