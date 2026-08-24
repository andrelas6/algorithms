package arrays_lists

import "testing"

// Spec under test:
//   - compress CONSECUTIVE runs: "aabcccccaaa" -> "a2b1c5a3" (note 'a' appears
//     twice, as two separate runs -- this is run-length encoding, not frequency
//     counting)
//   - every run gets a count, including runs of length 1
//   - if the compressed form is NOT smaller than the original, return the original
//     (equal length counts as "not smaller")
//
// The statement says the input is only letters a-z / A-Z, so there are no digit or
// symbol cases here -- that behaviour is undefined by the problem.
func TestStringCompression(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		// the sample from the statement
		{"aabcccccaaa", "a2b1c5a3"},

		// degenerate inputs
		{"", ""},
		{"a", "a"},
		{"aa", "aa"},

		// the shrink threshold: "aaa" -> "a3" is 2 < 3, so it wins
		{"aaa", "a3"},
		{"aaaa", "a4"},

		// compressed is the SAME length -> return the original
		{"aabb", "aabb"},
		{"aabbcc", "aabbcc"},
		{"aabbccdd", "aabbccdd"},
		{"abbb", "abbb"},

		// compressed is LONGER -> return the original
		{"ab", "ab"},
		{"abc", "abc"},
		{"aab", "aab"},
		{"ababab", "ababab"},
		{"abcdef", "abcdef"},

		// compressed is genuinely smaller
		{"aaabbb", "a3b3"},
		{"aaabbbccc", "a3b3c3"},
		{"aabbb", "a2b3"},
		{"abbbb", "a1b4"},
		{"abbbbbbc", "a1b6c1"},

		// counts of two digits
		{"aaaaaaaaaa", "a10"},
		{"aaaaaaaaaaaa", "a12"},
		{"aaaaaaaaab", "a9b1"},

		// run at the very start / very end must not be dropped
		{"aaaab", "a4b1"},
		{"abbbb", "a1b4"},
		{"aaaabbbb", "a4b4"},

		// casing: 'a' and 'A' are different characters, so they are different runs
		{"aA", "aA"},
		{"aaAA", "aaAA"},
		{"aaaAAA", "a3A3"},
		{"aAaA", "aAaA"},

		// longer mixed input
		{"aabbccccccdd", "a2b2c6d2"},
		{"wwwwaaadexxxxxxywww", "w4a3d1e1x6y1w3"},
	}

	for _, c := range cases {
		if got := stringCompression(c.in); got != c.want {
			t.Errorf("stringCompression(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
