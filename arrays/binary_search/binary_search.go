package binarysearch

func getLength32(nums []int32) int32 {
	return int32(len(nums))
}

func binarySearch(nums []int32, target int32) int32 {
	high := getLength32(nums) - 1
	low := int32(0)

	for low <= high {
		middle := low + (high-low)/2 // to avoid overflowing

		if target == nums[middle] {
			return middle
		}

		if target < nums[middle] {
			high = middle - 1
		} else {
			low = middle + 1
		}
	}

	return -1
}
