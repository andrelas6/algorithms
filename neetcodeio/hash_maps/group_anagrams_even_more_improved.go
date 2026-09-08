package hashmaps

func groupAnagramsEvenMoreImproved(strs []string) [][]string {
	// loop over str
	// loop over next str
	// use hash table for character matching
	// once a string is paired with another, it shouldn't be checked again
	alphabetOrderMap := make(map[byte]int)
	alphabetOrderMap['a'] = 0

	for i := 1; i < 26; i++ {
		char := byte('a') + byte(i)
		alphabetOrderMap[char] = i
	}

	var result [][]string

	// visitedMap := make(map[string]struct{})

	// somehow build the signature that turns ACT and CAT into the same thing
	sigToWord := make(map[[26]int][]string)

	for _, word := range strs {
		sig := getSig(word, alphabetOrderMap)
		anagram, ok := sigToWord[sig]

		if ok {
			sigToWord[sig] = append(anagram, word)
		} else {
			sigToWord[sig] = []string{word}
		}
	}

	// convert
	for _, anagram := range sigToWord {
		result = append(result, anagram)
	}

	return result
}

func getSig(word string, alphabetMap map[byte]int) [26]int {
	var charCountArr [26]int
	for i, _ := range word {
		byteChar := word[i]
		position, _ := alphabetMap[byteChar]

		charCountArr[position]++
	}
	return charCountArr
}
