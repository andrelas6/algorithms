package twopointers

func isAlphabeticPalindrome(code string) bool {
	// Write your code here
	// text: a3b%sc -> textCharsOnly -> isPalindrome?
	// uppercase, lowercase, 0-9, symbols
	// 0 - 1000 chars
	// exclude some ASCII codes from it -> what is below 33 and above 126? IDK maybe spaces?
	// only printable ASCII letters -> all good for Go's utf8
	// Input Z -> 1 is palin
	// Input abc123cba -> is palin
	// empty string? should be 0?
	// only numbers? should be 0?
	// only symbols? should be 0?\
	// Z123 is palin?
	// Zz is palindrome? yes it's case-insensitive

	// use two-pointer strategy -> front end end
	// ignore symbols -> it's like they don't have any index

	// A1b2B!a -> P1 A P2 a (true) -> P1 b P2 B (skip 1 and !)
	// A1b2Ba is essentially same
	// need to walk to next valid element (is char)
	// how do I verify that element is char? since code[i] is byte, I can do byte comparison

	var pointer2 = len(code) - 1 // lowercase won't guarantee same length
	for pointer1 := 0; pointer1 < pointer2; {
		letterStart := code[pointer1]
		letterEnd := code[pointer2]

		// a a
		// isStartValid = true isEndValid = true
		// p1 = 1 p2 = 7
		// b b
		// isStartValid = true isEndValid = true
		// p1 = 2 p2 = 6
		// c c
		// isStartValid = true isEndValid = true
		// p1 = 3 p2 = 5
		// 1 3
		// isStartValid = false isEndInvalid = false
		// p1 = 4 p2 4
		// comparison must be case insensitive
		var isStartValid = (letterStart >= 'a' && letterStart <= 'z') || (letterStart >= 'A' && letterStart <= 'Z')
		var isEndValid = (letterEnd >= 'a' && letterEnd <= 'z') || (letterEnd >= 'A' && letterEnd <= 'Z')

		if isStartValid && isEndValid {
			caseDiff := 'a' - 'A'
			lowercaseStart := letterStart | byte(caseDiff)
			lowercaseEnd := letterEnd | byte(caseDiff)
			if lowercaseStart != lowercaseEnd {
				return false
			} else {
				pointer1++
				pointer2--
			}
		}
		// validation -> can I compare these two? -> should be letters
		// will be false for digits and symbols
		if !isStartValid {
			pointer1++
		}

		if !isEndValid {
			pointer2--
		}
	}

	return true
}
