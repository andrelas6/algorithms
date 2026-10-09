package carfleet

import "slices"

func carFleet(target int, position []int, speed []int) int {
	// O(n * m) Time
	// (O m) Space
	// iterate on 0 to target
	// iterate on positions
	// SLOW: ever iteration, sum speed[i] to position[i]
	// at the end, iterate on position, add to a set
	// find size of set

	// optimize: remove the M of the time complexity
	// sort desc
	// for each car, compare with the next one and check if they meet before or at target
	// use a stack to get always the last car checked with the actual
	// Time O(n log n)
	// Space O(n) stack & fleets slice

	fleets := make([][]int, 0)

	for i := 0; i < len(position); i++ {
		fleets = append(fleets, []int{position[i], speed[i]})
	}

	slices.SortFunc(fleets, func(a, b []int) int {
		if a[0] < b[0] {
			return 1
		}

		if a[0] > b[0] {
			return -1
		}

		return 0
	})

	stack := make([][]int, 0)
	stack = append(stack, fleets[0])
	for _, fleet := range fleets {
		car := stack[len(stack)-1]
		iters := float64((target - car[0])) / float64(car[1])

		stack = append(stack, fleet)
		fleetIters := float64((target - fleet[0])) / float64(fleet[1])

		if fleetIters <= iters {
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack)
}
