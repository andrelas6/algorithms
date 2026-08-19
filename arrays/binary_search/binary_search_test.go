package binarysearch

import "testing"

func TestBinarySearch(t *testing.T) {
	cases := []struct {
		name   string
		nums   []int32
		target int32
		want   int32
	}{
		// examples from the statement
		{"example 1", []int32{1, 2, 3, 4, 5}, 3, 2},
		{"example 2", []int32{2, 4, 6, 8, 10, 12, 14, 16}, 16, 7},

		// samples
		{"empty array", []int32{}, 5, -1},
		{"single element, found", []int32{10}, 10, 0},
		{"single element, absent", []int32{10}, 5, -1},

		// every index of the same array must be findable
		{"first index", []int32{2, 4, 6, 8, 10, 12, 14, 16}, 2, 0},
		{"index 1", []int32{2, 4, 6, 8, 10, 12, 14, 16}, 4, 1},
		{"index 3 (the first midpoint)", []int32{2, 4, 6, 8, 10, 12, 14, 16}, 8, 3},
		{"index 6", []int32{2, 4, 6, 8, 10, 12, 14, 16}, 14, 6},

		// absent: below, above, and in the gaps
		{"below the range", []int32{2, 4, 6, 8, 10, 12, 14, 16}, 1, -1},
		{"above the range", []int32{2, 4, 6, 8, 10, 12, 14, 16}, 17, -1},
		{"gap in the middle", []int32{0, 10, 20, 30, 40, 50, 60, 70, 80, 90}, 35, -1},
		{"gap near the start", []int32{0, 10, 20, 30, 40, 50, 60, 70, 80, 90}, 5, -1},

		// two elements, both parities of the range
		{"two elements, first", []int32{1, 2}, 1, 0},
		{"two elements, second", []int32{1, 2}, 2, 1},
		{"two elements, absent", []int32{1, 2}, 3, -1},

		// negative values are in range per the constraints
		{"negatives", []int32{-1000000000, -5, 0, 7, 1000000000}, -5, 1},
		{"negative target absent", []int32{-1000000000, -5, 0, 7, 1000000000}, -6, -1},
	}

	for _, c := range cases {
		if got := binarySearch(c.nums, c.target); got != c.want {
			t.Errorf("%s: binarySearch(%v, %d) = %d, want %d", c.name, c.nums, c.target, got, c.want)
		}
	}
}

// Every index of a large array must be findable, and no lookup may hang.
func TestBinarySearchExhaustive(t *testing.T) {
	const n = 100000
	nums := make([]int32, n)
	for i := range nums {
		nums[i] = int32(i) * 2 // even numbers 0..2n-2
	}

	for i := int32(0); i < n; i++ {
		if got := binarySearch(nums, i*2); got != i {
			t.Fatalf("binarySearch(nums, %d) = %d, want %d", i*2, got, i)
		}
		if got := binarySearch(nums, i*2+1); got != -1 { // odd values are absent
			t.Fatalf("binarySearch(nums, %d) = %d, want -1", i*2+1, got)
		}
	}
}
