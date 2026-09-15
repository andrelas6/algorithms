package strings

import (
	"fmt"
	"testing"
)

// Encode's output format is your choice, so these tests assert the properties
// any correct encoding must have rather than an exact string. Once Decode
// exists, the round-trip test becomes the real check and most of this stays as
// a guard on the format itself.

func encodeOf(strs []string) string {
	s := &Solution{}
	return s.Encode(strs)
}

// Same input, same output. Nothing may depend on map iteration order, a clock,
// or anything else non-deterministic.
func TestEncodeIsDeterministic(t *testing.T) {
	inputs := [][]string{
		{},
		{""},
		{"Hello", "World"},
		{"a", "b", "c"},
		{"neet", "code", "love", "you"},
		{"#", "##", "###"},
	}

	for _, in := range inputs {
		first := encodeOf(in)
		for range 5 {
			if got := encodeOf(in); got != first {
				t.Errorf("Encode(%q) is not deterministic: %q then %q", in, first, got)
			}
		}
	}
}

// The whole point of the problem: two different lists must never encode to the
// same string. Every pair here is chosen to collide under some naive scheme -
// joining on a separator, ignoring empty strings, or assuming the separator
// cannot appear in the data.
func TestEncodeIsInjective(t *testing.T) {
	inputs := [][]string{
		// empty list vs a list holding one empty string vs two empty strings
		{},
		{""},
		{"", ""},
		{"", "", ""},

		// the classic separator collision
		{"a", "b"},
		{"ab"},
		{"a,b"},
		{"a#b"},
		{"a", "", "b"},

		// a payload that looks like an encoding
		{"1#a"},
		{"1", "a"},
		{"3#abc"},
		{"abc"},
		{"2#hi", "3#bye"},

		// separators as the entire payload
		{"#"},
		{"##"},
		{"#", "#"},
		{","},
		{",", ","},

		// digits next to data, where a length prefix could be misread
		{"12"},
		{"1", "2"},
		{"12", "34"},
		{"1234"},

		// order matters
		{"a", "b", "c"},
		{"c", "b", "a"},

		// leading and trailing empties
		{"a", ""},
		{"", "a"},

		// statement samples
		{"Hello", "World"},
		{"neet", "code", "love", "you"},

		// non-ASCII: the length prefix must agree with how the payload is
		// sliced back, bytes or runes, but consistently
		{"héllo"},
		{"h", "éllo"},
		{"日本語"},
		{"日", "本語"},
	}

	seen := make(map[string][]string, len(inputs))

	for _, in := range inputs {
		enc := encodeOf(in)

		if prev, ok := seen[enc]; ok {
			t.Errorf("Encode(%q) and Encode(%q) both produced %q; different lists must encode differently",
				prev, in, enc)
			continue
		}
		seen[enc] = in
	}
}

// Long strings and many strings: the length prefix stops being a single digit,
// which is where fixed-width assumptions break.
func TestEncodeHandlesMultiDigitLengths(t *testing.T) {
	cases := []struct {
		name string
		in   []string
	}{
		{"9 chars", []string{"123456789"}},
		{"10 chars", []string{"1234567890"}},
		{"99 chars", []string{str(99)}},
		{"100 chars", []string{str(100)}},
		{"1000 chars", []string{str(1000)}},
		{"mixed widths", []string{"a", str(12), str(345), "b"}},
		{"many strings", manyStrings(200)},
	}

	seen := make(map[string]string, len(cases))

	for _, c := range cases {
		enc := encodeOf(c.in)

		if enc == "" && len(c.in) > 0 {
			t.Errorf("%s: Encode returned an empty string for a non-empty list", c.name)
		}
		if prev, ok := seen[enc]; ok {
			t.Errorf("%s collided with %s", c.name, prev)
		}
		seen[enc] = c.name
	}
}

func str(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + i%26)
	}
	return string(b)
}

func manyStrings(n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = fmt.Sprintf("s%d", i)
	}
	return out
}

// Encode must not hold on to or modify the caller's slice contents.
func TestEncodeDoesNotModifyInput(t *testing.T) {
	in := []string{"Hello", "World", ""}
	before := append([]string(nil), in...)

	encodeOf(in)

	for i := range in {
		if in[i] != before[i] {
			t.Fatalf("Encode modified its input: %q became %q", before, in)
		}
	}
	if len(in) != len(before) {
		t.Fatalf("Encode changed the input length: %d, want %d", len(in), len(before))
	}
}

// nil and the empty slice are different values but the same list: no panic, and
// they may encode the same way.
func TestEncodeAcceptsNilAndEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Encode panicked on nil or empty input: %v", r)
		}
	}()

	nilEnc := encodeOf(nil)
	emptyEnc := encodeOf([]string{})

	if nilEnc != emptyEnc {
		t.Errorf("Encode(nil) = %q but Encode([]string{}) = %q; both are empty lists", nilEnc, emptyEnc)
	}
}

// --- Decode ---------------------------------------------------------------
//
// Decode's input is whatever Encode produced, so these tests go through the
// round trip rather than hand-written encoded strings. That keeps them working
// whatever format you settle on, and the round trip is the actual contract:
// decode(encode(x)) == x.

func decodeOf(encoded string) []string {
	s := &Solution{}
	return s.Decode(encoded)
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		in   []string
	}{
		// statement samples
		{"hello world", []string{"Hello", "World"}},
		{"neetcode", []string{"neet", "code", "love", "you"}},

		// degenerate
		{"empty list", []string{}},
		{"one empty string", []string{""}},
		{"two empty strings", []string{"", ""}},
		{"three empty strings", []string{"", "", ""}},
		{"single string", []string{"a"}},

		// empties mixed in, at every position
		{"empty first", []string{"", "a", "b"}},
		{"empty middle", []string{"a", "", "b"}},
		{"empty last", []string{"a", "b", ""}},

		// duplicates must survive, in order
		{"duplicates", []string{"a", "a", "a"}},
		{"duplicates interleaved", []string{"a", "b", "a", "b"}},

		// order must survive
		{"ascending", []string{"a", "b", "c"}},
		{"descending", []string{"c", "b", "a"}},

		// the separator inside the data
		{"hash only", []string{"#"}},
		{"many hashes", []string{"#", "##", "###"}},
		{"hash inside", []string{"a#b", "c#d"}},

		// payloads that look like encodings
		{"looks encoded", []string{"4#abc"}},
		{"looks encoded twice", []string{"2#hi", "3#bye"}},
		{"digits only", []string{"12", "34"}},
		{"digit then hash", []string{"1#"}},

		// multi-digit lengths
		{"nine chars", []string{"123456789"}},
		{"ten chars", []string{"1234567890"}},
		{"hundred chars", []string{str(100)}},
		{"thousand chars", []string{str(1000)}},
		{"mixed widths", []string{"a", str(12), str(345), "b"}},

		// many strings
		{"two hundred strings", manyStrings(200)},

		// non-ASCII: the prefix unit and the slicing unit must agree
		{"accented", []string{"héllo"}},
		{"split accented", []string{"h", "éllo"}},
		{"japanese", []string{"日本語"}},
		{"mixed scripts", []string{"日", "本語", "abc", ""}},
		{"emoji", []string{"🙂", "a🙂b"}},

		// whitespace and control characters are ordinary data
		{"spaces", []string{" ", "  ", "a b"}},
		{"newlines", []string{"a\nb", "\n"}},
		{"tabs", []string{"\t", "a\tb"}},
	}

	for _, c := range cases {
		got := decodeOf(encodeOf(c.in))

		if !sameStrings(got, c.in) {
			t.Errorf("%s: Decode(Encode(%q)) = %q, want %q", c.name, c.in, got, c.in)
		}
	}
}

// Randomised round trip over an alphabet stuffed with the characters most
// likely to break a scheme: digits, the separator, and multi-byte runes.
func TestRoundTripRandom(t *testing.T) {
	alphabet := []string{"a", "b", "#", "0", "1", "9", "é", "日", "🙂", " ", "\n"}

	seed := uint64(12345)
	next := func(n int) int {
		seed = seed*6364136223846793005 + 1442695040888963407
		return int((seed >> 33) % uint64(n))
	}

	for trial := 0; trial < 2000; trial++ {
		count := next(6)
		in := make([]string, count)
		for i := range in {
			length := next(8)
			s := ""
			for range length {
				s += alphabet[next(len(alphabet))]
			}
			in[i] = s
		}

		got := decodeOf(encodeOf(in))

		if !sameStrings(got, in) {
			t.Fatalf("round trip failed for %q: got %q", in, got)
		}
	}
}

// An empty encoding decodes to an empty list, not a list holding one empty
// string.
func TestDecodeOfEmptyEncoding(t *testing.T) {
	got := decodeOf(encodeOf(nil))

	if len(got) != 0 {
		t.Errorf("Decode(Encode(nil)) = %q, want an empty list", got)
	}
}

// Decoding is a pure read of the encoded string.
func TestDecodeDoesNotDependOnCallOrder(t *testing.T) {
	in := []string{"a", "#", "12", ""}
	enc := encodeOf(in)

	first := decodeOf(enc)
	for range 5 {
		if got := decodeOf(enc); !sameStrings(got, first) {
			t.Errorf("Decode(%q) is not deterministic: %q then %q", enc, first, got)
		}
	}
}
