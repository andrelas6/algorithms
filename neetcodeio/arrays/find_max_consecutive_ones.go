package arrays

func findMaxConsecutiveOnes(nums []int) int {
	// 1,1,0,1,1
	// max number of consecutive 1s
	// only need the count

	// 1st solution
	/*
		init finalCount, currentCount
		read every value, if 1, add to currentCount, if 0, don't count
		when value is 0, if currentCount > finalCount, replace.

		return finalCount
	*/

	var highestCountYet int
	var currentCount int
	for i := 0; i < len(nums); i++ {
		// 0
		value := nums[i]

		// last one -> needs to make the count
		if value == 1 {
			// true -> cc 1
			currentCount++
		}

		if currentCount > highestCountYet {
			highestCountYet = currentCount
		}

		if value == 0 {
			currentCount = 0
		}
	}

	return highestCountYet
}
