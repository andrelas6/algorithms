package arrays_lists

import "fmt"

/*
 * DESCRIPTION
 *
 * String compression:
 * Implement a method to perform basic string compression using the counts of repeated characters.
 * For example, the string aabcccccaaa would become a2b1c5a3.
 * If the "compressed" string would not become smaller than te original string,
 * your method should return the original string. You can assume the string has only uppercase and lowercase letters (a - z).
 */

// Example
// aabcccccaaa -> a2b1c5a3 -> SMALLER RETURN
// aabb -> a2b2 -> aabb
// ab -> a1b1 -> ab
// aaaaaaaaaaaa -> a10 -> a10
//
// Brute force
// create new string
// loop on given string -> get the element as (FOUND) -> compare with next, if same, ++, otherwise finish that found
// do the same for the next found
// as the loop moves, keep building the return string
// compare lengths of both strings, then return
//
// "aabb"
// found = nil
// count = 0
// "a" != null -> found = a, count = 1
// "a" != null a == a -> found == a, count = 2
// null check "a" != "b" ->  newStr = "a2" -> count = 1
// null check "b" != "b" ->  count = 2
// end of loop -> newStr = "a2b2"
// a2b2 > aabb ? aabb : a2b2
//
// aaaaaaaaa
// a != nil -> found = a, count = 1
// a == a -> count = 2
// ...
// a == a -> count = 10
// end of loop -> newStr = "found+counter"
// a10 < aaaaaaa ? a10 : aaaaaaa

// aaaaaaaa
func stringCompression(str string) string {
	var letter rune
	var count int
	var newStr string

	if len(str) == 0 {
		return ""
	}

	letter = rune(str[0])

	// aabb
	// a2b1
	for _, s := range str {
		// assign
		if letter != s {
			newStr += fmt.Sprintf("%c%d", letter, count)
			letter = s
			count = 0
		}
		count++
	}

	newStr += fmt.Sprintf("%c%d", letter, count)

	if len(newStr) < len(str) {
		return newStr
	} else {
		return str
	}
}
