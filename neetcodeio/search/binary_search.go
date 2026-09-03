package search

func search(nums []int, target int) int {
	return BinarySearch(nums, 0, len(nums)-1, target)
}

func BinarySearch(nums []int, left, right, target int) int {
	// 1,2,3,4,5 4
	//
	for left <= right {
		middle := (left + right) / 2

		if nums[middle] == target {
			return middle
		}

		if target > nums[middle] { // to the right
			return BinarySearch(nums, middle+1, right, target)
		} else {
			return BinarySearch(nums, left, middle-1, target)
		}
	}

	return -1
}
