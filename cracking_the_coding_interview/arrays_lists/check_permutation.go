package arrays_lists

/* DESCRIPTOIN
 * 2. Check permuntation: given two strings, write a method to decide if one is a permutation of the other
 */

// Example
// "abcd" and "cbda" -> are permutation of the other
// "abc" and "abd" -> not a permutation
// "" and "ab" -> not a permutatoin
// again, symbols and numbers count? probably
//
// Brute Force
// hash table -> store each rune and their count
// loop over other string -> reduce each element from hash table
// if every value of hash table is zero, is permutation
// otherwise, no
// three loops -> one to store each unique value, other to check, other to count
//
// Optimize
//
// check len of both -> must be equal -> not, early exit
// use two loops -> one to store, the other to check
// second loop -> remove values as they become zero
// check hash table length
// "abcd" and "cbda"
// { a: 1, b: 1, c: 1, d: 1 }, then
// 1st -> "c", c gets 0, then removed, then { a: 1, b: 1, d: 1 }
// ...
// Nth -> "d", d gets 0, then removed, then { }
// check length with len() -> O(1) time and O(1) space

// LEBOWIT -> Listen, Example, Brute Force, Optimize, Walk through, Implement and Test
// 1st "abca" and "bcaa"
// 2nd "abcd" and "bcad"
func checkPermutation(str1, str2 string) bool {
	// this might not work 100% because some strings have one element but length 2 i.e. "é" but I don't worry about that here
	// but if I have to, I would use a rune count instead, which might cover more but not 100% i.e. "é" == "é" might be false becuase of the second byte (an accent thing)
	if len(str1) != len(str2) {
		return false
	}

	charMap := make(map[rune]int)
	// char count
	for _, s1 := range str1 {
		v, ok := charMap[s1]

		// "a" ? no, then charMap.a = 1
		// "b" ? no, then charMap.b = 1
		// ..
		// "a" ? yes, charMap.a = 2
		if ok {
			charMap[s1] = v + 1
		} else {
			charMap[s1] = 1
		}
	}

	// { a: 2, b: 1, c: 1 }
	for _, s2 := range str2 {
		v, ok := charMap[s2]

		// "t" ? no
		if ok {
			// "b" == 1? yes, then delete
			// "c" == 1? yes, then delete
			// "a" == 1? yes, then delete
			if v == 1 {
				delete(charMap, s2)
			} else {
				// "aa" == 1? no, then - 1, then { a: 1 }
				charMap[s2] = v - 1
			}
		} else {
			return false
		}
	}

	// { } -> len 0

	return len(charMap) == 0
}
