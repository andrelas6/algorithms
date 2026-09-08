package twopointer

/*
 * YELLOW
 * I couldn't come up with the logic on which pointer to move. I got a hint that it's the smaller and then it made sense. I thought
 * that if both bars have equal height, I needed to pick the one furthest from the middle to move. But apparently that's not needed (submission passed). Dunno why.
 *
 * Solution in 25 min
 * Time: O(n)
 * Space: O(1)
 */

func maxArea(heights []int) int {
	var left, right int
	right = len(heights) - 1

	var biggestWaterAmount int

	// calculate, then move smaller pointer only
	for left < right {
		leftBar := heights[left]
		rightBar := heights[right]

		// width * height
		currentWaterAmount := (right - left) * min(leftBar, rightBar)
		if biggestWaterAmount < currentWaterAmount {
			biggestWaterAmount = currentWaterAmount
		}

		// the moving
		if moveLeft(left, right, heights) {
			left++
		} else {
			right--
		}
	}

	return biggestWaterAmount
}

func moveLeft(left, right int, nums []int) bool {
	if nums[left] < nums[right] {
		return true
	} else if nums[right] < nums[left] {
		return false
	}

	return true
}
