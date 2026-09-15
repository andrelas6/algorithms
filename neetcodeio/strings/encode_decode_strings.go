package strings

import (
	"strconv"
	"strings"
)

type Solution struct{}

/*
 * YELLOW: I had a solution but I overcomplicated too much. I had in my mind that I needed to use go without any imported libraries
 * from std lib, which is incorrect. I can use them apparently in those challenges. It's just important to explain their time complexity.
 * Time complexity encode: O(m) m is chars in all strings - which is what is iterated over.
 * Space encode: O(m) the chars in all strings which becomes the encoded string
 * Time complexity decode: O(m), where m is the chars in the encoded string. Append and convertion are O(1).
 * Space: O(m), where n is amount of chars - because since I am substringing the encoded string, the encoded string is kept in memory. Thus O(m).
 */
func (s *Solution) Encode(strs []string) (encoded string) {
	var encodedBuilder strings.Builder
	for _, str := range strs {
		length := len(str)
		encodedBuilder.WriteString(strconv.Itoa(length))
		encodedBuilder.WriteByte('#')
		encodedBuilder.WriteString(str)
	}

	return encodedBuilder.String()
}

func (s *Solution) Decode(encoded string) (decoded []string) {
	i := 0

	for i < len(encoded) {
		j := i
		for encoded[j] != '#' {
			j++
		}

		encodedPartLength, _ := strconv.Atoi(encoded[i:j])
		start := j + 1
		decoded = append(decoded, encoded[start:start+encodedPartLength])
		i = (start + encodedPartLength)
	}

	return decoded
}
