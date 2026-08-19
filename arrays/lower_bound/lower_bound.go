package lowerbound

func findFirstOccurrence(nums []int32, target int32) int32 {
	// Write your code here
	// ordered and find -> binary search
	// duplicates ???
	// state: low, high and best found value
	// invariant: if target is in array, it is within low-high bounds
	low := 0
	high := len(nums) - 1
	var found int32 = -1
	// start 18:48
	// end 18:54

	for low <= high {
		middle := low + (high-low)/2

		if nums[middle] == target {
			found = int32(middle)
		}

		if target > nums[middle] {
			low = middle + 1
		} else {
			high = middle - 1
		}
	}

	return found
}
