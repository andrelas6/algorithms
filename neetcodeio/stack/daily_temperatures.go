package stack

/*
 * RED: I could the brute force but I had absolutely no idea how to use the monotonic stack.
 * What is this term "monotonic" btw?
 * Time: O(n) because every element is visited once
 * Space: O(n) due to the stack
 */
func dailyTemperatures(temperatures []int) []int {
	// result is all zeroes - fine when temp has no warmer day in the future

	// use a stack to remember the previous day
	// popping an element means checking with the previous day, then updating the count
	// so this is like backwards traversal but with O(n) time

	type TempAndIndex struct {
		Temp  int
		Index int
	}
	stack := make([]TempAndIndex, 0)
	result := make([]int, len(temperatures))

	// 30, 39, 30, 36, 40

	for i := 0; i < len(temperatures); i++ {
		currentTemp := temperatures[i]
		var previous TempAndIndex
		if len(stack) != 0 {
			previous = stack[len(stack)-1]
		}

		// [39 30, 36]
		// 40 > 36
		// 36 -> 1
		// [39 30]
		// 40 > 30,  30 -> 2, [39]
		// 40 < 39, 39 -> 3, []
		for len(stack) != 0 && currentTemp > previous.Temp {
			previous = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			result[previous.Index] = i - previous.Index

			if len(stack) != 0 {
				previous = stack[len(stack)-1]
			}
		}

		stack = append(stack, TempAndIndex{Temp: temperatures[i], Index: i})

		// if currentTemp > previous.Temp {
		// 	stack = stack[:len(stack) - 1]
		// 	result[previousTemp.Index] = i - previousTemp.Index
		// }
	}

	return result
}
