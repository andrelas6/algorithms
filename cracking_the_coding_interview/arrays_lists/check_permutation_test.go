package arrays_lists

import "testing"

// Spec under test, from the notes at the top of check_permutation.go:
//   - every character counts, symbols and digits included
//   - no case-folding is described, so 'a' and 'A' are distinct characters
//   - permutation is order-insensitive but multiplicity-sensitive:
//     "aab" is not a permutation of "abb"
func TestCheckPermutation(t *testing.T) {
	cases := []struct {
		in1, in2 string
		want     bool
	}{
		// samples from the statement
		{"abcd", "cbda", true},
		{"abc", "abd", false},
		{"", "ab", false},

		// walk-through cases from the notes
		{"abca", "bcaa", true},
		{"abcd", "bcad", true},

		// degenerate inputs
		{"", "", true},
		{"a", "a", true},
		{"a", "b", false},
		{"a", "", false},

		// identical strings are permutations of themselves
		{"abc", "abc", true},
		{"aaaa", "aaaa", true},

		// length mismatch -> early exit
		{"abc", "ab", false},
		{"ab", "abc", false},
		{"aab", "aabb", false},

		// same length, same character set, different multiplicity
		{"aab", "abb", false},
		{"aabb", "aaab", false},
		{"aaab", "abbb", false},

		// same length, disjoint character sets
		{"abc", "xyz", false},
		{"aaa", "bbb", false},

		// reversal is a permutation
		{"abcdef", "fedcba", true},
		{"ab", "ba", true},

		// case sensitivity: upper and lower are distinct
		{"aA", "Aa", true},
		{"ab", "AB", false},
		{"Abc", "abC", false},
		{"aAbB", "BbAa", true},

		// symbols and digits count, they are not ignored
		{"a1b2", "2b1a", true},
		{"a1b2", "ab12", true},
		{"%$^", "^$%", true},
		{"a!b", "ab!", true},
		{"a!b", "ab?", false},
		{"123", "321", true},
		{"123", "312", true},
		{"1234", "1235", false},

		// symbol vs letter of the same length
		{"!!!", "aaa", false},
		{"a b", "b a", true},

		// whitespace is a character like any other
		{"ab ", "ab", false},
		{" ab", "ab ", true},

		// longer inputs
		{"thequickbrownfox", "xofnworbkciuqeht", true},
		{"thequickbrownfox", "xofnworbkciuqehz", false},

		// repeated-heavy inputs
		{"aabbccdd", "dcbadcba", true},
		{"aabbccdd", "aabbccde", false},

		// multi-byte runes: comparison is per rune, not per byte
		{"éa", "aé", true},
		{"éé", "éa", false},
		{"日本", "本日", true},
		{"日本", "日日", false},
	}

	for _, c := range cases {
		if got := checkPermutation(c.in1, c.in2); got != c.want {
			t.Errorf("checkPermutation(%q, %q) = %v, want %v", c.in1, c.in2, got, c.want)
		}
	}
}
