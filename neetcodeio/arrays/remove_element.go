package arrays

func removeElement(nums []int, val int) int {
	// retunr amount of non val values
	// shift val to the right

	// in-place -> two pointers
	// [2]
	// [1,2]
	// [2,1]
	// [2,1,1] 2
	var j int
	// 2
	// [3,2,2,3] 2
	// i = 0, j = 0, 3 != 2? True -> nums[0] = nums[0] -> [3,2,2,3]
	// i = 1, j = 1, 2 != 2? False -> i = 1, j = 1, [3,2,2,3]
	// i = 2, j = 1, 2 != 2? False -> j = 1, [3,2,2,3]
	// i = 3, j = 1, 3 != 2? True -> nums[1] = nums[3], j = 3 -> [3,3,2,2]
	// [1,2,3,4,5] 3
	// 1 != 3 -> j = 1
	// i = 1 j = 1, 2 != 3
	// i = 2, j = 2, 3 != 3? False -> j = 2
	// i = 3, j = 2, 4 != 3? nums[2] = nums[3] [1,2,4,4,5], j = 3
	// i = 4, j = 3, 4 != 3? nums [3] = nums[4] [1,2,4,5,5], j = 4
	for i := 0; i < len(nums); i++ {
		// value is val, nothing to do
		value := nums[i]
		if value != val {
			nums[j] = nums[i]
			j++
		}
	}

	return j
}
