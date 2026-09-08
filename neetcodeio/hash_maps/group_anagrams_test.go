package hashmaps

import (
	"math/rand"
	"sort"
	"strings"
	"testing"
)

var impls = []struct {
	name string
	fn   func([]string) [][]string
}{
	{"original", groupAnagrams},
	{"improved", groupAnagramsImproved},
	{"even_more", groupAnagramsEvenMoreImproved},
}

// map iteration order is random, so groups come back in any order, and words within a
// group come back in any order. Canonicalise both before comparing.
func canon(gs [][]string) string {
	out := make([]string, 0, len(gs))
	for _, g := range gs {
		c := append([]string(nil), g...)
		sort.Strings(c)
		out = append(out, strings.Join(c, ","))
	}
	sort.Strings(out)
	return strings.Join(out, " | ")
}

func TestGroupAnagrams(t *testing.T) {
	cases := []struct {
		in   []string
		want [][]string
	}{
		// samples from the statement
		{[]string{"act", "pots", "tops", "cat", "stop", "hat"},
			[][]string{{"act", "cat"}, {"pots", "tops", "stop"}, {"hat"}}},
		{[]string{"x"}, [][]string{{"x"}}},
		{[]string{""}, [][]string{{""}}},
		{[]string{"eat", "tea", "tan", "ate", "nat", "bat"},
			[][]string{{"eat", "tea", "ate"}, {"tan", "nat"}, {"bat"}}},

		// degenerate
		{[]string{}, nil},

		// duplicate words in the input - each copy must appear in the output
		{[]string{"ab", "ab"}, [][]string{{"ab", "ab"}}},
		{[]string{"a", "a", "a"}, [][]string{{"a", "a", "a"}}},
		{[]string{"ab", "ba", "ab"}, [][]string{{"ab", "ba", "ab"}}},

		// empty strings are anagrams of each other
		{[]string{"", "", "a"}, [][]string{{"", ""}, {"a"}}},

		// same letters, different counts -> NOT anagrams
		{[]string{"aab", "abb"}, [][]string{{"aab"}, {"abb"}}},
		{[]string{"a", "aa", "aaa"}, [][]string{{"a"}, {"aa"}, {"aaa"}}},

		// same length, no letters in common
		{[]string{"abc", "xyz"}, [][]string{{"abc"}, {"xyz"}}},

		// every word its own group / one big group
		{[]string{"a", "b", "c"}, [][]string{{"a"}, {"b"}, {"c"}}},
		{[]string{"abc", "bca", "cab", "acb"}, [][]string{{"abc", "bca", "cab", "acb"}}},
	}

	for _, impl := range impls {
		t.Run(impl.name, func(t *testing.T) {
			for _, c := range cases {
				in := append([]string(nil), c.in...)
				got, want := canon(impl.fn(in)), canon(c.want)
				if got != want {
					t.Errorf("%s(%v)\n got: %s\nwant: %s", impl.name, c.in, got, want)
				}
			}
		})
	}
}

// small alphabet so anagrams actually collide
func randWords(n, maxLen int, rng *rand.Rand) []string {
	out := make([]string, n)
	for i := range out {
		b := make([]byte, 1+rng.Intn(maxLen))
		for j := range b {
			b[j] = byte('a' + rng.Intn(4))
		}
		out[i] = string(b)
	}
	return out
}

func TestImplsAgree(t *testing.T) {
	rng := rand.New(rand.NewSource(3))
	for trial := 0; trial < 500; trial++ {
		in := randWords(1+rng.Intn(20), 5, rng)
		want := canon(groupAnagrams(append([]string(nil), in...)))
		for _, impl := range impls[1:] {
			got := canon(impl.fn(append([]string(nil), in...)))
			if got != want {
				t.Fatalf("%s disagrees on %v\n got: %s\nwant: %s", impl.name, in, got, want)
			}
		}
	}
}

// LeetCode constraints: up to 10^4 words, each up to 100 chars.
// The O(n^2) versions are the reason this is a benchmark and not a test.
func benchInput() []string {
	rng := rand.New(rand.NewSource(4))
	out := make([]string, 10000)
	for i := range out {
		b := make([]byte, 20)
		for j := range b {
			b[j] = byte('a' + rng.Intn(26))
		}
		out[i] = string(b)
	}
	return out
}

func BenchmarkGroupAnagrams(b *testing.B) {
	in := benchInput()
	for _, impl := range impls {
		b.Run(impl.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				impl.fn(in)
			}
		})
	}
}
