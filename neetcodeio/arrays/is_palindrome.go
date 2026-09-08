package arrays

/*
 * 1st cycle
 * GREEN
 * mistake: forgot during implemntation to only apply lowercase to chars
 * time complexity: O(log n)
 * space complexity: O(1)
 */
func isPalindrome(s string) bool {
	// anana
	// nana
	// ignore case
	// ignore non alphanumrec chars
	// two pointer, compare

	var i, j int
	i = 0
	j = len(s) - 1

	// while the pointers are still apart
	for i < j {
		// Question: why do I need to convert?
		// first, check pointers vailidity
		if !isValueValid(i, s) {
			i++
			continue
		}

		if !isValueValid(j, s) {
			j--
			continue
		}

		charI := s[i]
		charJ := s[j]

		if lowercase(charI) != lowercase(charJ) {
			return false
		}

		i++
		j--
	}

	return true
}

func isValueValid(pointer int, s string) bool {
	val := s[pointer]
	if (val >= 'a' && val <= 'z') ||
		(val >= 'A' && val <= 'Z') ||
		(val >= '0' && val <= '9') {
		return true
	}

	return false
}

func lowercaseWithBytes(char byte) byte {
	if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
		lowercaseOffset := 'a' - 'A'

		// converts to lowercase
		// convert to uppercase would be char &^ lowecaseOffset
		return char | byte(lowercaseOffset)
	}

	return char
}

func lowercase(char byte) byte {
	if char >= 'A' && char <= 'Z' {
		return char + ('a' - 'A')
	}

	return char
}
