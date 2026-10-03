package kokoeatingbananas

import "math"

func minEatingSpeed(piles []int, h int) int {
	// k - minimum k
	// each hour is an iteration
	// iterate as less as possible to eat all bannas
	// piles 1,4,3,2
	// find max value in piles -> that is the max eating per hour rate
	// binary search on 1..max eating per hour rate (1 is the minimum)
	// for each value in bs, find the amount of hours, store
	// repeat until L > R

	// find the max
	var max int = math.MinInt
	for _, p := range piles {
		if p > max {
			max = p
		}
	}

	// binary search
	var left, right, mid int

	right = max - 1
	var result int = math.MaxInt
	for left <= right {
		mid = (left + right) / 2
		var hoursToFinishEating int

		// check eating hour rate
		for _, p := range piles {
			hours := int(math.Ceil(float64(p) / float64(mid+1)))
			hoursToFinishEating += hours
		}
		if hoursToFinishEating <= h {
			result = min(result, mid+1)
			right = mid - 1
		} else {
			left = mid + 1
		}
	}

	return result
}
