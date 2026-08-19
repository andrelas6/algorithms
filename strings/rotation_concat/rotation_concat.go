package rotationconcat

import "strings"

func isNonTrivialRotation(s1 string, s2 string) bool {
	// Write your code here
	// a b c d e a b c d e
	//   - - - - - b c d e a
	//    - - - - - c d e a b
	//     - - - - - d e a b c
	//      - - - - - e a b c d

	// constraints:
	// len s1 == len s2
	// s1 != s2

	if len(s1) != len(s2) {
		return false
	}

	if s1 == s2 {
		return false
	}

	// build s1 + s1
	// check if s2 is in (s1 + s1)

	s1s1 := s1 + s1

	// var windowStart int = 1
	// var windowEnd = windowStart + len(s1)
	// "ab" "ba"
	// p1 = 0 p2 = 1 s1s1 = "abab"
	// s1s1[1:2] -> "ba"
	// "ba" == "ba" true
	// "abcde" "cdeab"
	// p1 = 1 p2 = 5 s1s1 = abcdeabcde
	// substring = abcdeabcde[1:6] -> bcdea
	// bcdea == cdeab false
	// p1 = 2 p2 = 7 substring abcdeabcde[2:7] -> cdeab
	// true

	// "ab" "ba"
	// p1 = 1 p2 = 3 substring = abab[1:3], len(s1s1) - 1 = 4 - 1 = 3
	// ba == ba true

	return strings.Contains(s1s1, s2)
	// or
	// for windowEnd < len(s1s1) {
	// 	// grab the substring
	// 	// compare both
	// 	// include 0, exclude
	// 	substring := s1s1[windowStart:windowEnd]
	// 	if substring == s2 {
	// 		return true
	// 	}
	// 	// this is infinite, need to walk!
	// 	windowStart++
	// 	windowEnd++
	// }

	// return false
}
