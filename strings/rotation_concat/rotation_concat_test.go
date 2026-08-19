package rotationconcat

import "testing"

func TestIsNonTrivialRotation(t *testing.T) {
	cases := []struct {
		s1, s2 string
		want   bool
	}{
		// samples from the statement
		{"abcde", "cdeab", true},
		{"a", "a", false},
		{"a", "b", false},

		// every non-trivial rotation of the same string
		{"abcde", "bcdea", true},
		{"abcde", "deabc", true},
		{"abcde", "eabcd", true},
		{"abcde", "abcde", false}, // identity

		// shortest non-trivial case
		{"ab", "ba", true},
		{"ab", "ab", false},

		// periodic strings: a NON-ZERO rotation can equal s1, so the
		// identity check must run before the window scan
		{"abab", "abab", false},
		{"abab", "baba", true},
		{"abcabc", "abcabc", false},
		{"abcabc", "bcabca", true},
		{"abcabc", "cabcab", true},

		// all-identical: every rotation is the identity
		{"aa", "aa", false},
		{"aaaa", "aaaa", false},

		// same multiset of letters, not a rotation
		{"abc", "acb", false},
		{"aaab", "aaba", true},
		{"aaab", "baaa", true},

		// length mismatch (excluded by constraints, but the trick depends on it)
		{"abc", "ab", false},
	}

	for _, c := range cases {
		if got := isNonTrivialRotation(c.s1, c.s2); got != c.want {
			t.Errorf("isNonTrivialRotation(%q, %q) = %v, want %v", c.s1, c.s2, got, c.want)
		}
	}
}
