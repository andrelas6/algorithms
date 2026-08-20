package greedybyend

import (
	"slices"
	"sort"
)

func maximizeNonOverlappingMeetings(meetings [][]int32) int32 {
	// Write your code here
	// overlapping things -> sorting by end time
	copy := slices.Clone(meetings)
	sort.Slice(copy, func(i, j int) bool { return copy[i][1] < copy[j][1] })
	var good int32 = 0

	// 1Start, 2End -> compare
	// overlap  -> bad, need to discard this value somehow
	// no overlap -> awesome, need to add to the count
	var lastEndVal int32 = -1
	for i := 0; i < len(copy); i++ {
		currentStartVal := copy[i][0]
		if lastEndVal <= currentStartVal {
			// no overlap - very good
			// let's say a pervious iteration had overlap and this current one has no overlap
			// lastEndVal MUST NOT BE I, or can it? actually yes, it's i. GOod easy.
			lastEndVal = copy[i][1]
			good++
		} else {
			// we have an overlap
			// currentStart is bad
			// i++ will occurr so I don't need to discard here
			// but it should never be part of lastEndVal, ever
			// how to guarantee that? what makes it part of the last val? how will I set last val and when? now? NO. ONly int he else right?

		}

	}
	return good
}
