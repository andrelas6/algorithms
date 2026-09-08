package twopointer

// GREEN: got correct solution, YELLOW

/*
 * 1st cycle: arrays
 * MEDIUM
 * GREEN but maybe yellow due to overcomplicating my map implementation. But I didn't need a hint to fix it so possibly it's green.
 * mistake: overcomplicated data structure. I made my map more complicated than it shoudl be because I found the solution to be very similar to
 * a previous one I had done. So I kinda "copied" without thinking.
 *
 * time complexity: O(n)
 * space complexity: O(n)
 */

func twoSum(numbers []int, target int) []int {
	// loop i
	// inner loop j (i + 1)
	// with two pointers
	// [1,2,3,4] target 3
	//  ^     ^
	//  |
	valAndIndex := make(map[int]int)
	result := make([]int, 2)

	// convert array to map
	// value -> index
	for i, n := range numbers {
		valAndIndex[n] = i
	}

	for i, n := range numbers {
		diff := target - n
		// 3 - 1 = 2
		index, ok := valAndIndex[diff]
		// map[2] -> ql teu index? 1
		// [0 + 1, 1 + 2]

		if index != i && ok {
			result[0] = i + 1
			result[1] = index + 1
			return result
		}
	}

	return result
}
