package arrays

func hasDuplicate(nums []int) bool {
	numAppearanceCount := make(map[int]int)

	for _, n := range nums {
		_, ok := numAppearanceCount[n]

		if !ok {
			numAppearanceCount[n] = 1
		} else {
			return true
		}
	}

	return false
}
