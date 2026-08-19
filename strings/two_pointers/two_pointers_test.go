package twopointers

import "testing"

func TestIsAlphabeticPalindrome(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		// samples from the statement
		{"A1b2B!a", true},
		{"Z", true},
		{"abc123cba", true},

		// degenerate inputs
		{"", true},
		{"a", true},
		{"1", true},
		{"aa", true},
		{"ab", false},

		// no letters at all -> vacuously a palindrome
		{"!!!", true},
		{"123", true},
		{"$%^", true},

		// case insensitivity
		{"Zz", true},
		{"aXbbxA", true},

		// pointers must skip non-letters on both sides
		{"a!b", false},
		{"!ab!", false},
		{"!a!", true},
		{"1a2a3", true},
		{"ab!!ba", true},
		{"!@#a1a#@!", true},

		// odd/even length palindromes
		{"aba", true},
		{"abba", true},

		// near misses
		{"abca", false},
		{"a1b2c", false},
	}

	for _, c := range cases {
		if got := isAlphabeticPalindrome(c.in); got != c.want {
			t.Errorf("isAlphabeticPalindrome(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}
