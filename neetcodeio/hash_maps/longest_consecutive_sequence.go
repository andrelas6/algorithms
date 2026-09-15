package hashmaps

/*
 * Yellow - I found an almost good solution but needed the hint so that I didn't need to iterate on all elements but only the starters
 * Time: O(n)
 * Space: O(2n)
 */
func longestConsecutive(nums []int) int {
	// to make get O(1)
	mapValueToPresence := make(map[int]struct{})
	// find sequence starters - so that iteration only goes on the starters
	sequenceStarters := make([]int, 0)

	//  map value -> occurrence
	for _, v := range nums {
		mapValueToPresence[v] = struct{}{}
	}

	// if there is a value to the left, then it can't start a sequence
	for k, _ := range mapValueToPresence {
		_, ok := mapValueToPresence[k-1]

		if !ok {
			sequenceStarters = append(sequenceStarters, k)
		}
	}

	maxLength := 0

	for _, value := range sequenceStarters {
		temp := value
		lengthCounter := 1
		for {
			temp += 1
			_, ok := mapValueToPresence[temp]

			if ok {
				lengthCounter++
			} else {
				break
			}
		}
		if maxLength < lengthCounter {
			maxLength = lengthCounter
		}
	}

	return maxLength
}
