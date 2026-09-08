package arrays

func isAnagram(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}

	charCountS := make(map[rune]int)

	for i := 0; i < len(s); i++ {
		char := rune(s[i])
		charCountS[char]++
	}

	// now delete every char from t and if t is empty, return true

	for i := 0; i < len(t); i++ {
		char := rune(t[i])
		val, ok := charCountS[char]

		if !ok {
			return false
		} else {
			if val == 1 {
				delete(charCountS, char)
			} else {
				charCountS[char] = val - 1
			}
		}
	}

	return len(charCountS) == 0
}
