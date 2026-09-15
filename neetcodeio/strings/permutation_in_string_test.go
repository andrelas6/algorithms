package strings

import (
	"bytes"
	"math/rand/v2"
	"slices"
	"testing"
)

// inclusionBruteForce is an independent oracle: sort the letters of every
// len(s1) window of s2 and compare with s1 sorted. No counting and no sliding,
// so a bug in the window bookkeeping cannot hide behind the same mistake here.
func inclusionBruteForce(s1, s2 string) bool {
	if len(s1) > len(s2) {
		return false
	}

	want := []byte(s1)
	slices.Sort(want)

	for start := 0; start+len(s1) <= len(s2); start++ {
		window := []byte(s2[start : start+len(s1)])
		slices.Sort(window)
		if bytes.Equal(window, want) {
			return true
		}
	}
	return false
}

func TestCheckInclusionExamples(t *testing.T) {
	cases := []struct {
		name   string
		s1, s2 string
		want   bool
	}{
		// samples from the statement
		{"neetcode sample 1", "abc", "lecabee", true},
		{"neetcode sample 2", "abc", "lecaabee", false},
		{"leetcode sample 1", "ab", "eidbaooo", true},
		{"leetcode sample 2", "ab", "eidboaoo", false},

		// degenerate
		{"single char present", "a", "a", true},
		{"single char absent", "a", "b", false},
		{"single char inside longer s2", "z", "abcz", true},
		{"s1 longer than s2", "abc", "ab", false},
		{"s1 much longer than s2", "abcdef", "a", false},

		// equal lengths: the only window is the whole of s2
		{"equal, identical", "abc", "abc", true},
		{"equal, rearranged", "abc", "cab", true},
		{"equal, one letter off", "abc", "abd", false},
	}

	for _, c := range cases {
		if got := checkInclusion(c.s1, c.s2); got != c.want {
			t.Errorf("%s: checkInclusion(%q, %q) = %v, want %v", c.name, c.s1, c.s2, got, c.want)
		}
	}
}

// The match can sit anywhere, including the very first and very last window. The
// last window is the one an off-by-one in the loop bound skips.
func TestCheckInclusionMatchPosition(t *testing.T) {
	cases := []struct {
		name   string
		s1, s2 string
		want   bool
	}{
		{"first window", "abc", "bcaxxxx", true},
		{"last window", "abc", "xxxxcba", true},
		{"middle window", "abc", "xxbacxx", true},
		{"second window", "ab", "xbax", true},
		{"second to last window", "ab", "xxbax", true},
		{"only after many slides", "xyz", "aaaaaaaaaaaaaaaaaaaazyx", true},
	}

	for _, c := range cases {
		if got := checkInclusion(c.s1, c.s2); got != c.want {
			t.Errorf("%s: checkInclusion(%q, %q) = %v, want %v", c.name, c.s1, c.s2, got, c.want)
		}
	}
}

// Near-misses: the right letters are all there, but not in one window of the
// right size, or with the wrong counts. Checking "does every letter of s1 appear
// somewhere" passes all of these; they should all be false.
func TestCheckInclusionNearMisses(t *testing.T) {
	cases := []struct {
		name   string
		s1, s2 string
	}{
		{"letters split by another", "ab", "acb"},
		{"letters far apart", "abc", "axxbxxc"},
		{"same letters, wrong counts", "aab", "abbxabb"},
		{"one extra copy needed", "aa", "aba"},
		{"right letters, one short", "abcd", "abcxbcdxacd"},
		{"window has a stranger", "abc", "abxc"},
		{"repeated letter only once in window", "aabb", "abxbaxab"},
	}

	for _, c := range cases {
		if checkInclusion(c.s1, c.s2) {
			t.Errorf("%s: checkInclusion(%q, %q) = true, want false", c.name, c.s1, c.s2)
		}
	}
}

// Repeated letters are where a count-per-letter window goes wrong: a letter
// leaving the window must drop by one, not vanish, and a letter entering must
// add to an existing count.
func TestCheckInclusionRepeatedLetters(t *testing.T) {
	cases := []struct {
		name   string
		s1, s2 string
		want   bool
	}{
		{"double letter", "aa", "baa", true},
		{"double letter at start", "aa", "aab", true},
		{"all one letter", "aaa", "aaaa", true},
		{"all one letter, not enough", "aaaa", "aaab", false},
		// the window slides past an 'a' while another 'a' stays inside
		{"left char still in window", "aab", "xaaxaab", true},
		{"s2 floods one letter", "abc", "aaaaaaaabc", true},
		{"s2 floods one letter, match reversed", "abc", "aaaaaaaacb", true},
		{"mixed counts", "aabbc", "xxcbabaxx", true},
		{"mixed counts, one short", "aabbc", "xxcbabxax", false},
	}

	for _, c := range cases {
		if got := checkInclusion(c.s1, c.s2); got != c.want {
			t.Errorf("%s: checkInclusion(%q, %q) = %v, want %v", c.name, c.s1, c.s2, got, c.want)
		}
	}
}

// Every s1 against every s2 over a small alphabet, compared with the brute force.
// A tiny alphabet forces lots of repeats and near-misses; a fixed seed keeps it
// repeatable.
func TestCheckInclusionAgainstBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(5, 6))

	randomString := func(n int, alphabet string) string {
		b := make([]byte, n)
		for i := range b {
			b[i] = alphabet[rng.IntN(len(alphabet))]
		}
		return string(b)
	}

	for range 3000 {
		alphabet := []string{"ab", "abc", "abcdefghijklmnopqrstuvwxyz"}[rng.IntN(3)]
		s1 := randomString(1+rng.IntN(6), alphabet)
		s2 := randomString(1+rng.IntN(20), alphabet)

		got, want := checkInclusion(s1, s2), inclusionBruteForce(s1, s2)
		if got != want {
			t.Fatalf("checkInclusion(%q, %q) = %v, brute force = %v", s1, s2, got, want)
		}
	}
}

// Every lowercase letter at once, so the count map is as big as it gets.
func TestCheckInclusionFullAlphabet(t *testing.T) {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	reversed := []byte(alphabet)
	slices.Reverse(reversed)

	cases := []struct {
		name   string
		s1, s2 string
		want   bool
	}{
		{"reversed alphabet", alphabet, string(reversed), true},
		{"alphabet after padding", alphabet, "zzzz" + string(reversed), true},
		{"alphabet missing one letter", alphabet, "abcdefghijklmnopqrstuvwxy" + "y", false},
	}

	for _, c := range cases {
		if got := checkInclusion(c.s1, c.s2); got != c.want {
			t.Errorf("%s: checkInclusion(%q, %q) = %v, want %v", c.name, c.s1, c.s2, got, c.want)
		}
	}
}

// Constraints allow both strings up to 10^4 long. Answers here are known by
// construction, so no brute force is needed at this size.
func TestCheckInclusionAtMaxSize(t *testing.T) {
	const n = 10_000

	repeat := func(s string, count int) string {
		return string(bytes.Repeat([]byte(s), count))
	}

	cases := []struct {
		name   string
		s1, s2 string
		want   bool
	}{
		// match only in the very last window, after ~10^4 slides
		{"short s1, match at the end", "abc", repeat("x", n-3) + "cba", true},
		{"short s1, no match", "abc", repeat("ab", n/2), false},
		// equal lengths, 10^4 each
		{"equal length, rearranged", repeat("ab", n/2), repeat("ba", n/2), true},
		{"equal length, one letter off", repeat("ab", n/2), repeat("ba", n/2-1) + "bb", false},
		// half-size s1 sliding across the whole of s2
		{"half size s1, match at the end", repeat("a", n/2), repeat("b", n/2) + repeat("a", n/2), true},
		{"half size s1, one short", repeat("a", n/2), repeat("ab", n/2), false},
	}

	for _, c := range cases {
		if got := checkInclusion(c.s1, c.s2); got != c.want {
			t.Errorf("%s: checkInclusion(len %d, len %d) = %v, want %v", c.name, len(c.s1), len(c.s2), got, c.want)
		}
	}
}
