package arrays_lists

/*
 * DESCRIPTION
 *
 * Palindrome Permutation: Given a string, write a functoin to check
 * if it's a permutation of a palindrome.
 * The palindrome does not need to be limited to just dictionary words. You can ignore casing and non-letter characters.
 *
 * Input: Tact coa
 * Output: True (perms: taco cat, atco cta)
 */

// Example
// Tact coa -> taco cat
// brute force
// How can tact coa be a palindrome? character count most likely
// if aa, tt, cc, o -> palindrome. every letter has a even count except max 1
// brute force
// use a map for the char repeat count O(1) access and walk over the string
// if length == 0 -> false
// "tact coa" ->
// "t" -> { t : 1 }, { t:1, a:1, c:1 } -> { t: 2, a:1, c:1 } -> { t: 2, c: 2, a:2, o: 1 }
// if two elements have odd count, return false
// otherwise, return true
//
// loop -> one map for the count
// loop again -> check how may times odd counts appear
// early return
// no optimization yet
//
// walk through
// "tact coa"
// if letter: "t" -> { t : 1 }, { t:1, a:1, c:1 } -> { t: 2, a:1, c:1 } -> { t: 2, c: 2, a:2, o: 1 }
// second loop -> walk the map! check each element and if count % 2 != 0 two times, return false
// otherwise return true

func palindromePermutation(str string) bool {

	charCount := make(map[rune]int)
	charDiff := 'a' - 'A'

	// first, get the count
	for _, s := range str {
		s = s | charDiff

		if s >= 'a' && s <= 'z' {
			v, ok := charCount[s]

			if ok {
				charCount[s] = v + 1
			} else {
				charCount[s] = 1
			}
		}
	}

	// find out if more than 2 times odd count
	//
	oddCount := 0
	for _, v := range charCount {
		if v%2 != 0 {
			oddCount += 1
		}

		if oddCount > 1 {
			return false
		}
	}

	return true
}
