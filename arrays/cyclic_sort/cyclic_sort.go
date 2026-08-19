package cyclicsort

func findSmallestMissingPositive(orderNumbers []int32) int32 {
	// Write your code here
	var result int32
	if len(orderNumbers) == 0 {
		return 1
	}
	// sort array
	// walk sorted array until an element is out of place or all elements are walked

	// put every number in their own place
	for i := 0; i < len(orderNumbers); {
		n_index := int32(i)
		// is this invalid? go to next index
		if int32(len(orderNumbers)) < orderNumbers[i] || orderNumbers[i] < 1 || isInCorrectPlace(orderNumbers[i], n_index) {
			i++
			continue
		}

		// weird case -> duplicate value: I try to put it in correct place, and it succeeds
		// but then the swapped value is in wrong position, I swap again,
		// infinite loop
		// thanks for nothing hackerrank

		wannabe_index := orderNumbers[i] - 1 // where 3 should be
		// can I swap? if I can't I should skip
		// how? a index 2 means the array must have length 3 at least

		save_this := orderNumbers[wannabe_index] // will be moved elsewhere I don't care
		save_this_too := orderNumbers[i]         // will be put in correct place
		if save_this == save_this_too {
			i++
			continue
		}
		orderNumbers[i] = save_this
		orderNumbers[wannabe_index] = save_this_too
	}

	for i, n := range orderNumbers {
		if n-1 != int32(i) {
			result = int32(i + 1)
			break
		}
	}

	if result == 0 {
		return int32(len(orderNumbers) + 1)
	}

	return result
}

func isInCorrectPlace(value int32, position int32) bool {
	// all good
	if value-1 == position {
		return true
	}

	// out of place
	return false
}
