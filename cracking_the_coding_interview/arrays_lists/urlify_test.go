package arrays_lists

import "testing"

// Spec under test, from the statement and the notes at the top of urlify.go:
//   - every ' ' within the true length becomes "%20"
//   - l is the "true" length: anything past it is padding and must be dropped
//   - all other characters pass through untouched
func TestUrlify(t *testing.T) {
	cases := []struct {
		in   string
		l    int
		want string
	}{
		// canonical CTCI sample
		{"Mr John Smith    ", 13, "Mr%20John%20Smith"},

		// sample from the notes
		{"something with spaces  ", 23, "something%20with%20spaces%20%20"},

		// degenerate inputs
		{"", 0, ""},
		{"a", 1, "a"},
		{" ", 1, "%20"},

		// true length is 0 -> everything is padding
		{"   ", 0, ""},
		{"abc", 0, ""},

		// no spaces at all -> string is unchanged
		{"abc", 3, "abc"},
		{"nospaceshere", 12, "nospaceshere"},

		// single space in each position
		{" ab", 3, "%20ab"},
		{"a b", 3, "a%20b"},
		{"ab ", 3, "ab%20"},

		// consecutive spaces are each replaced
		{"a  b", 4, "a%20%20b"},
		{"  ab", 4, "%20%20ab"},
		{"ab  ", 4, "ab%20%20"},

		// nothing but spaces, all inside the true length
		{"  ", 2, "%20%20"},
		{"   ", 3, "%20%20%20"},

		// the padding past l must be dropped, spaces included
		{"ab   ", 2, "ab"},
		{"a b     ", 3, "a%20b"},
		{"Mr John Smith    ", 7, "Mr%20John"},
		{"Mr John Smith    ", 8, "Mr%20John%20"},

		// l stopping mid-word
		{"hello world", 5, "hello"},
		{"hello world", 6, "hello%20"},
		{"hello world", 7, "hello%20w"},

		// digits and symbols pass through untouched
		{"a1 b2", 5, "a1%20b2"},
		{"%20 x", 5, "%20%20x"},
		{"a!b c?d", 7, "a!b%20c?d"},
		{"1 2 3", 5, "1%202%203"},

		// other whitespace is not a space and must not be replaced
		{"a\tb", 3, "a\tb"},
		{"a\nb", 3, "a\nb"},
		{"a\tb c", 5, "a\tb%20c"},

		// longer input, many replacements
		{"the quick brown fox", 19, "the%20quick%20brown%20fox"},
	}

	for _, c := range cases {
		if got := urlify(c.in, c.l); got != c.want {
			t.Errorf("urlify(%q, %d) = %q, want %q", c.in, c.l, got, c.want)
		}
	}
}
