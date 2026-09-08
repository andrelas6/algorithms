package hashmaps

func groupAnagramsImproved(strs []string) [][]string {
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

	visitedMap := make(map[string]struct{})

	// somehow build the signature that turns ACT and CAT into the same thing
	wordToSig := make(map[string][26]int)

	for _, word := range strs {
		sig := getSigImproved(word, alphabetOrderMap)
		wordToSig[word] = sig
	}

	for i, word := range strs {
		if _, visited := visitedMap[word]; visited {
			continue
		} else {
			visitedMap[word] = struct{}{}
		}

		sigA := wordToSig[word]
		anagram := make([]string, 0)
		anagram = append(anagram, word)
		for j := i + 1; j < len(strs); j++ {
			otherWord := strs[j]
			sigB := wordToSig[otherWord]
			if sigA == sigB {
				anagram = append(anagram, otherWord)
				visitedMap[otherWord] = struct{}{}
			}
		}

		result = append(result, anagram)
	}

	return result
}

func getSigImproved(word string, alphabetMap map[byte]int) [26]int {
	var charCountArr [26]int
	for i, _ := range word {
		byteChar := word[i]
		position, _ := alphabetMap[byteChar]

		charCountArr[position]++
	}
	return charCountArr
}

// improvement next:
// no inner loop
