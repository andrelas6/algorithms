package arrays_lists

import "testing"

// Every `in` below is already padded to its final size, exactly as the statement
// promises: len(in) == l + 2*(spaces in in[:l]). The padding bytes are spaces,
// which is what makes the problem interesting -- they look like input but are not.
func TestUrlifyBytes(t *testing.T) {
	cases := []struct {
		in   string
		l    int
		want string
	}{
		// canonical CTCI sample
		{"Mr John Smith    ", 13, "Mr%20John%20Smith"},

		// degenerate inputs
		{"", 0, ""},
		{"a", 1, "a"},
		{"   ", 1, "%20"},

		// no spaces at all -> bytes are unchanged
		{"abc", 3, "abc"},
		{"nospaceshere", 12, "nospaceshere"},

		// single space in each position
		{" ab  ", 3, "%20ab"},
		{"a b  ", 3, "a%20b"},
		{"ab   ", 3, "ab%20"},

		// consecutive spaces are each replaced
		{"a  b    ", 4, "a%20%20b"},
		{"  ab    ", 4, "%20%20ab"},
		{"ab      ", 4, "ab%20%20"},

		// nothing but spaces, all inside the true length
		{"      ", 2, "%20%20"},
		{"         ", 3, "%20%20%20"},

		// padding past l must be dropped, spaces included
		{"hello world  ", 11, "hello%20world"},

		// digits and symbols pass through untouched
		{"a1 b2  ", 5, "a1%20b2"},
		{"%20 x  ", 5, "%20%20x"},
		{"1 2 3    ", 5, "1%202%203"},

		// other whitespace is not a space and must not be replaced
		{"a\tb", 3, "a\tb"},
		{"a\tb c  ", 5, "a\tb%20c"},

		// longer input, many replacements
		{"the quick brown fox      ", 19, "the%20quick%20brown%20fox"},
	}

	for _, c := range cases {
		buf := []byte(c.in)
		if got := string(urlifyBytes(buf, c.l)); got != c.want {
			t.Errorf("urlifyBytes(%q, %d) = %q, want %q", c.in, c.l, got, c.want)
		}
	}
}

// The point of this version is that it writes into the caller's buffer instead of
// building a new one. Returning a correct-looking slice that lives somewhere else
// passes the table above but misses the exercise entirely, so check the backing
// array identity directly.
func TestUrlifyBytesWritesInPlace(t *testing.T) {
	buf := []byte("Mr John Smith    ")

	got := urlifyBytes(buf, 13)

	if len(got) == 0 {
		t.Fatal("urlifyBytes returned an empty slice, expected 17 bytes")
	}
	if &got[0] != &buf[0] {
		t.Error("urlifyBytes allocated a new backing array, expected an in-place write")
	}
	if string(buf) != "Mr%20John%20Smith" {
		t.Errorf("caller's buffer = %q, want %q", string(buf), "Mr%20John%20Smith")
	}
}

// The table above pads every buffer to the byte. Real callers do not: they hand you
// a scratch buffer, or round up, or reuse something they already had. A buffer with
// MORE than the required room must still work, which only holds if the write start
// and the returned length are both derived from the space count rather than read off
// len(buf). Everything past the result is leftover and must not be returned.
func TestUrlifyBytesOverPaddedBuffer(t *testing.T) {
	cases := []struct {
		in   string
		l    int
		want string
	}{
		// needs 17, given 23
		{"Mr John Smith          ", 13, "Mr%20John%20Smith"},
		// needs 5, given 12
		{"a b         ", 3, "a%20b"},
		// no spaces at all: needs 3, given 9
		{"abc      ", 3, "abc"},
		// needs 13, given 20
		{"hello world         ", 11, "hello%20world"},
		// needs 9, given 16
		{"1 2 3           ", 5, "1%202%203"},
	}

	for _, c := range cases {
		buf := []byte(c.in)
		if got := string(urlifyBytes(buf, c.l)); got != c.want {
			t.Errorf("urlifyBytes(%q, %d) = %q, want %q", c.in, c.l, got, c.want)
		}
	}
}
