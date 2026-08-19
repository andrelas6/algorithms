package lowerbound

import (
	"math/rand"
	"testing"
)

func TestFindFirstOccurrence(t *testing.T) {
	cases := []struct {
		name   string
		nums   []int32
		target int32
		want   int32
	}{
		// example + samples from the statement
		{"example", []int32{1, 2, 3, 4, 5}, 3, 2},
		{"empty array", []int32{}, 5, -1},
		{"single element, found", []int32{3}, 3, 0},
		{"single element, absent", []int32{3}, 7, -1},

		// duplicates: the whole point
		{"all duplicates", []int32{3, 3, 3, 3, 3}, 3, 0},
		{"run at the start", []int32{2, 2, 2, 5, 9}, 2, 0},
		{"run at the end", []int32{1, 4, 7, 7, 7}, 7, 2},
		{"run in the middle", []int32{1, 3, 5, 5, 5, 8, 9}, 5, 2},
		{"run of two", []int32{1, 2, 2, 3}, 2, 1},
		{"long run, even length", []int32{0, 4, 4, 4, 4, 4, 4, 9}, 4, 1},

		// absent
		{"all elements less", []int32{1, 2, 3}, 10, -1},
		{"all elements greater", []int32{5, 6, 7}, 1, -1},
		{"gap in the middle", []int32{1, 2, 4, 5}, 3, -1},

		// every index of one array must be findable
		{"first index", []int32{1, 2, 3, 4, 5}, 1, 0},
		{"last index", []int32{1, 2, 3, 4, 5}, 5, 4},

		// negatives, per the constraints
		{"negatives with duplicates", []int32{-1000000000, -5, -5, 0, 1000000000}, -5, 1},
		{"min value", []int32{-1000000000, -5, 0}, -1000000000, 0},
		{"max value", []int32{-5, 0, 1000000000}, 1000000000, 2},
	}

	for _, c := range cases {
		if got := findFirstOccurrence(c.nums, c.target); got != c.want {
			t.Errorf("%s: findFirstOccurrence(%v, %d) = %d, want %d", c.name, c.nums, c.target, got, c.want)
		}
	}
}

// linear scan, obviously correct, used as the oracle
func firstOccurrenceBrute(nums []int32, target int32) int32 {
	for i, v := range nums {
		if v == target {
			return int32(i)
		}
	}
	return -1
}

// Random sorted arrays with heavy duplication, checked against the brute force.
func TestFindFirstOccurrenceAgainstBruteForce(t *testing.T) {
	r := rand.New(rand.NewSource(1))

	for trial := 0; trial < 2000; trial++ {
		n := r.Intn(40)
		nums := make([]int32, n)
		v := int32(-5)
		for i := range nums {
			v += int32(r.Intn(3)) // 0 keeps duplicates coming
			nums[i] = v
		}

		for target := int32(-7); target <= v+2; target++ {
			want := firstOccurrenceBrute(nums, target)
			if got := findFirstOccurrence(nums, target); got != want {
				t.Fatalf("findFirstOccurrence(%v, %d) = %d, want %d", nums, target, got, want)
			}
		}
	}
}
