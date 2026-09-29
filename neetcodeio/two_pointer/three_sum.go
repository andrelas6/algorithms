package twopointer

import "slices"

/*
 * RED
 * Time: O(n^2)
 * Space: O(1) but depends on how slices.Sort works.
 */
func threeSum(nums []int) [][]int {
	// sort array first, then usen two pointers
	// but edge case with repeating numbers is tough

	// sort ascending
	slices.Sort(nums)
	result := make([][]int, 0)

	for i := 0; i < len(nums); i++ {
		if i != 0 && nums[i] == nums[i-1] {
			continue
		}

		var left, right int
		left = i + 1
		right = len(nums) - 1
		// pointers loop
		for left < right {
			sum := nums[i] + nums[left] + nums[right]
			if sum < 0 {
				left++
			}
			if sum > 0 {
				right--
			}

			if sum == 0 {
				result = append(result, []int{nums[i], nums[left], nums[right]})

				currentLeft := nums[left]
				for left < right && nums[left] == currentLeft {
					left++
				}
			}
		}
	}

	return result

}
