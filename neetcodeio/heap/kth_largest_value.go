package heap

import (
	"cmp"
	"slices"
)

/*
 * GREEN - pretty easy to implement
 * Time: O(n * log(n))
 * Space: O(1) in-place sorting
 */
func findKthLargest(nums []int, k int) int {
	// sort, then use index to find the value
	// n * log(n)

	slices.SortFunc(nums, func(a, b int) int {
		return cmp.Compare(a, b)
	})

	length := len(nums)
	return nums[length-k]
}

/*
 * RED - I had no idea how to implement this improved solution
 * Time: O(n) average case, O(n^2) worst case
 * Space: O(1) in-place sorting
 */
func findKthLargestImproved(nums []int, k int) int {
	// index to be found from backwards
	k = len(nums) - k

	return quickSelect(nums, 0, len(nums)-1, k)
}

func quickSelect(nums []int, left, right, k int) int {
	pivot := right
	pointer := left
	for i := left; i < right; i++ {
		// a bit not intuitive, but it is quicksort and it helps when p != i
		if nums[i] <= nums[pivot] {
			nums[i], nums[pointer] = nums[pointer], nums[i]
			pointer++
		}
	}
	nums[pointer], nums[pivot] = nums[pivot], nums[pointer]
	// pivot goes to pointer

	if pointer > k {
		// k is on the left
		return quickSelect(nums, left, pointer-1, k)
	} else if pointer < k {
		// k is on the right
		return quickSelect(nums, pointer+1, right, k)
	} else {
		// value found
		return nums[pointer]
	}
}
