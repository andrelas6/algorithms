package slidingwindow

// BROKEN AS IT STANDS - kept only to show the shape of the slow version.
// `j = i` was removed but `charToIndex = make(...)` stayed, and those two were a pair:
// clearing the map is only safe if j rewinds to rebuild it. Without the rewind, chars
// still inside [i..j] get wiped and stop registering as duplicates.
// Fails on "bbabaaa" (returns 3, want 2). Restore `j = i` to get the original back.
func lengthOfLongestSubstring(s string) int {
	// hashmap char -> index
	// sliding window
	// result state
	charToIndex := make(map[byte]int)
	var i, j int
	// z x y z x y z
	var result int

	// " "
	// x y z x y z
	for j < len(s) {
		// is there duplicate?

		index, ok := charToIndex[s[j]]
		// found duplication
		if ok {
			// calculate first, then set

			i = index + 1
			j = i
			// from i to j inclusive in a 0 indexed slice

			charToIndex = make(map[byte]int)
		} else {
			result = max(result, (j-i)+1)
			charToIndex[s[j]] = j
			j++
		}
	}

	return result
}

// optimized
/*
 * YELLOW -> I had a solution but it was really inneficient due to the multiple allocations for the new map and the j that kept going backwards
 * Mistakes: I couldn't figure out how to optimize without help. Visualizing the window and figuring out that I didn't need to rescan every value to
 * calculate the length didn't come easily.
 *
 * SLIDING WINDOW IS NOT RESCAN!!! My bad solution restarted it and that shouldn't be the case.
 *
 * The fix is lazy deletion: after i jumps forward the map still holds chars left of i, but
 * they are harmless - `index >= i` recognises a stale entry instead of paying to delete it.
 * Same trick shows up in heap problems (Dijkstra, Task Scheduler): push duplicates, discard
 * the outdated ones on pop.
 *
 * State has to match the plan. char -> index means i jumps in one hop and nothing is deleted.
 * A set of in-window chars would mean i steps forward one at a time, deleting as it goes
 * (also O(n), amortised). The slow version mixed the two: stored indexes AND cleared the map.
 *
 * Time: O(n) - each char enters and leaves the window once, i and j only move forward.
 * Space: O(1) - map is keyed by byte over printable ASCII, so <= 95 entries whatever len(s) is.
 */
func lengthOfLongestSubstringOptimized(s string) int {
	// hashmap char -> index
	// sliding window
	// result state
	charToIndex := make(map[byte]int)
	var i, j int
	// z x y z x y z
	var result int

	// " "
	// abba
	for j < len(s) {
		// is there duplicate?

		index, ok := charToIndex[s[j]]
		// foudn duplication and INSIDE the window
		if ok && (index >= i) {
			// calculate first, then set

			i = index + 1
			// from i to j inclusive in a 0 indexed slice
		}
		result = max(result, (j-i)+1)
		charToIndex[s[j]] = j
		j++
	}

	return result
}
