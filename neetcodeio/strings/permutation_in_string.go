package strings

/*
 * Yellow - I could figure out the solution with the sliding window and the map, but I was not sharp in the implementation
 * Time: O(26 * n), where n is the length of s2 and m is the alphabet
 * Space: O(26), because I keep two maps that are as big as the alphabet
 */
func checkInclusion(s1 string, s2 string) bool {
	// 1st map runes and their count in s2
	// use sliding window
	// left is 0, right is len(s1)
	// iterate with sliding window shifting
	// first iteration l = 0, r = 3, go to the end, add to a map, then another loop to compare
	// if not, l = 0 key is deleted, next is added (r = 4)
	// second iteration l = 1, r = 4
	// check the map
	// check matches

	// build matches and s1/s2 map for the first three characters

	s1CharToOccurrence := make(map[rune]int)
	s2CharToOccurrence := make(map[rune]int)

	if len(s1) > len(s2) {
		return false
	}

	for i, c := range s1 {
		s1CharToOccurrence[c]++

		s2Char := rune(s2[i])
		s2CharToOccurrence[s2Char]++
	}

	// sliding window of s1 length
	var left, right, matches int
	left = 0
	right = len(s1) - 1
	for right < len(s2) {
		for k, v1 := range s1CharToOccurrence {
			if v2, ok := s2CharToOccurrence[k]; ok && v1 == v2 {
				matches++
			}
		}

		if matches == len(s1CharToOccurrence) {
			return true
		}
		if right == len(s2)-1 {
			return false
		}
		matches = 0
		// delete seen element
		leftChar := rune(s2[left])
		leftVal, _ := s2CharToOccurrence[leftChar]
		if leftVal == 1 {
			delete(s2CharToOccurrence, leftChar)
		} else {
			s2CharToOccurrence[leftChar]--
		}

		// move window
		left++
		right++
		rightChar := rune(s2[right])
		s2CharToOccurrence[rightChar]++
	}

	return false
}
