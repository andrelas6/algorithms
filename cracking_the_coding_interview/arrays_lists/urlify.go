package arrays_lists

/*
 * DESCRIPTION
 *
 * URLify. Write a method to replaec all spaces in a string with %20. You may assume that the string has sufficient space at the end to hold the additional characters, and that you are given the "true" lenght of the string.
 *
 */

// input "something with spaces  " 23 -> "something%20with%20spaces%20%20"
// brute force
// loop over each char, if space -> replace indexed element %20, continue until end
// return string
//
// optimize
// if string is epmty, return empty
// something about this challenge does not click - I can build an array with 23 slots and add each
// goal seems to make operatoin in place, but not possible heere
// because ' ' this is one byte, while "%20" is three bytes. I have three runes against one for the space.
func urlify(wannaBeUrl string, l int) string {
	response := ""

	// "some thing"
	for i := 0; i < l; i++ {
		char := wannaBeUrl[i]
		// s == ' '? No, then add to ""
		// ..
		// ' ' == ' '? Yes, then add "%20" -> some%20thing
		if char == byte(' ') { // 1 byte
			response += "%20" // 3 bytes
		} else {
			response += string(char)
		}
	}

	return response
}

// urlifyBytes is the in-place version: the one the book is actually asking for.
//
// CONTRACT
//   - buf holds trueLength real bytes, followed by exactly enough padding to fit
//     the expansion: len(buf) == trueLength + 2*(number of spaces in buf[:trueLength])
//   - every ' ' in buf[:trueLength] becomes "%20", written into buf itself
//   - returns the populated prefix of buf; no new backing array is allocated
//
// CONSTRAINTS
//   - O(n) time, O(1) extra space
//   - no strings.Builder, no append, no second slice
//
// THE INSIGHT (why forward fails)
//
//	Going left to right, writing "%20" over a single ' ' clobbers the two bytes
//	after it -- bytes you have not read yet. Going right to left, the write index
//	starts at the end of the padding and the read index at trueLength-1, so the
//	writer stays ahead of the reader and never destroys unread input.
//
// SHAPE
//   - one index for reading, one for writing
//   - walk read backwards from trueLength-1 down to 0
//   - space  -> write '0', '2', '%' (mind the order when moving backwards)
//   - other  -> copy the byte across
//   - return the correct slice of buf at the end
//
// "some thing  " -> some%20thing%20%20" -> BIGGER
// brute force
// create a slice, add every byte, but when finding a space, add %, then 2, then 0, the continue
// that means O(n), but space is O(n) because of second array of bytes
//
// brute force 2
// create space count variable
// walk the whole byte array, for every space found, add 2 to the space count
// then create a new byte array with length equal to the first + space count
// walk the same array again
// add the byte if it's not a space to the position
// when encountering space, add %, 2 and 0 to the next 3 positions, advance the loop 3 slots
// "some thing" -> walk s, o, m and e, then add %, 2 and 0, then counter goes from 5 to 8, continue
// O(n) time but O(1) space as well but overwriting chars will be a problem.
//
// Optimization
// walk the whole array from end to beginning (considering we are not traversing right-to-left languages)
// "some thing"
func urlifyBytes(buf []byte, trueLength int) []byte {
	spaceCount := 0

	for i := 0; i < trueLength; i++ {
		if buf[i] == ' ' {
			spaceCount += 2
		}
	}

	// "something   " -> 10 -> 12 -> i[10] == b[12]
	// "somethin g " ->
	lastBufElIndex := trueLength - 1
	offsetIndex := lastBufElIndex + spaceCount
	for i := lastBufElIndex; i >= 0; i-- {
		// i[10] -> " " -> offsetIndex = 14 -> i[14] = 0, i[13] = 2, i[12] = 0
		// i[9] -> g -> offsetIndex = 9 -> i[11] = g -> all good, i = 8, offsetIndex = 10
		// i[8] -> " " -> offsetIndex = 10 -> i[10] = 0, i[9] = 2, i[8] = 0, i = 7, offsetIndex = 7
		char := buf[i]
		if char == ' ' {
			buf[offsetIndex] = '0'
			buf[offsetIndex-1] = '2'
			buf[offsetIndex-2] = '%'

			offsetIndex -= 3
		} else {
			buf[offsetIndex] = char
			offsetIndex -= 1
		}
	}

	return buf[:(spaceCount + trueLength)]
}
