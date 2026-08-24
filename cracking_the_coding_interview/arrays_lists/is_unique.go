package arrays_lists

/* DESCRIPTION
 *
 * Is Unique: implement an algorith to determine fia string has all unique characters. What if you cannot use additoinal data structures?
 *
 */

// INTERVIEW notes
// str -> all unique chasr in there?
// array_lists -> is every char unique there?
// "12313231%%%$$^^" -> should consider symbols and numbers?o
// assumption: symbos and ints are not valid

// loop over every char of the list
// get the char, compare with all previous "walked" chars
// comparison -> is there already -> return false
// comparison -> is not htere? -> store and continue iterating
// design decisions
// store: how? array -> comparison means another walk. Hash tables -> comparison is O(1)
//
// optimized version
//
// loop over every char of the list
// get the char, compare with all previous "walked" chars
// comparison -> check hash table -> is there? return false
// comparison -> is not in hash table -> store in hash table and continue iterating
// O(n) time and O(n) space
func isUnique(str string) bool {
	walkedChars := make(map[rune]bool)

	for _, s := range str {
		if (s >= 'a' && s <= 'z') || (s >= 'A' && s <= 'Z') {
			_, ok := walkedChars[s]

			if ok {
				return false
			}

			walkedChars[s] = true
		}
	}

	return true
}
