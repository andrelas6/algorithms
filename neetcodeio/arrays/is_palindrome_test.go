package arrays

import "testing"

func TestIsPalindrome(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		// samples
		{"classic", "A man, a plan, a canal: Panama", true},
		{"not a palindrome", "race a car", false},
		{"was it a car", "Was it a car or a cat I saw?", true},

		// degenerate
		{"empty", "", true},
		{"single letter", "a", true},
		{"single space", " ", true},
		{"only punctuation", ".,!?", true},

		// digits count as alphanumeric
		{"digits palindrome", "12321", true},
		{"digits not", "12345", false},
		{"letter vs digit", "0P", false},
		{"mixed alnum", "a1b2b1a", true},

		// case insensitivity
		{"upper/lower", "Aa", true},
		{"AbBa", "AbBa", true},

		// punctuation on the outside / inside
		{"leading junk", "!!aba!!", true},
		{"junk in the middle", "ab!!ba", true},
		{"junk hides nothing", "a!b", false},

		// even / odd length
		{"even", "abba", true},
		{"odd", "aba", true},
		{"near miss", "abca", false},
	}

	for _, c := range cases {
		if got := isPalindrome(c.in); got != c.want {
			t.Errorf("%s: isPalindrome(%q) = %v, want %v", c.name, c.in, got, c.want)
		}
	}
}

func TestLowercaseHelpers(t *testing.T) {
	// both helpers must agree on every byte they claim to handle
	for c := 0; c < 128; c++ {
		b := byte(c)
		if lowercase(b) != lowercaseWithBytes(b) {
			t.Errorf("byte %d (%q): lowercase = %q, lowercaseWithBytes = %q",
				c, b, lowercase(b), lowercaseWithBytes(b))
		}
	}
}
