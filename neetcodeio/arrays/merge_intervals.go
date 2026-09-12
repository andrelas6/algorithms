package arrays

import "slices"

/*
 * RED: I needed to know that sorting was the trick and I needed to know that I could use a sort function from the std library.
 * Time: O(n * log n)
 * Space: O(n)
 */
func merge(intervals [][]int) [][]int {
	slices.SortFunc(intervals, func(a, b []int) int {
		if a[0] < b[0] {
			return -1
		} else if a[0] > b[0] {
			return 1
		} else {
			return 0
		}
	})

	var merged [][]int
	merged = append(merged, intervals[0])

	// overlapping checks
	for i := 1; i < len(intervals); i++ {
		start := intervals[i][0]
		lastEnd := merged[len(merged)-1][1]

		if start <= lastEnd {
			end := intervals[i][1]
			merged[len(merged)-1][1] = max(lastEnd, end)
		} else {
			// add item
			merged = append(merged, intervals[i])
		}
	}

	return merged
}
