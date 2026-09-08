package hashmaps

/*
 * YELLOW -> brute force works but is terribly inneficient
 *
 * Time: O(n * n^2)
 */

func groupAnagrams(strs []string) [][]string {
	// loop over str
	// loop over next str
	// use hash table for character matching
	// once a string is paired with another, it shouldn't be checked again

	wordToCharCount := make(map[string]map[rune]int)
	var result [][]string

	// calculate char count

	for _, str := range strs {
		wordToCharCount[str] = getCharsCountMap(str)
	}

	inAnagramGroup := make(map[string]struct{})

	// "act","pots","tops", "cat"
	for i, str := range strs {
		if _, isInGroup := inAnagramGroup[str]; isInGroup {
			// don't check it again
			continue
		} else {
			inAnagramGroup[str] = struct{}{}
		}
		anagram := []string{str}
		// go over all other strings
		for j := i + 1; j < len(strs); j++ {
			otherStr := strs[j]
			if len(str) == len(otherStr) && isAnagram(wordToCharCount[str], wordToCharCount[otherStr]) {
				anagram = append(anagram, otherStr)
				inAnagramGroup[otherStr] = struct{}{}
			}
		}

		result = append(result, anagram)
	}

	return result
}

func isAnagram(aCharCount, bCharCount map[rune]int) bool {
	for char, bCount := range bCharCount {
		aCount, ok := aCharCount[char]

		if !ok {
			return false
		} else {
			if aCount != bCount {
				return false
			}
		}
	}

	return true
}

func getCharsCountMap(str string) map[rune]int {
	charCount := make(map[rune]int)
	for _, c := range str {
		_, ok := charCount[c]

		if ok {
			charCount[c] += 1
		} else {
			// charCount[c]++ maybe works and I don't need if/else
			charCount[c] = 1
		}
	}

	return charCount
}
