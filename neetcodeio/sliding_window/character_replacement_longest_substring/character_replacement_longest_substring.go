package characterreplacementlongestsubstring

import "math"

func characterReplacement(s string, k int) int {
	// use sliding window
	// within a sliding window, find if there is a longest valid substring given k
	// longest valid substring: substring length - most common character occurrence <= k
	// if condition matches, move right edge, left edge stays
	// if condition does not match, move left and right edges
	// when right edge goes out of bounds, stop and return

	var left, right int
	var result int
	var charToCount map[byte]int = make(map[byte]int)

	for right < len(s) {
		// is current substring valid?
		mostCommonCount := getMostCommonCharCount(s[right], charToCount)
		if (right-left)+1-mostCommonCount <= k {
			result = max(right-left+1, result)
		} else {
			charToCount[s[left]]--
			left++
		}
		right++
	}

	return result
}

// if right only moves, great
// if left move too?
func getMostCommonCharCount(char byte, charToCount map[byte]int) int {
	charToCount[char]++
	var actualMax = math.MinInt
	for _, v := range charToCount {
		actualMax = max(actualMax, v)
	}
	return actualMax
}
