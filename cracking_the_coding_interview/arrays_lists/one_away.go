package arrays_lists

/*
 * DESCRIPTION
 *
 * One Away: There are three types of edits that can be performed on strings:
 * insert a character, remove a character, or replace a character.
 * Given two strings, write a functoin to check if they are one edit (or zero edits) away.
 *
 * example: pale, ple -> true
 * example: pale, pale -> true
 * example: pales, pale -> true
 * example: bale, pale -> true
 * example: pales, bale -> false
 */

// Brute force
// length diff can't be bigger than 1
// they should be equal sequence wise except for the change
// p a l e -> this walks faster because it's bigger
// p   l e -> this walks slower because it's equal size
// p a l e -> this walks the same as the other one
// b a l e ->
// optimize
// if length is equal, no need for tracking two values
// if length diff is > 1, return early
// distinguishing insertion, deletion or replacement is not relevant here
//
// walk through
//
// pointer 1, pointer 2
// create diffCounter
// loop
// p a l e -> p1 = p, p1 = a (diff) -> p1 = l (walks faster), p2 = e
// p l e -> p2 = p, p2 = l (diff), p2 = l, p2 = e
// if diffCounter > 1, return false
// otherwise, return true
//
// loop
// p a l e -> p1 = p (diff), ..., p1 = ends (nowhere to go)
// b a l e s -> p2 = b, ..., p2 = s
// if diffCounter > 1, return false
// return false, diff (p and b) an diff (diff length)
//
// optmize 2
// if length 1 != length 2 -> there can be no diffs at all (earlier return)
// if length 1 == length 2 -> there can be only one diff

func oneAway(str1, str2 string) bool {
	var diffCount int

	var lengthDiff = len(str1) - len(str2)
	if lengthDiff > 1 || lengthDiff < -1 {
		return false
	}

	len1 := len(str1)
	len2 := len(str2)
	longer := str1
	shorter := str2
	if len1 < len2 {
		shorter = str1
		longer = str2
	}
	// p a l e
	// p l e
	// i = 0, j = 0, p == p -> i = 1, j = 1, diffCount = 0
	// i = 1, j = 1, a == l -> i = 2, j = 1, diffCount = 1
	// i = 2, j = 1, l == l -> i = 3, j = 2, diffCount = 1
	// i = 3, j = 2, e == e -> i == 4, j == 3, diffCount = 1
	// i == 4 -> break look
	//
	// p a l e s
	// p a l e
	// i = 0, j = 0, p == p -> i == 1, j = 1, diffCount = 0
	// i 1, j 1, a == a
	// ...
	// i = 4, j
	//
	//
	// bales
	// ales
	// i = 0, j = 0, b == p -> diffCount = 1, i = 1, j = 0
	// i = 1, j = 0, a == p -> diffCount = 2, i = 2, j = -1
	//
	var j int // smaller
	for i := 0; i < len(longer); i++ {
		if j >= len(shorter) {
			break
		}
		if longer[i] != shorter[j] {
			diffCount++
			if len1 != len2 {
				j--
			}
		}

		if diffCount > 1 {
			return false
		}
		j++
	}

	return true

}
