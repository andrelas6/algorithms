package arrays_lists

import "testing"

// Spec under test: true when the two strings are at most ONE edit apart, where an
// edit is insert / remove / replace. Zero edits (identical) counts as true.
//
// Nothing in the statement mentions casing or ignoring characters, so every byte
// counts and 'a' and 'A' are different.
func TestOneAway(t *testing.T) {
	cases := []struct {
		in1, in2 string
		want     bool
	}{
		// samples from the statement
		{"pale", "ple", true},
		{"pale", "pale", true},
		{"pales", "pale", true},
		{"bale", "pale", true},
		{"pales", "bale", false},
		{"pale", "bake", false},

		// degenerate inputs
		{"", "", true},
		{"", "a", true},
		{"a", "", true},
		{"", "ab", false},
		{"ab", "", false},
		{"a", "a", true},
		{"a", "b", true},
		{"a", "ab", true},
		{"ab", "a", true},

		// length differs by more than one -> always false
		{"abc", "abcde", false},
		{"abcde", "abc", false},
		{"a", "abcd", false},

		// replace, same length
		{"abc", "abc", true},
		{"abc", "abd", true},
		{"abc", "axc", true},
		{"xbc", "abc", true},
		{"abc", "xyc", false},
		{"abc", "xyz", false},

		// remove, in every position
		{"abc", "bc", true},
		{"abc", "ac", true},
		{"abc", "ab", true},

		// insert, in every position
		{"bc", "abc", true},
		{"ac", "abc", true},
		{"ab", "abc", true},

		// repeated characters: the skip must not get confused
		{"aaa", "aa", true},
		{"aa", "aaa", true},
		{"aab", "ab", true},
		{"ab", "aab", true},
		{"aaa", "aaa", true},
		{"aaa", "aab", true},
		{"aaa", "abb", false},

		// the mismatch is not at the first character
		{"abcdef", "abcdf", true},
		{"abcdef", "abdef", true},
		{"abcde", "abXde", true},
		{"abcde", "abde", true},
		{"abcde", "abcd", true},

		// two edits that look like one until you keep walking
		{"abcd", "abde", false},
		{"cat", "cs", false},
		{"abcde", "abdef", false},

		// transposition is two edits, not one
		{"ab", "ba", false},
		{"abcd", "abdc", false},
		{"abc", "bca", false},

		// casing counts
		{"Pale", "pale", true},
		{"PALE", "pale", false},

		// digits, symbols and spaces are ordinary characters
		{"a1", "a2", true},
		{"a b", "ab", true},
		{"a b", "a  b", true},
		{"a!b", "a?b", true},
		{"a!b", "a?c", false},

		// longer inputs
		{"the quick brown fox", "the quick brown fx", true},
		{"the quick brown fox", "the quick brown fax", true},
		{"the quick brown fox", "the quick brown", false},
	}

	for _, c := range cases {
		if got := oneAway(c.in1, c.in2); got != c.want {
			t.Errorf("oneAway(%q, %q) = %v, want %v", c.in1, c.in2, got, c.want)
		}
	}
}
